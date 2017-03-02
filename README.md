# PirriGo
An (in-progress) rewrite of the Pirri irrigation and garden management application, in Go and Angular.

- Live Demo: http://backwardhat.com:8001
- login/password are demo/demo

## Features
- Activate irrigation zones on a schedule, or manually
- Track total seconds spent watering by zone, day of the week, and other metrics
- Estimate monthly costs to operate your irrigation system
- Using WeatherUnderground, leverage current weather data to automatically regulate watering (increase in heat waves, reduce when rain will do the job for you)


## Notes

#### TODO
- Provide hardware assembly instructions
- Better automate initial stand-up of database
- Finish making the front end less offensive.
- Make the calendar interactive (like in the original version, but prettier and more functional)

#### Manual submission of a station run with curl 

(when not using authentication)
```
curl -X POST -d "{\"StationID\": 2, \"Duration\": 60}" http://localhost:8001/station/run
```
