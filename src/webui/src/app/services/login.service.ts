// import { Injectable, OnInit } from "@angular/core";
// import { HttpClient, HttpHeaders, HttpParams } from "@angular/common/http";
// import "rxjs/add/operator/map";
// import { Globals } from "../services/globals";
//
// @Injectable()
// export class LoginService implements OnInit {
// constructor() //
// private _http: HttpClient;
// // private _globals: Globals,
// // private _watrapi: WatrapiService,
// // private _params: ParamsService,
// // private _cookies: CookieService,
// // private _sidenav: SidenavService
// {
// }
//
// ngOnInit() {}
//
// public loginMessage: string;
//
// user: string;
// pass: string;
//
// public login(username: string, password: string): void {
// const header = new HttpHeaders({
// Authorization: this.buildToken(username, password)
// });
// this.makeLoginRequest(header);
// }
//
// makeLoginRequest(header: HttpHeaders): void {
// this._http
// .post(
// this._globals.apiRoot + "auth",
// { observe: "response" },
// { headers: header }
// )
// .map(result => result)
//
// .subscribe(resp => {
// if (resp["status"] == "success") {
// // this._globals.headers = header;
// // this._globals.authenticated = true;
// // this.afterAuth();
// // this._cookies.set("user_key", header.get("Authorization"));
// } else {
// this.loginMessage = "Login Failed";
// this._globals.authenticated = false;
// }
// });
// }
//
// buildToken(username: string, password: string): string {
// return "Bearer " + btoa(username + ":" + password);
// }
//
// public logout() {
// // this._cookies.delete("user_key");
// // window.location.reload(true);
// }
//
// private afterAuth() {
// // this._watrapi.loadDevices();
// // this._watrapi.loadUser(this._params.getParams());
// // this._sidenav.setNav(
// // this._cookies.get("lastTab") != ""
// // ? { name: this._cookies.get("lastTab") }
// // : { name: "Dashboard" }
// );
//
// // this._watrapi.getDeviceReport(this._params.getReportParams(this._globals.currentDevice.Serial));
// }
// }
