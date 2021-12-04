import { Component, OnInit } from '@angular/core';
import { CalendarEvent } from 'calendar-utils';
import { ApiClientService } from 'src/app/services/apiclient.service';
import { StationSchedule } from 'src/app/structs/station-schedule';
import * as moment from 'moment';


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

  eventClicked({ event }: { event: CalendarEvent }): void {
    console.log(event);
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
          let start: moment.Moment = d.add(hm[0], 'h').add(hm[1], "m")
          let end: Date = moment(start.toDate()).add(event.Duration, "s").toDate()

          let newEvent = {
            "id": event.ID,
            "start": start.toDate(),
            "end": end,
            "title": "Station Run: " + event.StationID,
            "color": this.colors.blue,
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


  convertMilIntTo12h(mt: number | string): string[] {
    mt = mt.toString()
    return [mt.substring(0, mt.length - 2), mt.substring(mt.length - 2, mt.length)];
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