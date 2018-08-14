import { Injectable } from "@angular/core";
import { HttpHeaders } from "@angular/common/http";
import { Chart } from "chart.js";

@Injectable()
export class Globals {
    // General

    public authHeader = "Basic ZGVtbzpkZW1v";
    public headers = new HttpHeaders({'Authorization': "Basic ZGVtbzpkZW1v"});
    public zones = [];
    public authenticated = true;
    public currentTab = "Zones";

    public apiHost = "http://localhost:8001";
    public convertDate(date: Date): string {
        return Math.round(date.getTime() / 1000).toString();
    }
}
