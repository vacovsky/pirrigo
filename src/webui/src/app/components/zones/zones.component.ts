import { Component, AfterViewInit } from "@angular/core";
import { StationService } from "../../services/station.service";
import { Globals } from "../../services/globals.service";

@Component({
  selector: "app-zones",
  templateUrl: "./zones.component.html",
  styleUrls: ["./zones.component.css"]
})
export class ZonesComponent implements AfterViewInit {
  constructor(private _globals: Globals, private _stations: StationService) {}

  ngAfterViewInit() {
    this._stations.Load();
  }
}
