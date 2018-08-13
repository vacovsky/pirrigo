import { Injectable } from "@angular/core";
import { HttpHeaders } from "@angular/common/http";
import { Chart } from "chart.js";

@Injectable()
export class Globals {
    // General

    public authenticated = true;
    public convertDate(date: Date): string {
        return Math.round(date.getTime() / 1000).toString();
    }
}
