import { Injectable } from "@angular/core";
import { Globals } from "./globals.service";
import { HttpClient, HttpHeaders, HttpParams } from "@angular/common/http";
import { map } from "rxjs/operators";

@Injectable({
  providedIn: "root"
})
export class StationService {
  constructor(private _globals: Globals, private _http: HttpClient) {}
  // params: HttpParams
  Load(): any {
    const uriStem = "/station/all";
    const result = this._http
      .get(this._globals.apiHost + uriStem, {
        headers: this._globals.headers,
        // params: params
      })
      .pipe(map(_ => result))
      .subscribe(data => {
        this._globals.zones = data["stations"];
      });
    return result;
  }
}
