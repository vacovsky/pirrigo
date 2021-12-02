import { Component, OnInit } from '@angular/core';
import { ApiClientService } from 'src/app/services/apiclient.service';
import { StationHistory } from 'src/app/structs/station-history';
import * as moment from 'moment';

@Component({
  selector: 'app-history',
  templateUrl: './history.component.html',
  styleUrls: ['./history.component.css']
})
export class HistoryComponent implements OnInit {
  displayedColumns: string[] = ['Station ID', 'Schedule ID', 'Start Time', 'Duration'];
  history: StationHistory[];

  constructor(
    private _api: ApiClientService
  ) { }

  ngOnInit(): void {
    this.loadHistory()
  }

  parseDateTimeForHumans(date: Date): string {
    return `${moment(date).fromNow()} (${date})`
  }

  loadHistory() {
    this._api.getRunHistory(undefined, -168).subscribe(data => {
      this.history = data.history;
    })
  }
}
