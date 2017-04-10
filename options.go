package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"reflect"
	"strconv"
	"time"

	"github.com/jason0x43/go-alfred"
)

// OptionsCommand shows and sets configuration options
type OptionsCommand struct{}

// About returns information about a command
func (c OptionsCommand) About() alfred.CommandDef {
	return alfred.CommandDef{
		Keyword:     "options",
		Description: "Workflow options",
		IsEnabled:   true,
	}
}

// Items ...
func (c OptionsCommand) Items(arg, data string) (items []alfred.Item, err error) {
	ct := reflect.TypeOf(config)
	cfg := reflect.Indirect(reflect.ValueOf(config))

	for i := 0; i < ct.NumField(); i++ {
		field := ct.Field(i)
		desc := field.Tag.Get("desc")
		if desc == "" {
			continue
		}

		name, value := alfred.SplitCmd(arg)
		if !alfred.FuzzyMatches(field.Name, name) {
			continue
		}

		switch field.Name {
		case "Service":
			if name == "Service" {
				if alfred.FuzzyMatches(string(serviceDarkSky), value) {
					items = append(items, makeStringChoice("Service", string(serviceDarkSky)))
				}

				if alfred.FuzzyMatches(string(serviceWunderground), value) {
					items = append(items, makeStringChoice("Service", string(serviceWunderground)))
				}

				return
			}

			items = append(items, alfred.Item{
				Title:        fmt.Sprintf("Service: %v", config.Service),
				Autocomplete: "Service ",
				Subtitle:     desc,
			})

		case "Units":
			if name == "Units" {
				if alfred.FuzzyMatches(string(unitsMetric), value) {
					items = append(items, makeStringChoice("Units", string(unitsMetric)))
				}

				if alfred.FuzzyMatches(string(unitsUS), value) {
					items = append(items, makeStringChoice("Units", string(unitsUS)))
				}

				return
			}

			items = append(items, alfred.Item{
				Title:        fmt.Sprintf("Units: %v", config.Units),
				Autocomplete: "Units ",
				Subtitle:     desc,
			})

		case "Location":
			if name == "Location" {
				if value == "" {
					items = append(items, alfred.Item{
						Title:    "Location: " + config.Location.Name,
						Subtitle: "Enter a new city/state or ZIP",
					})
				} else {
					var location Geocode
					if location, err = Locate(value); err != nil {
						return
					}

					opts := config
					o := reflect.Indirect(reflect.ValueOf(&opts))
					o.FieldByName("Location").Set(reflect.ValueOf(location.Location()))

					items = append(items, alfred.Item{
						Title:    location.Name,
						Subtitle: fmt.Sprintf("(%f, %f)", location.Latitude, location.Longitude),
						Arg: &alfred.ItemArg{
							Keyword: "options",
							Mode:    alfred.ModeDo,
							Data:    alfred.Stringify(&opts),
						},
					})
				}

				return
			}

			items = append(items, alfred.Item{
				Title:        "Location: " + config.Location.Name,
				Subtitle:     desc,
				Autocomplete: "Location ",
			})

		case "Icons":
			if name == "Icons" {
				var dirs []os.FileInfo
				if dirs, err = ioutil.ReadDir("icons"); err != nil {
					return
				}
				for _, dir := range dirs {
					if dir.IsDir() && alfred.FuzzyMatches(dir.Name(), value) {
						items = append(items, makeIconChoice("Icons", dir.Name()))
					}
				}
				return
			}

			items = append(items, alfred.Item{
				Title:        "Icons: " + config.Icons,
				Subtitle:     desc,
				Autocomplete: "Icons ",
			})

		case "DateFormat":
			if name == "DateFormat" {
				for _, tf := range DateFormats {
					items = append(items, makeStringChoice("DateFormat", tf))
				}
				return
			}

			items = append(items, alfred.Item{
				Title:        "DateFormat: " + config.DateFormat,
				Subtitle:     desc,
				Autocomplete: "DateFormat ",
			})

		case "TimeFormat":
			if name == "TimeFormat" {
				for _, tf := range TimeFormats {
					items = append(items, makeStringChoice("TimeFormat", tf))
				}
				return
			}

			items = append(items, alfred.Item{
				Title:        "TimeFormat: " + config.TimeFormat,
				Subtitle:     desc,
				Autocomplete: "TimeFormat ",
			})

		default:
			item := alfred.Item{
				Title:        field.Name,
				Subtitle:     desc,
				Autocomplete: field.Name,
			}

			itemArg := &alfred.ItemArg{
				Keyword: "options",
				Mode:    alfred.ModeDo,
			}

			switch field.Type.Name() {
			case "bool":
				f := cfg.FieldByName(field.Name)
				if name == field.Name {
					item.Title += " (press Enter to toggle)"
				}

				// copy the current options, update them, and use as the arg
				opts := config
				o := reflect.Indirect(reflect.ValueOf(&opts))
				newVal := !f.Bool()
				o.FieldByName(field.Name).SetBool(newVal)
				item.Arg = itemArg
				item.Arg.Data = alfred.Stringify(&opts)
				item.AddCheckBox(f.Bool())

			case "int":
				item.Autocomplete += " "

				if value != "" {
					val, err := strconv.Atoi(value)
					if err != nil {
						return items, err
					}
					item.Title += fmt.Sprintf(": %d", val)

					// copy the current options, update them, and use as the arg
					opts := config
					o := reflect.Indirect(reflect.ValueOf(&opts))
					o.FieldByName(field.Name).SetInt(int64(val))
					item.Arg = itemArg
					item.Arg.Data = alfred.Stringify(opts)
				} else {
					f := cfg.FieldByName(field.Name)
					val := f.Int()
					item.Title += fmt.Sprintf(": %v", val)
					if name == field.Name {
						item.Title += " (type a new value to change)"
					}
				}

			case "string":
				f := cfg.FieldByName(field.Name)
				item.Autocomplete += " "
				item.Title += ": " + f.String()

				if name == field.Name {
					opts := config
					o := reflect.Indirect(reflect.ValueOf(&opts))
					o.FieldByName(field.Name).SetString(value)
					item.Arg = itemArg
					item.Arg.Data = alfred.Stringify(&opts)
				}
			}

			items = append(items, item)
		}
	}

	alfred.FuzzySort(items, arg)

	return
}

// Do ...
func (c OptionsCommand) Do(data string) (out string, err error) {
	if err = json.Unmarshal([]byte(data), &config); err != nil {
		return
	}

	if err = alfred.SaveJSON(configFile, &config); err != nil {
		log.Printf("Error saving config: %s\n", err)
		return "Error updating config", err
	}

	// Clear the cache to allow data to be requestsed with the new options
	cache.Time = time.Unix(0, 0)
	if err = alfred.SaveJSON(cacheFile, &cache); err != nil {
		log.Printf("Error saving cache: %s\n", err)
	}

	return "Updated config", err
}

func makeStringChoice(fieldName, value string) alfred.Item {
	opts := config
	o := reflect.Indirect(reflect.ValueOf(&opts))
	currentValue := o.FieldByName(fieldName).String()
	o.FieldByName(fieldName).SetString(value)
	item := alfred.Item{
		Title: value,
		Arg: &alfred.ItemArg{
			Keyword: "options",
			Mode:    alfred.ModeDo,
			Data:    alfred.Stringify(&opts),
		},
	}
	item.AddCheckBox(currentValue == value)
	return item
}

func makeIconChoice(fieldName, value string) alfred.Item {
	item := makeStringChoice(fieldName, value)
	item.Icon = path.Join("icons", value, "tstorms.png")
	return item
}
