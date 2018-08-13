import { BrowserModule } from "@angular/platform-browser";
import { NgModule } from "@angular/core";

import { AppComponent } from "./app.component";
import { ScheduleComponent } from "./components/schedule/schedule.component";
import { SidenavComponent } from "./components/sidenav/sidenav.component";
import { Globals } from "./services/globals.service";

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

@NgModule({
    declarations: [AppComponent, SidenavComponent, ScheduleComponent],
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
        MatDialogModule
    ],
    providers: [Globals, HttpClient, MediaMatcher],
    bootstrap: [AppComponent]
})
export class AppModule {}
