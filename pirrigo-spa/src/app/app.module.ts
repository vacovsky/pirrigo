import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { HttpClient, HttpClientModule } from '@angular/common/http';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { NgxChartsModule } from '@swimlane/ngx-charts'
import { GaugeModule } from 'angular-gauge';
import { MatGridListModule } from '@angular/material/grid-list';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatNativeDateModule } from '@angular/material/core';
import { FormsModule } from '@angular/forms';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { AppComponent } from './app.component';
import { StationsComponent } from './components/stations/stations.component';
import { HistoryComponent } from './components/history/history.component';
import { CalendarComponent } from './components/calendar/calendar.component';
import { LogsComponent } from './components/logs/logs.component';
import { SettingsComponent } from './components/settings/settings.component';
import { UsageCalculatorComponent } from './components/usage-calculator/usage-calculator.component';
import { ChartTransformService } from './services/chart-transform.service';
import { ApiClientService } from './services/apiclient.service';
import { GlobalsService } from './services/globals.service';
import { MatCardModule } from '@angular/material/card';
import { MatExpansionModule } from '@angular/material/expansion';


@NgModule({
  declarations: [
    AppComponent,
    StationsComponent,
    HistoryComponent,
    CalendarComponent,
    LogsComponent,
    SettingsComponent,
    UsageCalculatorComponent
  ],
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    HttpClientModule,
    NgxChartsModule,
    GaugeModule.forRoot(),
    MatGridListModule,
    MatDatepickerModule,
    MatFormFieldModule,
    MatNativeDateModule,
    FormsModule,
    MatProgressSpinnerModule,
    MatCardModule,
    MatExpansionModule
  ],
  providers: [
    HttpClient,
    ChartTransformService,
    ApiClientService,
    GlobalsService

  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
