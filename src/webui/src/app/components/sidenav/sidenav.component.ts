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
        { name: "Dashboard", icon: "equalizer" },
        { name: "Alerts", icon: "warning" },
        { name: "Pumps", icon: "transform" },
        { name: "Tank", icon: "opacity" },
        { name: "Device", icon: "developer_board" },
        { name: "Account", icon: "account_box" },
        { name: "Settings", icon: "settings" },
        { name: "Support", icon: "help_outline" }
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
        // this._watrapi.loadDevices(); // matIconRegistry.addSvgIconSet(domSanitizer.bypassSecurityResourceUrl('/assets/mdi.svg'));
    }

    ngAfterViewInit() {
        // this._sidenav.changeDevice(this._globals.currentDevice) ? this._globals.currentDevice : '';
    }
}
