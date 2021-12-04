import { Component, OnInit } from '@angular/core';
import { CalendarEvent } from 'calendar-utils';
import { ApiClientService } from 'src/app/services/apiclient.service';
import { StationSchedule } from 'src/app/structs/station-schedule';
import * as moment from 'moment';
import { SchedulerLike } from 'rxjs';
// import { endianness } from 'os';

@Component({
  selector: 'app-calendar',
  templateUrl: './calendar.component.html',
  styleUrls: ['./calendar.component.css']
})
export class CalendarComponent implements OnInit {

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
    private _api: ApiClientService
  ) { }

  ngOnInit(): void {
    this.viewDate = new Date();
    this.loadSchedule()
  }

  loadSchedule() {
    this._api.getStationSchedules().subscribe(
      (data) => {
        this.events = this.convertScheduleToCalendarEvents(data.stationSchedules)
        console.log(this.events)
      }
    )
  }


  loadstuff(event: StationSchedule): any {
    let dayIsActiveHash = {
      "Sunday": event.Sunday,
      "Monday": event.Monday,
      "Tuesday": event.Tuesday,
      "Wednesday": event.Wednesday,
      "Thursday": event.Thursday,
      "Friday": event.Friday,
      "Saturday": event.Saturday,
    }
    return dayIsActiveHash
  }

  convertScheduleToCalendarEvents(schedule: StationSchedule[]): CalendarEvent[] {
    let events: CalendarEvent[] = [];
    for (let i = -8; i < this.DOW.length; i++) {
      let d: moment.Moment = moment(new Date()).add(i, "d")
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
          let start: moment.Moment = d
            .add(
              (Math.floor(event.StartTime / 10), "h"))
            .add(
              (event.StartTime - (Math.floor(event.StartTime / 60)), "m")
            )
          // console.log(start.format())
          let end: Date = moment(start.toDate()).add((event.Duration / 60), "m").toDate()
          // console.log(moment(end).format())

          let newEvent = {
            "id": event.ID,
            "start": start.toDate(),
            "end": end,
            "title": "Station Run: " + event.StationID,
            "color": this.colors.red,
            // "actions": EventAction[],
            "allDay": false,
            // "cssClass": string,
            "resizable": {
              "beforeStart": false,
              "afterEnd": false,
            },
            "draggable": false,
            // "meta": MetaType,
          }
          events.push(newEvent)
          // console.log(newEvent)
        }
      }
    }
    return events
  }


  getDOWForMoment(date: Date): string | undefined {
    return this.DOW.find(day => day == moment(date).format('dddd'))
  }


  // {"ID":1,
  // "StartDate":"2017-03-10T17:08:40Z",
  // "EndDate":"2027-03-10T08:00:00Z",
  // "Sunday":true,  "Monday":false,"Tuesday":true,"Wednesday":false,"Thursday":true,"Friday":false,"Saturday":true,
  // "StationID":4,
  // "StartTime":100,
  // "Duration":1800,
  // "Repeating":false}
  DOW: string[] = [
    "Sunday",
    "Monday",
    "Tuesday",
    "Wednesday",
    "Thursday",
    "Friday",
    "Saturday"
  ]
}