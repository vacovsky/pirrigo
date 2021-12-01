import { Component, OnInit } from '@angular/core';
import { ApiClientService } from 'src/app/services/apiclient.service';
import { Station } from 'src/app/structs/station';
@Component({
  selector: 'app-stations',
  templateUrl: './stations.component.html',
  styleUrls: ['./stations.component.css']
})
export class StationsComponent implements OnInit {

  panelOpenState = false;

  stations: Station[];

  constructor(
    private _api: ApiClientService
  ) { }

  ngOnInit(): void {
    this.loadStations()
  }

  loadStations() {
    this._api.getAllStations().subscribe((data) => {
      this.stations = data.stations
    })
  }
}
