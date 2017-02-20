# PirriGo
An (in-progress) rewrite of the Pirri irrigation and garden management application, in Go.

## Notes

#### Manual submission of a station run with curl
```
curl -X POST -d "{\"StationID\": 2, \"Duration\": 15}" http://localhost:8001/station/run
```
#### TODO
- Create feature list
- Include schematic for building the hardware
- Autopopulate GPIOs in database
- Finish refactoring old frontend so it talks to new backend.
- Add new calendar interface (preferably something Angular-y)




