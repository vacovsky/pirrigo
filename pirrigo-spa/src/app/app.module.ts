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
import { MatTabsModule } from '@angular/material/tabs';
import { MatIconModule } from '@angular/material/icon';
import { MatTableModule } from '@angular/material/table';
import { adapterFactory } from 'angular-calendar/date-adapters/date-fns';
import { CalendarModule, DateAdapter } from 'angular-calendar';
import { StatusComponent } from './components/stations/status/status.component';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatSliderModule } from '@angular/material/slider';
import { MatListModule } from '@angular/material/list';
import { MatDialogModule } from '@angular/material/dialog';
import { EditScheduleDialog } from './components/calendar/calendar.component';
import { MatButtonModule } from '@angular/material/button';
import { MatInputModule } from '@angular/material/input';
import { NgxMaterialTimepickerModule } from 'ngx-material-timepicker';
import { MatSelectModule } from '@angular/material/select';



@NgModule({
  declarations: [
    AppComponent,
    StationsComponent,
    HistoryComponent,
    CalendarComponent,
    LogsComponent,
    SettingsComponent,
    UsageCalculatorComponent,
    StatusComponent,
    EditScheduleDialog
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
    MatInputModule,
    MatNativeDateModule,
    FormsModule,
    MatProgressSpinnerModule,
    MatButtonModule,
    MatCardModule,
    MatExpansionModule,
    MatSelectModule,
    NgxMaterialTimepickerModule,
    MatIconModule,
    MatTabsModule,
    MatTableModule,
    CalendarModule.forRoot({
      provide: DateAdapter,
      useFactory: adapterFactory,
    }),
    MatProgressBarModule,
    MatSliderModule,
    MatListModule,
    MatDialogModule
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
