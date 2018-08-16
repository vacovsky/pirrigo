import { BrowserModule } from "@angular/platform-browser";
import { NgModule } from "@angular/core";

import { AppComponent } from "./app.component";
import { ScheduleComponent } from "./components/schedule/schedule.component";
import { SidenavComponent } from "./components/sidenav/sidenav.component";
import { Globals } from "./services/globals.service";
import { StationService } from "./services/station.service";

import { HttpClient, HttpClientModule } from "@angular/common/http";
import { MediaMatcher } from "@angular/cdk/layout";
import { ReactiveFormsModule } from "@angular/forms";
import { FormControl } from "@angular/forms";

import {
    MatSidenavModule,
    MatIconModule,
    MatToolbarModule,
    MatFormFieldModule,
    MatListModule,
    MatCardModule,
    MatGridListModule,
    MatButtonModule,
    MatOptionModule,
    MatSelectModule,
    // MatAutocompleteModule,
    MatInputModule,
    MatRadioModule,
    MatProgressSpinnerModule,
    MatExpansionModule,
    MatDialogModule,
    MatDialogRef
} from "@angular/material";
import { BrowserAnimationsModule } from "@angular/platform-browser/animations";
import { ZonesComponent } from "./components/zones/zones.component";
import { SidenavService } from "./services/sidenav.service";
import { CookieService } from "ngx-cookie-service";
import { PirriCalendarComponent } from "./components/calendar/calendar.component";
// import { FullCalendarModule } from "ng-fullcalendar";

@NgModule({
    declarations: [
        AppComponent,
        SidenavComponent,
        ScheduleComponent,
        ZonesComponent,
        PirriCalendarComponent
    ],
    imports: [
        MatProgressSpinnerModule,
        BrowserModule,
        BrowserAnimationsModule,
        MatSidenavModule,
        MatToolbarModule,
        MatListModule,
        MatSelectModule,
        MatCardModule,
        HttpClientModule,
        MatGridListModule,
        MatExpansionModule,
        MatButtonModule,
        MatIconModule,
        MatRadioModule,
        MatFormFieldModule,
        MatOptionModule,
        MatSelectModule,
        MatInputModule,
        MatDialogModule,
        BrowserModule,
        // FullCalendarModule
    ],
    providers: [
        Globals,
        HttpClient,
        MediaMatcher,
        StationService,
        SidenavService,
        CookieService
    ],
    bootstrap: [AppComponent]
})
export class AppModule {}
