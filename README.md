# PirriGo
An (in-progress) rewrite of the Pirri irrigation and garden management application, in Go and Angular.

- Live Demo: http://pirri.vacovsky.us
- login/password are demo/demo

## Screenshots


## Features
- Activate irrigation zones on a schedule, or manually
- Track total seconds spent watering by zone, day of the week, and other metrics
- Estimate monthly costs to operate your irrigation system
- Using WeatherUnderground, leverage current weather data to automatically regulate watering (increase in heat waves, reduce when rain will do the job for you)

## Setup

### Clone repository to your Pi
```bash
git clone https://github.com/vacoj/pirrigo.git
```

### Installing Go in Raspbian
```bash
wget https://storage.googleapis.com/golang/go1.7.linux-armv6l.tar.gz
sudo tar -C /usr/local -xzf go1.7.linux-armv6l.tar.gz
export PATH=/usr/local/go/bin:$PATH
```

### Install required Go packages
```bash
go get "github.com/davecgh/go-spew/spew"
go get "github.com/jinzhu/gorm"
go get "github.com/jinzhu/gorm/dialects/mysql"
go get "github.com/newrelic/go-agent"
go get "github.com/stianeikeland/go-rpio"
go get "github.com/streadway/amqp"
go get "gopkg.in/redis.v5"
```

### Install javascript libraries for web
```bash
cd pirrigo/src/pirri/static
bower install
```

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


##### In case you're feeling generous

Bitcoin address: 17FC2XP5phFFupAn2w4dGCwQuDUWRXE4EQ
