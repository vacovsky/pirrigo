import { Component } from "@angular/core";
import { Globals } from "./services/globals.service";

@Component({
    selector: "app-root",
    templateUrl: "./app.component.html",
    styleUrls: ["./app.component.css"]
})
export class AppComponent {
    constructor(private _globals: Globals) {}
}
