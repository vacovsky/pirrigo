import { Injectable } from "@angular/core";
import { CookieService } from "ngx-cookie-service";
import { Globals } from "./globals.service";

@Injectable({
  providedIn: "root"
})
export class SidenavService {
  constructor(private _cookies: CookieService, private _globals: Globals) {}

  setNav(nav: any = { name: "Zones" }): void {
    this._globals.currentTab = nav.name;
    this._cookies.set("lastTab", nav.name);
  }
}
