# pirrigo

An experimental rewrite of the pirri backend application for controlling irrigation and sprinkers, but this time with Go


#### Generate DB from scratch
See SqlBootstrap.sql

#### Manual submission of a station run with curl
```
curl -X POST -d "{\"StationID\": 2, \"Duration\": 15}" http://localhost:8001/station/run
```