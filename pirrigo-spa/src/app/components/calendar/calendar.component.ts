import { Component, OnInit, Inject } from '@angular/core';
import { CalendarEvent } from 'calendar-utils';
import { ApiClientService } from 'src/app/services/apiclient.service';
import { StationSchedule } from 'src/app/structs/station-schedule';
import * as moment from 'moment';
import { MatDialog, MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { MatCardModule } from '@angular/material/card';


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

  eventClicked({ event }: { event: CalendarEvent }): void {
    this.editingSchedule = JSON.parse(event.title.split(" | ")[1])
    let mstart = moment(event.start)
    let mend = moment(event.end)
    this.editingSchedule.StartDate = new Date();
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
            "title": `Station ${event.StationID} for ${event.Duration / 60} minutes<br/><br/><br/><br/> | ${JSON.stringify(event)}
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
      // width: '80%',
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
  selector: 'dialog-overview-example-dialog',
  templateUrl: `./dialog-overview-example-dialog.html`,
})
export class EditScheduleDialog {

  tempStartTime: string;

  constructor(
    private _api: ApiClientService,
    public dialogRef: MatDialogRef<EditScheduleDialog>,
    @Inject(MAT_DIALOG_DATA) public data: StationSchedule,
  ) { }

  setStartTime(event: any): void {
    console.log(event)
    this.tempStartTime = event
  }

  setDuration(event: any): void {
    this.data.Duration = event.value * 60
  }

  formatSliderLabel(value: number) {
    return `${value}m`;
  }

  submitScheduleChange(schedule: StationSchedule): void {
    console.log(schedule)
  }


  onNoClick(): void {
    this.dialogRef.close();
    console.log(this.data)
  }

}