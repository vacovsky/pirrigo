import {
    Component,
    OnInit,
    ChangeDetectorRef,
    AfterViewInit
} from "@angular/core";
import { Chart } from "chart.js";
import { MediaMatcher } from "@angular/cdk/layout";
import { MatIconRegistry } from "@angular/material/icon";
import { DomSanitizer } from "@angular/platform-browser";
import { Globals } from "../../services/globals.service";

@Component({
    selector: "app-sidenav",
    templateUrl: "./sidenav.component.html",
    styleUrls: ["./sidenav.component.css"]
})
export class SidenavComponent implements AfterViewInit {
    mobileQuery: MediaQueryList;

    navItems: object = [
        { name: "Zones", icon: "equalizer" },
        { name: "Schedule", icon: "warning" },
        { name: "Calendar", icon: "transform" },
        { name: "Water Usage", icon: "opacity" },
        { name: "Weather", icon: "developer_board" },
        { name: "Settings", icon: "settings" }
    ];

    private _mobileQueryListener: () => void;

    constructor(
        changeDetectorRef: ChangeDetectorRef,
        media: MediaMatcher,
        private matIconRegistry: MatIconRegistry,
        private domSanitizer: DomSanitizer,
        private _globals: Globals
    ) {
        this.mobileQuery = media.matchMedia("(max-width: 600px)");
        this._mobileQueryListener = () => changeDetectorRef.detectChanges();
        this.mobileQuery.addListener(this._mobileQueryListener);
    }

    ngAfterViewInit() {}
}
