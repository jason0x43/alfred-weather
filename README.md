alfred-weather
==============

Deprecation Notice
------------------

The Wunderground API has been [shut down since 2018-09-06 per their announcement](https://apicommunity.wunderground.com/weatherapi/topics/end-of-service-for-the-weather-underground-api).

Additionally, [Dark Sky was purchased by Apple](https://blog.darksky.net/dark-sky-has-a-new-home/), and has sunset their API, issueing no new keys, and will cease servie at the end of 2021.

As a result, this Alfred Workflow can only be enjoyed by anyone who already has a Dark Sky API key.

* * * 

An [Alfred][alfred] workflow for showing weather forecasts

![Screenshot](doc/daily.png?raw=true)

This workflow lets you access weather forecasts from [Dark Sky][darksky] and the
[Weather Underground][wund]. 

Installation
------------

Download the latest workflow package from the [releases page](https://github.com/jason0x43/alfred-weather/releases) and double click it — Alfred will take care of the rest.

Setup
-----

The workflow has one top-level command, `wtr`, and three sub-commands, daily (`wtd`), hourly (`wth`), and options (`wto`). The first thing you'll need to do is configure some options.

![Options](doc/options.png?raw=true)

Some options are, well, optional, but the Service and related Key options, and a default location, are required. You need an API key for the service you choose. Both of the currently supported services (Dark Sky and Weather Underground) are free to use (for a reasonable number of requests per day).

* [Dark Sky API](https://darksky.net/dev/)
* [Weather Underground API](https://www.wunderground.com/member/registration?mode=api_signup)

Once you've entered the service key, selection the "Location" option then enter a ZIP code or city name, then wait a couple of seconds. When it looks like your desired location has been found, press Enter to save it.

Usage
-----

The `wtd` keyword will show a forecast for the next several days.

![Daily forecast](doc/daily.png?raw=true)

The `wth` keyword will show a forecast for the next several hours.

![Hourly forecast](doc/hourly.png?raw=true)

In either case, you can enter a location query to get the forecast for somewhere other than your default location.

![Name query](doc/daily_name.png?raw=true)

![ZIP query](doc/daily_zip.png?raw=true)

Actioning a day in the daily forecast will jump to an hourly forecast for that day, if hourly data is available. Actioning the list heading will jump back to the daily forecast.

If there are any active weather alerts, they'll show at the top of the forecast. Actioning an alert will open more detailed information in a browser window.

If there is a newer version of the workflow available, a message will be displayed at the top of the result list. Actioning it will open a release page for the new version in a browser window.

![Update notice](doc/update.png?raw=true)


Credits
-------

The package includes a number of icon sets from the Weather Underground and
from [weathericonsets.com][icons] (I'm not up to drawing weather icons yet).
Each set includes an `info.json` file that gives a short description and
provides a source URL for the icon set.

[alfred]: http://www.alfredapp.com
[icons]: http://www.weathericonsets.com
[wund]: http://www.weatherunderground.com
[darksky]: http://darksky.net
