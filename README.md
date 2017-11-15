# PirriGo
Irrigation and garden management application, in Go and Angular.  Based on <a href="https://github.com/vacoj/pirri">Pirri</a>.

## Live Demo: http://pirri.vacovsky.us
#### login/password are demo / demo

## Screenshots

### Login page

<a href="https://raw.githubusercontent.com/vacoj/pirrigo/master/screenshots/login.png"><img src="https://raw.githubusercontent.com/vacoj/pirrigo/master/screenshots/login.png" alt="Stations" height="150"></a>

### Main / Stations

<a href="https://raw.githubusercontent.com/vacoj/pirrigo/master/screenshots/main.png"><img src="https://raw.githubusercontent.com/vacoj/pirrigo/master/screenshots/main.png" alt="Stations" height="150"></a>

### Calendar

<a href="https://raw.githubusercontent.com/vacoj/pirrigo/master/screenshots/calendar.png"><img src="https://raw.githubusercontent.com/vacoj/pirrigo/master/screenshots/calendar.png" alt="Stations"  height="150"></a>

### Station Scheduler

<a href="https://raw.githubusercontent.com/vacoj/pirrigo/master/screenshots/scheduler.png"><img src="https://raw.githubusercontent.com/vacoj/pirrigo/master/screenshots/scheduler.png" alt="Stations"  height="150"></a>

### History of Activity

<a href="https://raw.githubusercontent.com/vacoj/pirrigo/master/screenshots/history.png"><img src="https://raw.githubusercontent.com/vacoj/pirrigo/master/screenshots/history.png" alt="Stations" height="150"></a>

### Charted Run Stats

<a href="https://raw.githubusercontent.com/vacoj/pirrigo/master/screenshots/charts.png"><img src="https://raw.githubusercontent.com/vacoj/pirrigo/master/screenshots/charts.png" alt="Stations" height="150"></a>

### Usage

<a href="https://raw.githubusercontent.com/vacoj/pirrigo/master/screenshots/usage.png"><img src="https://raw.githubusercontent.com/vacoj/pirrigo/master/screenshots/usage.png" alt="Stations" height="150"></a>

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
go get "go.uber.org/zap"
```

### Install javascript libraries for web

```bash
cd pirrigo/src/pirri/static
bower install
```

### Setting up SQL for testing

```bash
docker run --name sqldumpster -p 3306:3306 -e MYSQL_ROOT_PASSWORD=root -d mysql:8
```
```sql
create database pirritesting;
GRANT ALL PRIVILEGES ON *.* TO 'pirritesting'@'%' IDENTIFIED BY 'pirritesting';
FLUSH PRIVILEGES;
```

### Configuration file (<a href="https://github.com/vacoj/pirrigo/blob/master/init/appconfig.json">example</a>)

```javascript
{
    "pirri": {
        "welcome_message":"L/P is demo/demo",
        "utc_offset": -8,
        "monitor_interval": 60,
        "rain_skip": true,
        "use_rabbit": false,
        "date_format": "2006-01-02 15:04:05"
    },
    "rabbitmq": {
        "server": "192.168.111.50",
        "port": "5672",
        "user": "rabbit",
        "secret": "bunnyrabbit!!",
        "task_queue": "pirri_go_work_dev"
    },
    "sql": {
        "dbtype": "mysql",
        "server": "localhost",
        "user": "pirritesting",
        "secret": "pirritesting",
        "db": "pirritesting",
        "port": "3306"
    },
    "newrelic": {
        "active": false,
        "license_path": "../../../nr_key"
    },
    "redis": {
        "host": "localhost",
        "port": "6379"
    },
    "web": {
        "port": "8001",
        "user": "demo",
        "secret": "demo"
    },
    "weather": {
        "wunderground_key": "",
        "state_abbreviation": "CA",
        "city": "Atascadero",
        "units": "imperial",
        "station": "KCAATASC69"
    },
    "debug": {
        "pirri": true,
        "gorm": false,
        "settings": false,
        "simulate_gpio": true,
        "log_path": "pirrigo.log"
    }
}
```


### Running app

```bash
cd pirrigo/src/pirri
./pirri /path/to/config.json
```

## Notes

#### TODO

- Provide hardware assembly instructions
- Better automate initial stand-up of database
- Finish making the front end less offensive.
- Make the calendar interactive (like in the original version, but prettier and more functional)

#### Manual submission of a station run with curl 

(when not using authentication)

``` bash
curl -X POST -d "{\"StationID\": 2, \"Duration\": 60}" http://localhost:8001/station/run
```


##### In case you're feeling generous

Bitcoin address: 17FC2XP5phFFupAn2w4dGCwQuDUWRXE4EQ
