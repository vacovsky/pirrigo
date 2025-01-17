import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { GlobalsService } from './globals.service'
import { Observable } from 'rxjs';
import { StationResponse, StationRunRequestBody, Station, StationStatus, StationRunJob } from '../structs/station';
import { StationHistoryResponse } from 'src/app/structs/station-history';
import { StationLogsResponse } from '../structs/station-logs';
import { StationSchedule, StationScheduleResponse } from '../structs/station-schedule';
import { Gpio, GpioResponse } from '../structs/gpio';
import { ChartData } from '../structs/chart-data';

@Injectable()
export class ApiClientService {

  constructor(
    private _http: HttpClient,
    private _globals: GlobalsService
  ) { }

  getRunHistory(station: any, earliest: number): Observable<StationHistoryResponse> {
    const uri = `${this._globals.uriStem}/history?station=${station}&earliest=${earliest}`;
    return this._http.get<StationHistoryResponse>(uri)
  }

  getStationLogs(): Observable<StationLogsResponse> {
    const uri = `${this._globals.uriStem}/logs/all`
    return this._http.get<StationLogsResponse>(uri)
  }

  getAllStations(): Observable<StationResponse> {
    const uri = `${this._globals.uriStem}/station/all`
    return this._http.get<StationResponse>(uri)
  }

  getStationSchedules(): Observable<StationScheduleResponse> {
    const uri = `${this._globals.uriStem}/schedule/all`
    return this._http.get<StationScheduleResponse>(uri)
  }

  getStationStatus(): Observable<StationStatus> {
    const uri = `${this._globals.uriStem}/status/run`
    return this._http.get<StationStatus>(uri)
  }

  postStationRun(body: StationRunRequestBody): void {
    const uri = `${this._globals.uriStem}/station/run`
    // TODO: return response and parse for errors
    this._http.post(
      uri, body, { headers: this._globals.headers }
    ).subscribe(() => {
      // console.log("Running Station:", body.StationID)
    })
  }

  cancelActiveStationRun() {
    const uri = `${this._globals.uriStem}/status/cancel`
    return this._http.get<StationStatus>(uri)
  }

  getStationRunQueue(): Observable<StationRunJob[]> {
    const uri = `${this._globals.uriStem}/status/queue`
    return this._http.get<StationRunJob[]>(uri)
  }

  cancelQueuedJob(body: any): Observable<StationRunJob[]> {
    const uri = `${this._globals.uriStem}/status/queue/remove`
    return this._http.post<StationRunJob[]>(
      uri, body, { headers: this._globals.headers }
    )
  }

  postStationScheduleChange(body: StationSchedule): Observable<StationScheduleResponse> {
    const uri = `${this._globals.uriStem}/schedule/edit`
    return this._http.post<StationScheduleResponse>(uri, body, { headers: this._globals.headers })
  }

  getStationForScheduleEdit(stationID: number): Observable<Station> {
    const uri = `${this._globals.uriStem}/station?stationid=${stationID}`
    return this._http.get<Station>(uri)
  }

  deleteStationScheduleItem(id: number): Observable<StationScheduleResponse> {
    const uri = `${this._globals.uriStem}/schedule/delete`
    return this._http.post<StationScheduleResponse>(
      uri, { "ID": id }, { headers: this._globals.headers }
    )
  }

  postStationChange(station: Station): Observable<any> {
    const uri = `${this._globals.uriStem}/station/edit`
    return this._http.post(uri, station, { headers: this._globals.headers })
  }

  getGPIOsForStationEdit(): Observable<GpioResponse> {
    const uri = `${this._globals.uriStem}/gpio/available`
    return this._http.get<GpioResponse>(uri)
  }


  getAllGPIOs(): Observable<GpioResponse> {
    const uri = `${this._globals.uriStem}/gpio/all`
    return this._http.get<GpioResponse>(uri)
  }



  deleteStation(id: number): Observable<StationScheduleResponse> {
    const uri = `${this._globals.uriStem}/station/delete`
    return this._http.post<StationScheduleResponse>(
      uri, { "ID": id }, { headers: this._globals.headers }
    )
  }

  setCommonWireGpio(gpio: number): Observable<GpioResponse> {
    const uri = `${this._globals.uriStem}/gpio/common/set`
    return this._http.post<GpioResponse>(
      uri, { "GPIO": gpio }, { headers: this._globals.headers }
    )
  }


  loadChartByID(chart: number, startDate: number, endDate: number): Observable<ChartData> {
    const uri = `${this._globals.uriStem}/stats/${4}`;
    // const uri = `${this._globals.uriStem}/tempchartdata?startTime=${startDate}&endTime=${endDate}`;
    return this._http.get<ChartData>(uri)
  }
  // /gpio/common/set


  // schedule/delete {ID: 134} POST

  // /schedule/edit POST
  //   stationSchedules: [{ID: 1, StartDate: "2017-03-10T17:08:40Z", EndDate: "2027-03-10T08:00:00Z", Sunday: true,…}
  // Duration: 1800
  // EndDate: "2027-03-10T08:00:00Z"
  // Friday: false
  // ID: 1
  // Monday: false
  // Repeating: false
  // Saturday: true
  // StartDate: "2017-03-10T17:08:40Z"
  // StartTime: 100
  // StationID: 4
  // Sunday: true
  // Thursday: true
  // Tuesday: true
  // Wednesday: false

}

// schedule/all
// { "stationSchedules": [{"ID":1,"StartDate":"2017-03-10T17:08:40Z","EndDate":"2027-03-10T08:00:00Z","Sunday":true,"Monday":false,"Tuesday":true,"Wednesday":false,"Thursday":true,"Friday":false,"Saturday":true,"StationID":4,"StartTime":100,"Duration":1800,"Repeating":false},{"ID":2,"StartDate":"2017-03-11T00:18:22Z","EndDate":"2027-03-10T08:00:00Z","Sunday":true,"Monday":false,"Tuesday":true,"Wednesday":false,"Thursday":true,"Friday":false,"Saturday":false,"StationID":5,"StartTime":100,"Duration":1800,"Repeating":false},{"ID":3,"StartDate":"2017-03-11T00:18:55Z","EndDate":"2027-03-10T08:00:00Z","Sunday":true,"Monday":false,"Tuesday":true,"Wednesday":false,"Thursday":true,"Friday":false,"Saturday":true,"StationID":6,"StartTime":100,"Duration":1800,"Repeating":false},{"ID":4,"StartDate":"2017-03-11T08:00:25Z","EndDate":"2027-03-11T08:00:00Z","Sunday":true,"Monday":false,"Tuesday":true,"Wednesday":false,"Thursday":true,"Friday":false,"Saturday":false,"StationID":7,"StartTime":100,"Duration":1800,"Repeating":false},{"ID":5,"StartDate":"2017-03-11T08:01:01Z","EndDate":"2027-03-11T08:00:00Z","Sunday":true,"Monday":false,"Tuesday":false,"Wednesday":true,"Thursday":false,"Friday":true,"Saturday":false,"StationID":9,"StartTime":100,"Duration":300,"Repeating":false},{"ID":6,"StartDate":"2017-03-11T08:01:23Z","EndDate":"2027-03-11T08:00:00Z","Sunday":true,"Monday":false,"Tuesday":false,"Wednesday":true,"Thursday":false,"Friday":true,"Saturday":false,"StationID":10,"StartTime":100,"Duration":300,"Repeating":false},{"ID":7,"StartDate":"2017-03-11T08:01:52Z","EndDate":"2027-03-11T08:00:00Z","Sunday":true,"Monday":false,"Tuesday":true,"Wednesday":false,"Thursday":true,"Friday":false,"Saturday":false,"StationID":17,"StartTime":100,"Duration":300,"Repeating":false},{"ID":9,"StartDate":"2017-03-11T08:02:38Z","EndDate":"2027-03-11T08:00:00Z","Sunday":true,"Monday":false,"Tuesday":true,"Wednesday":false,"Thursday":false,"Friday":true,"Saturday":false,"StationID":16,"StartTime":100,"Duration":300,"Repeating":false},{"ID":10,"StartDate":"2017-03-12T17:33:13Z","EndDate":"2027-03-12T08:00:00Z","Sunday":true,"Monday":false,"Tuesday":true,"Wednesday":false,"Thursday":false,"Friday":true,"Saturday":false,"StationID":15,"StartTime":1900,"Duration":300,"Repeating":false},{"ID":100,"StartDate":"2020-01-01T01:01:01.000000001Z","EndDate":"2150-01-01T01:01:01.000000001Z","Sunday":false,"Monday":false,"Tuesday":false,"Wednesday":false,"Thursday":false,"Friday":false,"Saturday":false,"StationID":100,"StartTime":0,"Duration":10,"Repeating":false},{"ID":104,"StartDate":"2017-03-12T17:33:13Z","EndDate":"2031-05-02T07:00:00Z","Sunday":true,"Monday":false,"Tuesday":true,"Wednesday":false,"Thursday":false,"Friday":true,"Saturday":false,"StationID":14,"StartTime":1930,"Duration":300,"Repeating":false},{"ID":105,"StartDate":"2017-03-12T17:33:13Z","EndDate":"2031-05-02T07:00:00Z","Sunday":true,"Monday":false,"Tuesday":true,"Wednesday":false,"Thursday":false,"Friday":true,"Saturday":false,"StationID":13,"StartTime":2000,"Duration":300,"Repeating":false}]}

// logs/all
// { "logs": [{"level":"debug","message":"Starting monitoring at interval","application":"PirriGo","interval":59,"time":""},{"level":"debug","message":"Queuing Task for GPIO activation in OfflineRunQueue for station","application":"PirriGo","gpio":16,"time":""},{"level":"debug","message":"Queuing Task for GPIO activation in OfflineRunQueue for station","application":"PirriGo","gpio":18,"time":""},{"level":"debug","message":"Executing task for station","application":"PirriGo","stationID":10,"time":""},{"level":"debug","message":"Logging task for station","application":"PirriGo","stationID":10,"startTime":100,"time":""},{"level":"debug","message":"Activating GPIOs","application":"PirriGo","commonWire":21,"gpio":18,"durationSeconds":300,"time":""},{"level":"debug","message":"Queuing Task for GPIO activation in OfflineRunQueue for station","application":"PirriGo","gpio":16,"time":""},{"level":"debug","message":"Queuing Task for GPIO activation in OfflineRunQueue for station","application":"PirriGo","gpio":18,"time":""},{"level":"debug","message":"Deactivating GPIOs","application":"PirriGo","commonWire":21,"gpio":18,"durationSeconds":0,"time":""},{"level":"debug","message":"Task execution complete for station","application":"PirriGo","stationID":10,"time":""},{"level":"debug","message":"Executing task for station","application":"PirriGo","stationID":10,"time":""},{"level":"debug","message":"Logging task for station","application":"PirriGo","stationID":10,"startTime":100,"time":""},{"level":"debug","message":"Activating GPIOs","application":"PirriGo","commonWire":21,"gpio":18,"durationSeconds":300,"time":""},{"level":"debug","message":"Deactivating GPIOs","application":"PirriGo","commonWire":21,"gpio":18,"durationSeconds":0,"time":""},{"level":"debug","message":"Task execution complete for station","application":"PirriGo","stationID":10,"time":""},{"level":"debug","message":"Executing task for station","application":"PirriGo","stationID":9,"time":""},{"level":"debug","message":"Logging task for station","application":"PirriGo","stationID":9,"startTime":100,"time":""},{"level":"debug","message":"Activating GPIOs","application":"PirriGo","commonWire":21,"gpio":16,"durationSeconds":300,"time":""},{"level":"debug","message":"Deactivating GPIOs","application":"PirriGo","commonWire":21,"gpio":16,"durationSeconds":0,"time":""},{"level":"debug","message":"Task execution complete for station","application":"PirriGo","stationID":9,"time":""},{"level":"debug","message":"Executing task for station","application":"PirriGo","stationID":9,"time":""},{"level":"debug","message":"Logging task for station","application":"PirriGo","stationID":9,"startTime":100,"time":""},{"level":"debug","message":"Activating GPIOs","application":"PirriGo","commonWire":21,"gpio":16,"durationSeconds":300,"time":""},{"level":"debug","message":"Deactivating GPIOs","application":"PirriGo","commonWire":21,"gpio":16,"durationSeconds":0,"time":""},{"level":"debug","message":"Task execution complete for station","application":"PirriGo","stationID":9,"time":""}]}
// status/run
// {"IsIdle":true,"IsManual":false,"StartTime":"0001-01-01T00:00:00Z","Duration":0,"ScheduleID":0,"StationID":0,"Cancel":false}

// gpio/available
// { "gpios": [{"ID":1,"GPIO":0,"Notes":"","Common":false},{"ID":2,"GPIO":1,"Notes":"","Common":false},{"ID":4,"GPIO":3,"Notes":"","Common":false},{"ID":8,"GPIO":7,"Notes":"","Common":false},{"ID":9,"GPIO":8,"Notes":"","Common":false},{"ID":10,"GPIO":9,"Notes":"","Common":false},{"ID":11,"GPIO":10,"Notes":"","Common":false},{"ID":12,"GPIO":11,"Notes":"","Common":false},{"ID":15,"GPIO":14,"Notes":"","Common":false},{"ID":16,"GPIO":15,"Notes":"","Common":false},{"ID":18,"GPIO":17,"Notes":"","Common":false},{"ID":20,"GPIO":19,"Notes":"","Common":false},{"ID":28,"GPIO":27,"Notes":"","Common":false}]}

// schedule/all
// { "stationSchedules": [{"ID":1,"StartDate":"2017-03-10T17:08:40Z","EndDate":"2027-03-10T08:00:00Z","Sunday":true,"Monday":false,"Tuesday":true,"Wednesday":false,"Thursday":true,"Friday":false,"Saturday":true,"StationID":4,"StartTime":100,"Duration":1800,"Repeating":false},{"ID":2,"StartDate":"2017-03-11T00:18:22Z","EndDate":"2027-03-10T08:00:00Z","Sunday":true,"Monday":false,"Tuesday":true,"Wednesday":false,"Thursday":true,"Friday":false,"Saturday":false,"StationID":5,"StartTime":100,"Duration":1800,"Repeating":false},{"ID":3,"StartDate":"2017-03-11T00:18:55Z","EndDate":"2027-03-10T08:00:00Z","Sunday":true,"Monday":false,"Tuesday":true,"Wednesday":false,"Thursday":true,"Friday":false,"Saturday":true,"StationID":6,"StartTime":100,"Duration":1800,"Repeating":false},{"ID":4,"StartDate":"2017-03-11T08:00:25Z","EndDate":"2027-03-11T08:00:00Z","Sunday":true,"Monday":false,"Tuesday":true,"Wednesday":false,"Thursday":true,"Friday":false,"Saturday":false,"StationID":7,"StartTime":100,"Duration":1800,"Repeating":false},{"ID":5,"StartDate":"2017-03-11T08:01:01Z","EndDate":"2027-03-11T08:00:00Z","Sunday":true,"Monday":false,"Tuesday":false,"Wednesday":true,"Thursday":false,"Friday":true,"Saturday":false,"StationID":9,"StartTime":100,"Duration":300,"Repeating":false},{"ID":6,"StartDate":"2017-03-11T08:01:23Z","EndDate":"2027-03-11T08:00:00Z","Sunday":true,"Monday":false,"Tuesday":false,"Wednesday":true,"Thursday":false,"Friday":true,"Saturday":false,"StationID":10,"StartTime":100,"Duration":300,"Repeating":false},{"ID":7,"StartDate":"2017-03-11T08:01:52Z","EndDate":"2027-03-11T08:00:00Z","Sunday":true,"Monday":false,"Tuesday":true,"Wednesday":false,"Thursday":true,"Friday":false,"Saturday":false,"StationID":17,"StartTime":100,"Duration":300,"Repeating":false},{"ID":9,"StartDate":"2017-03-11T08:02:38Z","EndDate":"2027-03-11T08:00:00Z","Sunday":true,"Monday":false,"Tuesday":true,"Wednesday":false,"Thursday":false,"Friday":true,"Saturday":false,"StationID":16,"StartTime":100,"Duration":300,"Repeating":false},{"ID":10,"StartDate":"2017-03-12T17:33:13Z","EndDate":"2027-03-12T08:00:00Z","Sunday":true,"Monday":false,"Tuesday":true,"Wednesday":false,"Thursday":false,"Friday":true,"Saturday":false,"StationID":15,"StartTime":1900,"Duration":300,"Repeating":false},{"ID":100,"StartDate":"2020-01-01T01:01:01.000000001Z","EndDate":"2150-01-01T01:01:01.000000001Z","Sunday":false,"Monday":false,"Tuesday":false,"Wednesday":false,"Thursday":false,"Friday":false,"Saturday":false,"StationID":100,"StartTime":0,"Duration":10,"Repeating":false},{"ID":104,"StartDate":"2017-03-12T17:33:13Z","EndDate":"2031-05-02T07:00:00Z","Sunday":true,"Monday":false,"Tuesday":true,"Wednesday":false,"Thursday":false,"Friday":true,"Saturday":false,"StationID":14,"StartTime":1930,"Duration":300,"Repeating":false},{"ID":105,"StartDate":"2017-03-12T17:33:13Z","EndDate":"2031-05-02T07:00:00Z","Sunday":true,"Monday":false,"Tuesday":true,"Wednesday":false,"Thursday":false,"Friday":true,"Saturday":false,"StationID":13,"StartTime":2000,"Duration":300,"Repeating":false}]}

// status/cancel
// {"IsIdle":false,"IsManual":false,"StartTime":"2021-11-30T19:30:22.653689271-08:00","Duration":300,"ScheduleID":0,"StationID":14,"Cancel":false}

// nodes/usage
// { "waterUsage": []}

// stats/1 (report chart thing)
// {"ReportType":1,"Labels":[],"Series":["Unscheduled","Scheduled"],"Data":[[],[]]}

// history?station=undefined&earliest=-168
// { "history": [{"ID":5879,"StationID":14,"ScheduleID":104,"Duration":300,"StartTime":"2021-11-30T19:30:22.622081323-08:00"},{"ID":5878,"StationID":15,"ScheduleID":10,"Duration":300,"StartTime":"2021-11-30T19:00:52.297780761-08:00"},{"ID":5877,"StationID":4,"ScheduleID":1,"Duration":1800,"StartTime":"2021-11-30T02:41:04.35674647-08:00"},{"ID":5876,"StationID":5,"ScheduleID":2,"Duration":1800,"StartTime":"2021-11-30T02:11:02.993172824-08:00"},{"ID":5875,"StationID":6,"ScheduleID":3,"Duration":1800,"StartTime":"2021-11-30T01:41:01.595348591-08:00"},{"ID":5874,"StationID":7,"ScheduleID":4,"Duration":1800,"StartTime":"2021-11-30T01:11:00.290915714-08:00"},{"ID":5873,"StationID":17,"ScheduleID":7,"Duration":300,"StartTime":"2021-11-30T01:05:59.150109505-08:00"},{"ID":5872,"StationID":16,"ScheduleID":9,"Duration":300,"StartTime":"2021-11-30T01:00:58.017370238-08:00"},{"ID":5871,"StationID":4,"ScheduleID":0,"Duration":900,"StartTime":"2021-11-29T13:28:27.457081357-08:00"},{"ID":5870,"StationID":4,"ScheduleID":0,"Duration":3600,"StartTime":"2021-11-29T12:36:36.010888224-08:00"},{"ID":5869,"StationID":13,"ScheduleID":105,"Duration":300,"StartTime":"2021-11-28T20:00:03.260215225-08:00"},{"ID":5868,"StationID":14,"ScheduleID":104,"Duration":300,"StartTime":"2021-11-28T19:30:11.639758838-08:00"},{"ID":5867,"StationID":15,"ScheduleID":10,"Duration":300,"StartTime":"2021-11-28T19:00:41.316314225-08:00"},{"ID":5866,"StationID":4,"ScheduleID":1,"Duration":1800,"StartTime":"2021-11-28T02:51:04.808628856-08:00"},{"ID":5865,"StationID":6,"ScheduleID":3,"Duration":1800,"StartTime":"2021-11-28T02:21:03.467760062-08:00"},{"ID":5864,"StationID":7,"ScheduleID":4,"Duration":1800,"StartTime":"2021-11-28T01:51:02.115574392-08:00"},{"ID":5863,"StationID":9,"ScheduleID":5,"Duration":300,"StartTime":"2021-11-28T01:46:00.98894247-08:00"},{"ID":5862,"StationID":10,"ScheduleID":6,"Duration":300,"StartTime":"2021-11-28T01:40:59.865621051-08:00"},{"ID":5861,"StationID":17,"ScheduleID":7,"Duration":300,"StartTime":"2021-11-28T01:35:58.725250242-08:00"},{"ID":5860,"StationID":16,"ScheduleID":9,"Duration":300,"StartTime":"2021-11-28T01:30:57.590024765-08:00"},{"ID":5859,"StationID":5,"ScheduleID":2,"Duration":1800,"StartTime":"2021-11-28T01:00:56.260172402-08:00"},{"ID":5858,"StationID":4,"ScheduleID":1,"Duration":1800,"StartTime":"2021-11-27T01:30:58.049628921-08:00"},{"ID":5857,"StationID":6,"ScheduleID":3,"Duration":1800,"StartTime":"2021-11-27T01:00:56.71600813-08:00"},{"ID":5856,"StationID":13,"ScheduleID":105,"Duration":300,"StartTime":"2021-11-26T20:00:02.930829987-08:00"},{"ID":5855,"StationID":14,"ScheduleID":104,"Duration":300,"StartTime":"2021-11-26T19:30:11.80334068-08:00"},{"ID":5854,"StationID":15,"ScheduleID":10,"Duration":300,"StartTime":"2021-11-26T19:00:41.504825026-08:00"},{"ID":5853,"StationID":9,"ScheduleID":5,"Duration":300,"StartTime":"2021-11-26T01:10:58.882579856-08:00"},{"ID":5852,"StationID":10,"ScheduleID":6,"Duration":300,"StartTime":"2021-11-26T01:05:57.700383144-08:00"},{"ID":5851,"StationID":16,"ScheduleID":9,"Duration":300,"StartTime":"2021-11-26T01:00:56.551446569-08:00"},{"ID":5850,"StationID":4,"ScheduleID":0,"Duration":3600,"StartTime":"2021-11-25T12:11:24.601626308-08:00"},{"ID":5849,"StationID":4,"ScheduleID":1,"Duration":300,"StartTime":"2021-11-25T01:21:02.486842521-08:00"},{"ID":5848,"StationID":5,"ScheduleID":2,"Duration":300,"StartTime":"2021-11-25T01:16:01.299771288-08:00"},{"ID":5847,"StationID":6,"ScheduleID":3,"Duration":300,"StartTime":"2021-11-25T01:11:00.127758201-08:00"},{"ID":5846,"StationID":7,"ScheduleID":4,"Duration":300,"StartTime":"2021-11-25T01:05:58.948692965-08:00"},{"ID":5845,"StationID":17,"ScheduleID":7,"Duration":300,"StartTime":"2021-11-25T01:00:57.357728884-08:00"},{"ID":5844,"StationID":9,"ScheduleID":5,"Duration":300,"StartTime":"2021-11-24T01:05:58.160022143-08:00"},{"ID":5843,"StationID":10,"ScheduleID":6,"Duration":300,"StartTime":"2021-11-24T01:00:57.05030406-08:00"},{"ID":5842,"StationID":13,"ScheduleID":105,"Duration":300,"StartTime":"2021-11-23T20:00:03.103985343-08:00"},{"ID":5841,"StationID":14,"ScheduleID":104,"Duration":300,"StartTime":"2021-11-23T19:30:12.286414733-08:00"},{"ID":5840,"StationID":15,"ScheduleID":10,"Duration":300,"StartTime":"2021-11-23T19:00:42.04122769-08:00"},{"ID":5839,"StationID":4,"ScheduleID":1,"Duration":300,"StartTime":"2021-11-23T01:26:02.538112495-08:00"},{"ID":5838,"StationID":5,"ScheduleID":2,"Duration":300,"StartTime":"2021-11-23T01:21:01.410563861-08:00"},{"ID":5837,"StationID":6,"ScheduleID":3,"Duration":300,"StartTime":"2021-11-23T01:16:00.301067944-08:00"},{"ID":5836,"StationID":7,"ScheduleID":4,"Duration":300,"StartTime":"2021-11-23T01:10:59.185611668-08:00"},{"ID":5835,"StationID":17,"ScheduleID":7,"Duration":300,"StartTime":"2021-11-23T01:05:58.066601978-08:00"},{"ID":5834,"StationID":16,"ScheduleID":9,"Duration":300,"StartTime":"2021-11-23T01:00:56.915246371-08:00"},{"ID":5833,"StationID":13,"ScheduleID":105,"Duration":300,"StartTime":"2021-11-21T20:00:02.976556335-08:00"},{"ID":5832,"StationID":14,"ScheduleID":104,"Duration":300,"StartTime":"2021-11-21T19:30:12.142865975-08:00"},{"ID":5831,"StationID":15,"ScheduleID":10,"Duration":300,"StartTime":"2021-11-21T19:00:42.863839759-08:00"},{"ID":5830,"StationID":4,"ScheduleID":1,"Duration":300,"StartTime":"2021-11-21T01:36:05.389854647-08:00"},{"ID":5829,"StationID":5,"ScheduleID":2,"Duration":300,"StartTime":"2021-11-21T01:31:04.223785882-08:00"},{"ID":5828,"StationID":6,"ScheduleID":3,"Duration":300,"StartTime":"2021-11-21T01:26:03.077590886-08:00"},{"ID":5827,"StationID":7,"ScheduleID":4,"Duration":300,"StartTime":"2021-11-21T01:21:01.921018802-08:00"},{"ID":5826,"StationID":9,"ScheduleID":5,"Duration":300,"StartTime":"2021-11-21T01:16:00.772086589-08:00"},{"ID":5825,"StationID":10,"ScheduleID":6,"Duration":300,"StartTime":"2021-11-21T01:10:59.631880624-08:00"},{"ID":5824,"StationID":16,"ScheduleID":9,"Duration":300,"StartTime":"2021-11-21T01:05:58.476126897-08:00"},{"ID":5823,"StationID":17,"ScheduleID":7,"Duration":300,"StartTime":"2021-11-21T01:00:57.32827719-08:00"},{"ID":5822,"StationID":4,"ScheduleID":1,"Duration":300,"StartTime":"2021-11-20T01:05:58.095372042-08:00"},{"ID":5821,"StationID":6,"ScheduleID":3,"Duration":300,"StartTime":"2021-11-20T01:00:56.966640587-08:00"},{"ID":5820,"StationID":13,"ScheduleID":105,"Duration":300,"StartTime":"2021-11-19T20:00:03.116894117-08:00"},{"ID":5819,"StationID":14,"ScheduleID":104,"Duration":300,"StartTime":"2021-11-19T19:30:11.665353996-08:00"},{"ID":5818,"StationID":15,"ScheduleID":10,"Duration":300,"StartTime":"2021-11-19T19:00:41.406878873-08:00"},{"ID":5817,"StationID":9,"ScheduleID":5,"Duration":300,"StartTime":"2021-11-19T01:10:59.200585101-08:00"},{"ID":5816,"StationID":10,"ScheduleID":6,"Duration":300,"StartTime":"2021-11-19T01:05:58.025135618-08:00"},{"ID":5815,"StationID":16,"ScheduleID":9,"Duration":300,"StartTime":"2021-11-19T01:00:56.877250882-08:00"},{"ID":5814,"StationID":4,"ScheduleID":1,"Duration":300,"StartTime":"2021-11-18T01:21:01.71689501-08:00"},{"ID":5813,"StationID":5,"ScheduleID":2,"Duration":300,"StartTime":"2021-11-18T01:16:00.593544671-08:00"},{"ID":5812,"StationID":6,"ScheduleID":3,"Duration":300,"StartTime":"2021-11-18T01:10:59.464190639-08:00"},{"ID":5811,"StationID":7,"ScheduleID":4,"Duration":300,"StartTime":"2021-11-18T01:05:58.356499623-08:00"},{"ID":5810,"StationID":17,"ScheduleID":7,"Duration":300,"StartTime":"2021-11-18T01:00:57.237178271-08:00"},{"ID":5809,"StationID":9,"ScheduleID":5,"Duration":300,"StartTime":"2021-11-17T01:05:58.591937991-08:00"},{"ID":5808,"StationID":10,"ScheduleID":6,"Duration":300,"StartTime":"2021-11-17T01:00:57.424427994-08:00"},{"ID":5807,"StationID":13,"ScheduleID":105,"Duration":300,"StartTime":"2021-11-16T20:00:02.745563818-08:00"},{"ID":5806,"StationID":14,"ScheduleID":104,"Duration":300,"StartTime":"2021-11-16T19:30:12.259737592-08:00"},{"ID":5805,"StationID":15,"ScheduleID":10,"Duration":300,"StartTime":"2021-11-16T19:00:41.94958043-08:00"},{"ID":5804,"StationID":4,"ScheduleID":1,"Duration":300,"StartTime":"2021-11-16T01:26:02.743855446-08:00"},{"ID":5803,"StationID":5,"ScheduleID":2,"Duration":300,"StartTime":"2021-11-16T01:21:01.594581588-08:00"},{"ID":5802,"StationID":6,"ScheduleID":3,"Duration":300,"StartTime":"2021-11-16T01:16:00.456584612-08:00"},{"ID":5801,"StationID":17,"ScheduleID":7,"Duration":300,"StartTime":"2021-11-16T01:10:59.313246182-08:00"},{"ID":5800,"StationID":16,"ScheduleID":9,"Duration":300,"StartTime":"2021-11-16T01:05:58.153888787-08:00"},{"ID":5799,"StationID":7,"ScheduleID":4,"Duration":300,"StartTime":"2021-11-16T01:00:56.991823024-08:00"},{"ID":5798,"StationID":13,"ScheduleID":105,"Duration":300,"StartTime":"2021-11-14T20:00:03.051526802-08:00"},{"ID":5797,"StationID":14,"ScheduleID":104,"Duration":300,"StartTime":"2021-11-14T19:30:12.576804061-08:00"},{"ID":5796,"StationID":15,"ScheduleID":10,"Duration":300,"StartTime":"2021-11-14T19:00:42.238838364-08:00"},{"ID":5795,"StationID":5,"ScheduleID":2,"Duration":300,"StartTime":"2021-11-14T01:36:05.329784754-08:00"},{"ID":5794,"StationID":6,"ScheduleID":3,"Duration":300,"StartTime":"2021-11-14T01:31:04.17558821-08:00"},{"ID":5793,"StationID":7,"ScheduleID":4,"Duration":300,"StartTime":"2021-11-14T01:26:03.005262164-08:00"},{"ID":5792,"StationID":9,"ScheduleID":5,"Duration":300,"StartTime":"2021-11-14T01:21:01.826616266-08:00"},{"ID":5791,"StationID":10,"ScheduleID":6,"Duration":300,"StartTime":"2021-11-14T01:16:00.667293079-08:00"},{"ID":5790,"StationID":17,"ScheduleID":7,"Duration":300,"StartTime":"2021-11-14T01:10:59.523484239-08:00"},{"ID":5789,"StationID":16,"ScheduleID":9,"Duration":300,"StartTime":"2021-11-14T01:05:58.360234348-08:00"},{"ID":5788,"StationID":4,"ScheduleID":1,"Duration":300,"StartTime":"2021-11-14T01:00:57.116268382-08:00"},{"ID":5787,"StationID":4,"ScheduleID":1,"Duration":300,"StartTime":"2021-11-13T01:05:57.668341902-08:00"},{"ID":5786,"StationID":6,"ScheduleID":3,"Duration":300,"StartTime":"2021-11-13T01:00:56.505171881-08:00"},{"ID":5785,"StationID":13,"ScheduleID":105,"Duration":300,"StartTime":"2021-11-12T20:00:02.662437079-08:00"},{"ID":5784,"StationID":14,"ScheduleID":104,"Duration":300,"StartTime":"2021-11-12T19:30:11.883999452-08:00"},{"ID":5783,"StationID":15,"ScheduleID":10,"Duration":300,"StartTime":"2021-11-12T19:00:42.194278921-08:00"},{"ID":5782,"StationID":9,"ScheduleID":5,"Duration":300,"StartTime":"2021-11-12T01:10:59.136819274-08:00"},{"ID":5781,"StationID":16,"ScheduleID":9,"Duration":300,"StartTime":"2021-11-12T01:05:57.980524272-08:00"},{"ID":5780,"StationID":10,"ScheduleID":6,"Duration":300,"StartTime":"2021-11-12T01:00:56.563206582-08:00"}]}