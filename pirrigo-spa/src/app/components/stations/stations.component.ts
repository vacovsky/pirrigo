import { Component, OnInit } from '@angular/core';
import { ApiClientService } from 'src/app/services/apiclient.service';
import { Station, StationStatus } from 'src/app/structs/station';
import { MwlGaugeObj } from 'src/app/structs/mwl-gauge-obj';
import { GlobalsService } from 'src/app/services/globals.service';



@Component({
  selector: 'app-stations',
  templateUrl: './stations.component.html',
  styleUrls: ['./stations.component.css']
})
export class StationsComponent implements OnInit {

  panelOpenState = false;
  status: StationStatus;
  runningGauge: MwlGaugeObj;
  stations: Station[];

  constructor(
    private _api: ApiClientService,
    private _globals: GlobalsService
  ) { }

  ngOnInit(): void {
    this.loadStations().then(() =>
      this.loadStationRunStatus()
    )
    setInterval(() => {
      this.loadStationRunStatus();
    }, this._globals.statusRefreshRateMs);
  }

  async loadStations() {
    this._api.getAllStations().subscribe((data) => {
      this.stations = data.stations
      console.log(this.stations)
    })
  }

  loadStationRunStatus() {
    this._api.getStationStatus().subscribe(data => {
      this.status = data
      console.log(this.status)
      this.runningGauge = this.gaugeFactory(
        2900,
        0,
        3600,
        data.Duration
      )
    })
  }

  findDateDiffForGaugeInSeconds(date: Date, duration: number): number {
    return 150

    // var duration = moment.duration(end.diff(startTime));
    // var hours = duration.asHours();
  }

  gaugeFactory(
    value: number,
    min: number,
    max: number,
    label: any = ((inp: number) => { return inp }),
    color: any = ((inp: number) => { return inp })
  ): MwlGaugeObj {
    let g = new MwlGaugeObj()
    g.Value = value
    g.Max = max
    g.Min = min
    g.Animated = true
    g.AnimationDuration = 5
    g.DialEndAngle = -180
    g.DialStartAngle = 180
    g.Label = label
    g.Color = color

    return g
  }

}



// async ngOnInit() {
//   let endDate = moment().unix();
//   let startDate = moment().subtract(4, 'd').unix();
//   await this._api.loadHistoryChart(startDate, endDate).subscribe((data) => {
//     this.curTemp = Number(data.Data[0][data.Data[0].length - 1])
//     this.curHum = Number(data.Data[1][data.Data[0].length - 1])
//     let na = data.Data[0].map(Number)
//     this.maxTemp = Math.max(...na)
//     this.minTemp = Math.min(...na)

//     this.gaugeCurTemp = this.gaugeFactory(
//       this.curTemp,
//       25,
//       125,
//       (d) => `${d}°F`,
//       (d) => {
//         if (d < 35) {
//           return "#000aff";
//         } else if (d >= 35 && d < 62) {
//           return "#0a8a8f";
//         } else if (d >= 62 && d < 90) {
//           return "#0a8f3b";
//         } else {
//           return "#0a8a8f";
//         }
//       }
//     )

//     this.gaugeCurHum = this.gaugeFactory(
//       this.curHum,
//       0,
//       100,
//       (d) => `${d}%rh`,
//       (d) => {
//         if (d < 35) {
//           return "#DD9155";
//         } else if (d >= 35 && d < 75) {
//           return "#088B58";
//         } else if (d > 75) {
//           return "#0898C6";
//         } else {
//           return "#000";
//         }
//       }
//     )

//   }, () => {
//     console.log("[RecentDataComponent] error making API call to [_api.loadHistoryChart]")
//   })

//   let endDateLy = moment().subtract(1, 'y').unix();
//   let startDateLy = moment().subtract(1, 'y').subtract(1, 'd').unix();
//   await this._api.loadHistoryChart(startDateLy, endDateLy).subscribe((data) => {
//     let na = data.Data[0].map(Number)
//     this.maxTempLy = Math.max(...na)
//     this.minTempLy = Math.min(...na)
//   }, () => {
//     console.log("[RecentDataComponent] error making API call to [_api.loadHistoryChart]")
//   })
// }

// gaugeFactory(
//   value: number,
//   min: number,
//   max: number,
//   label: any = ((inp: number) => { return inp }),
//   color: any = ((inp: number) => { return inp })
// ): MwlGaugeObj {
//   let g = new MwlGaugeObj()
//   g.Value = value
//   g.Max = max
//   g.Min = min
//   g.Animated = true
//   g.AnimationDuration = 1
//   g.DialEndAngle = 190
//   g.DialStartAngle = -10
//   g.Label = label
//   g.Color = color

//   return g
// }
// }
