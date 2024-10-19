import { Component, Inject, ViewChild, OnInit, AfterViewInit } from '@angular/core';
import { ApiClientService } from 'src/app/services/apiclient.service';
import { MatSliderModule } from '@angular/material/slider';
import { Station, StationStatus, StationProgressBar, StationRunJob, StationRunRequestBody } from 'src/app/structs/station';
import { MwlGaugeObj } from 'src/app/structs/mwl-gauge-obj';
import { GlobalsService } from 'src/app/services/globals.service';
import * as moment from 'moment';
import { MatDialog, MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
// import { CdkTextareaAutosize } from '@angular/cdk/text-field';
import { Gpio } from 'src/app/structs/gpio';
@Component({
  selector: 'app-stations',
  templateUrl: './stations.component.html',
  styleUrls: ['./stations.component.css']
})
export class StationsComponent implements OnInit {

  runRequestRunTime: number = 15;
  panelOpenState = false;
  status: StationStatus;
  runningGauge: MwlGaugeObj;
  stations: Station[];
  stationProgressBar: StationProgressBar;

  runQueue: StationRunJob[] = [];
  runRequest: StationRunRequestBody;

  constructor(
    private _api: ApiClientService,
    private _globals: GlobalsService,
    public dialog: MatDialog
  ) { }

  ngOnInit(): void {
    this.loadStations().then(() => {
      this.loadStationRunStatus().then(() => {
        this.loadStationsRunQueue()
      })
    })
    setInterval(() => {
      this.loadStationRunStatus();
      this.loadStationsRunQueue()
    }, this._globals.statusRefreshRateMs);
  }


  getPercentComplete(input: number): number {
    // console.log(input)
    if (this.status != undefined) {
      return Math.round(((input / 100) * this.status.Duration) / 60)
    }
    return 0
  }

  async loadStations() {
    this._api.getAllStations().subscribe((data) => {
      this.stations = data.stations
      // console.log(this.stations)
    })
  }

  async loadStationsRunQueue() {
    if (this.status != undefined) {
      this._api.getStationRunQueue().subscribe((data) => {
        let totalSeconds: number = 0;
        let currentRunRemainingSec: number = this.findDateDiffInSeconds(this.status.StartTime, this.status.Duration)
        totalSeconds += currentRunRemainingSec
        let qi = 0
        for (let job of data) {
          if (self.status != undefined) {
            let now = moment(new Date())
            job.startTime = now.add(totalSeconds, "s").fromNow()
            totalSeconds += currentRunRemainingSec
            job.queueIndex = qi
            qi++
            // console.log(totalSeconds, job.startTime, totalSeconds)
          }
        }
        this.runQueue = data
        // console.log(this.runQueue)
      })
    }
  }

  async loadStationRunStatus() {
    this._api.getStationStatus().subscribe(data => {
      this.status = data
      let tempPB = new StationProgressBar()
      tempPB.StationID = this.status.StationID
      tempPB.percentComplete = this.findDateDiffPercent(this.status.StartTime, this.status.Duration)
      this.stationProgressBar = tempPB
      // console.log(this.stationProgressBar)
    })
  }

  findDateDiffInSeconds(date: Date, duration: number): number {
    let now = moment(new Date());
    let end = moment(date).add(duration, "s");
    let durationDiff = moment.duration(now.diff(end));
    let sec = durationDiff.asSeconds()
    return -sec
  }

  findDateDiffPercent(date: Date, duration: number): number {
    let now = moment(new Date());
    let end = moment(date).add(duration, "s");
    let durationDiff = moment.duration(now.diff(end));
    let sec = durationDiff.asSeconds()
    return Math.round(100 - ((-sec / this.status.Duration) * 100))
  }

  runStation(station: number, seconds: number) {
    this.runRequest = new StationRunRequestBody()
    this.runRequest.Duration = seconds
    this.runRequest.StationID = station
    this._api.postStationRun(this.runRequest)
  }

  // status/cancel
  cancelRunStation() {
    this._api.cancelActiveStationRun().subscribe(data => {
      this.status = data
      let tempPB = new StationProgressBar()
      tempPB.StationID = this.status.StationID
      tempPB.percentComplete = this.findDateDiffPercent(this.status.StartTime, this.status.Duration)
      this.stationProgressBar = tempPB
      // console.log("Cancelled station run:", data)
    })
  }


  cancelJobInQueue(queueIndex: number) {
    this._api.cancelQueuedJob({ QueueIndex: queueIndex }).subscribe((data) => {
      this.loadStationsRunQueue()
    })
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
    g.DialEndAngle = -179
    g.DialStartAngle = 179
    g.Label = label
    g.Color = color
    return g
  }

  updateSliderValue(event: any): void {
    this.runRequestRunTime = event.value
  }

  formatSliderLabel(value: number) {
    return `${value}m`;
  }


  addNewStation(): void {
    this.ngOnInit()
  }


  openNewStationDialog(): void {
    let sch = new Station()
    this.openDialog(sch)
  }

  stationEditClicked(station: Station): void {
    this.openDialog(station)
  }

  openDialog(es: Station): void {
    const dialogRef = this.dialog.open(EditStationDialog, {
      data: es
    });

    dialogRef.afterClosed().subscribe(() => {
      this.ngOnInit()
    });
  }

}





@Component({
  styleUrls: ['./stations.component.css'],
  selector: 'dialog-scheduleform',
  templateUrl: `./dialog-stationedit.html`,
})
export class EditStationDialog implements OnInit, AfterViewInit {
  // @ViewChild('autosize') autosize: CdkTextareaAutosize;
  tempGPIOs: Gpio[];
  // tempStationsList: Station[];

  constructor(
    private _api: ApiClientService,
    // private _ngZone: NgZone,
    public dialogRef: MatDialogRef<EditStationDialog>,
    @Inject(MAT_DIALOG_DATA) public data: Station,
  ) { }


  ngOnInit() {
    this._api.getGPIOsForStationEdit().subscribe(stationdata => {
      this.tempGPIOs = stationdata.gpios
    })

    if (this.data.GPIO == undefined) {

    } else {

    }
  }

  setGPIO(event: any): void {
    this.data.GPIO = event.value
  }


  ngAfterViewInit() {
  }

  setStationGPIO(event: any): void {
    this.data.GPIO = event.value
    console.log(event.value)
  }

  formatSliderLabel(value: number) {
    return `${value}m`;
  }

  deleteStation(id: number): void {
    this._api.deleteStation(id).subscribe((d) => {
      this.dialogRef.close();
    })
  }

  submitStationEdit(station: Station): void {
    this._api.postStationChange(station).subscribe(() => {
      this.dialogRef.close();
    })
  }

  closeEditStation() {
    this.dialogRef.close();
  }

  onNoClick(): void {
    this.dialogRef.close();
    console.log(this.data)
  }

}