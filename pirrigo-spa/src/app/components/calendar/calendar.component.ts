import { Component, OnInit, Inject, AfterViewInit } from '@angular/core';
import { CalendarEvent } from 'calendar-utils';
import { ApiClientService } from 'src/app/services/apiclient.service';
import { StationSchedule } from 'src/app/structs/station-schedule';
import { Station } from 'src/app/structs/station';
import * as moment from 'moment';
import { MatDialog, MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';

@Component({
  selector: 'app-calendar',
  templateUrl: './calendar.component.html',
  styleUrls: ['./calendar.component.css']
})
export class CalendarComponent implements OnInit {

  editingSchedule: StationSchedule;
  viewDate: Date;
  events: CalendarEvent[];
  colors: any = {
    red: {
      primary: '#ad2121',
      secondary: '#FAE3E3',
    },
    blue: {
      primary: '#1e90ff',
      secondary: '#D1E8FF',
    },
    yellow: {
      primary: '#e3bc08',
      secondary: '#FDF1BA',
    },
  };

  constructor(
    private _api: ApiClientService,
    public dialog: MatDialog
  ) { }

  ngOnInit(): void {
    this.viewDate = new Date();
    this.loadSchedule()
  }

  loadSchedule() {
    this._api.getStationSchedules().subscribe(
      (data) => {
        this.events = this.convertScheduleToCalendarEvents(data.stationSchedules)
        // console.log(this.events)
      }
    )
  }


  openNewScheduleDialog(): void {
    let sch = new StationSchedule()
    sch.StartDate = moment(new Date().setHours(0, 0, 0, 0)).toDate();
    sch.EndDate = moment().add(15, "y").toDate();
    this.openDialog(sch)
  }

  eventClicked({ event }: { event: CalendarEvent }): void {
    this.editingSchedule = JSON.parse(event.title.split(" | ")[1])
    this.editingSchedule.StartDate = moment(new Date().setHours(0, 0, 0, 0)).toDate();
    this.editingSchedule.EndDate = moment().add(15, "y").toDate();
    this.openDialog(this.editingSchedule)
  }

  eventTimesChanged({ event }: { event: CalendarEvent }): void {
    console.log(event);
  }

  convertScheduleToCalendarEvents(schedule: StationSchedule[]): CalendarEvent[] {
    let events: CalendarEvent[] = [];
    for (let i = -8; i < this.DOW.length; i++) {
      let d: moment.Moment = moment(new Date().setHours(0, 0, 0, 0)).add(i, "d")
      for (let event of schedule) {
        if (
          (d.format('dddd') == "Sunday" && event.Sunday)
          || (d.format('dddd') == "Monday" && event.Monday)
          || (d.format('dddd') == "Tuesday" && event.Tuesday)
          || (d.format('dddd') == "Wednesday" && event.Wednesday)
          || (d.format('dddd') == "Thursday" && event.Thursday)
          || (d.format('dddd') == "Friday" && event.Friday)
          || (d.format('dddd') == "Saturday" && event.Saturday)
        ) {
          let hm = this.convertMilIntTo12h(event.StartTime)
          let start: moment.Moment = moment(d).add(hm[0], 'h').add(hm[1], "m")
          let end: Date = moment(start.toDate()).add(event.Duration, "s").toDate()
          let newEvent = {
            "id": event.ID,
            "start": start.toDate(),
            "end": end,
            "title": `Zone ${event.StationID} for ${event.Duration / 60} minutes<br/><br/><br/><br/> | ${JSON.stringify(event)}
            `,
            "color": this.colors.blue,
            "allDay": false,
            "resizable": {
              "beforeStart": false,
              "afterEnd": false,
            },
            "draggable": false,
          }
          events.push(newEvent)
        }
      }
    }
    return events
  }

  getDOWForMoment(date: Date): string | undefined {
    return this.DOW.find(day => day == moment(date).format('dddd'))
  }

  convertMilIntTo12h(mt: number | string): string[] {
    mt = mt.toString()
    return [mt.substring(0, mt.length - 2), mt.substring(mt.length - 2, mt.length)];
  }

  DOW: string[] = [
    "Sunday",
    "Monday",
    "Tuesday",
    "Wednesday",
    "Thursday",
    "Friday",
    "Saturday"
  ]


  openDialog(es: StationSchedule): void {
    const dialogRef = this.dialog.open(EditScheduleDialog, {
      data: es
    });

    dialogRef.afterClosed().subscribe(result => {
      this.ngOnInit()
    });
  }

  submitScheduleChange(schedule: StationSchedule): void {
    this._api.postStationScheduleChange(schedule).subscribe((data) => {
      this.events = this.convertScheduleToCalendarEvents(data.stationSchedules)
    })
  }
}


@Component({
  styleUrls: ['./calendar.component.css'],
  selector: 'dialog-scheduleform',
  templateUrl: `./dialog-scheduleform.html`,
})
export class EditScheduleDialog implements OnInit, AfterViewInit {

  tempStartTime: string;
  tempStation: Station;
  tempStationsList: Station[];

  constructor(
    private _api: ApiClientService,
    public dialogRef: MatDialogRef<EditScheduleDialog>,
    @Inject(MAT_DIALOG_DATA) public data: StationSchedule,
  ) { }


  ngOnInit() {
    if (this.data.StationID == undefined) {
      this._api.getAllStations().subscribe(stationdata => {
        this.tempStationsList = stationdata.stations
        console.log(this.tempStationsList)
      })

    } else {
      this._api.getStationForScheduleEdit(this.data.StationID).subscribe(stationdata => {
        this.tempStation = stationdata
      })
    }
  }

  ngAfterViewInit() {
    // this.tempStartTime = this.data.StartTime.toString()
  }

  setStationId(event: any): void {
    this.data.StationID = event.value
    // console.log(event.value)
  }

  setStartTime(event: string): void {
    this.tempStartTime = event
    this.data.StartTime = Number(event.replace(":", ""))
  }

  setDuration(event: any): void {
    this.data.Duration = event.value * 60
  }

  formatSliderLabel(value: number) {
    return `${value}m`;
  }

  deleteScheduleItem(id: number): void {
    this._api.deleteStationScheduleItem(id).subscribe((d) => {
      this.dialogRef.close();
    })
  }

  submitScheduleChange(schedule: StationSchedule): void {
    console.log(schedule)
    this._api.postStationScheduleChange(schedule).subscribe(() => {
      this.dialogRef.close();
      console.log(this.data)
    })
  }

  closeEditSchedule() {
    this.dialogRef.close();
  }

  onNoClick(): void {
    this.dialogRef.close();
    console.log(this.data)
  }

}