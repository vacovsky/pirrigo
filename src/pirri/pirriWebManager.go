package main

import (
	"encoding/json"
	//	"fmt"
	"io"
	//	"io/ioutil"
	"log"
	"net/http"
	"runtime"

	//	"github.com/davecgh/go-spew/spew"
)

func logTraffic() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}

func stationRun(rw http.ResponseWriter, req *http.Request) {
	//	var t Task
	var msr ManualStationRun

	dec := json.NewDecoder(req.Body)
	for {
		if err := dec.Decode(&msr); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		log.Printf("%s\n", msr.StationID)
	}

	//	fmt.Println(logTraffic())
	//	defer req.Body.Close()

	//	//	body, ERR := ioutil.ReadAll(req.Body)
	//	//	log.Println(string(body))

	//	ERR = json.NewDecoder(req.Body).Decode(&msr)
	//	//	ERR = json.Unmarshal(body, &msr)

	//	GormDbConnect()
	//	defer db.Close()

	//	if SETTINGS.PirriDebug {
	//		//		spew.Dump(req.Body)
	//		spew.Dump(t)
	//		spew.Dump(msr)
	//	}

	//	db.Where("id = ?", msr.StationID).Find(&t.Station)
	//	t.StationSchedule = StationSchedule{Duration: msr.Duration}

	//	blob, ERR := json.Marshal(&t)
	//	failOnError(ERR, ERR.Error())

	//	if SETTINGS.PirriDebug {
	//		fmt.Println(string(blob))
	//	}
	//	io.WriteString(rw, string(blob))
}
