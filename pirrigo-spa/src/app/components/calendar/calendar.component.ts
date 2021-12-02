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
  }

  loadSchedule() {
    this._api.getStationSchedules().subscribe(
      (data) => {
        this.events = this.convertScheduleToCalendarEvents(data.schedule)
      }
    )
  }

  convertScheduleToCalendarEvents(schedule: StationSchedule[]): CalendarEvent[] {
    let events: CalendarEvent[] = [];
    for (let event of schedule) {
      let newEvent = {
        "id": event.ID,
        "start": moment(event.StartDate).add((event.StartTime / 60), 'h').toDate(),
        "end": moment(event.StartDate).add((event.StartTime / 60) + (event.Duration / 60), 'h').toDate(),
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
    }
    return events
  }
}