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


  // "/stats/1": statsActivityByStation,
  // "/stats/2": statsActivityByDayOfWeek,
  // "/stats/3": statsActivityPerStationByDOW,
  // "/stats/4": statsStationActivity,
  loadChartByID(chart: number, startDate: number, endDate: number): Observable<ChartData> {
    const uri = `${this._globals.uriStem}/stats/${chart}?startDate=${startDate}&endDate=${endDate}`;
    return this._http.get<ChartData>(uri)
  }
}