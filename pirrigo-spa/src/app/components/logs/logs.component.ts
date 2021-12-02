import { Component, OnInit } from '@angular/core';
import { ApiClientService } from 'src/app/services/apiclient.service';
import { StationLogs } from 'src/app/structs/station-logs';


@Component({
  selector: 'app-logs',
  templateUrl: './logs.component.html',
  styleUrls: ['./logs.component.css']
})
export class LogsComponent implements OnInit {

  logs: StationLogs[];

  displayedColumns: string[] = [
    'time',
    'level',
    'message',
    'application',
    'interval'
  ];

  constructor(
    private _api: ApiClientService
  ) { }

  ngOnInit(): void {
    this._api.getStationLogs().subscribe(data => {
      this.logs = data.logs
    })
  }

}
