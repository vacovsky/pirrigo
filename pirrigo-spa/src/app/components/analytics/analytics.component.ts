import { Component, OnInit } from '@angular/core';
import { ApiClientService } from 'src/app/services/apiclient.service';
import * as moment from 'moment';
import { ChartTransformService } from 'src/app/services/chart-transform.service';

@Component({
  selector: 'app-analytics',
  templateUrl: './analytics.component.html',
  styleUrls: ['./analytics.component.css']
})
export class AnalyticsComponent implements OnInit {

  startDate: Date
  endDate: Date

  overallUsageChartData: any;
  zoneActivityChartData: any;
  usageByDOWChartData: any;

  chartLoading: boolean
  chartTab: number = 4;

  chart4Options: any = {
    // options
    legend: true,
    showLabels: true,
    animations: true,
    xAxis: true,
    yAxis: true,
    showYAxisLabel: true,
    showXAxisLabel: true,
    xAxisLabel: 'hour of the day',
    yAxisLabel: 'minutes',
    timeline: true,
  }


  chart2Options: any = {
    // options
    legend: true,
    showLabels: true,
    animations: true,
    xAxis: true,
    yAxis: true,
    showYAxisLabel: true,
    showXAxisLabel: true,
    xAxisLabel: 'day of week',
    yAxisLabel: 'minutes',
    timeline: true,
  }


  colorScheme = {
    domain: ['#5AA454', '#E44D25', '#CFC0BB', '#7aa3e5', '#a8385d', '#aae3f5']
  };

  constructor(
    private _api: ApiClientService,
    private _cts: ChartTransformService
  ) { }

  ngOnInit(): void {
    this.endDate = moment().toDate()
    this.startDate = moment().add(-14, "d").toDate()
    this.loadChart()
  }


  loadChart(): void {
    // if (this.chartTab == 4) {
    this.loadZoneUsageChart()
    // } else if (this.chartTab == 2) {
    this.loadStatsActivityByDayOfWeek()
    // }

  }

  // "/stats/1": statsActivityByStation,
  // "/stats/2": statsActivityByDayOfWeek,
  // "/stats/3": statsActivityPerStationByDOW,
  // "/stats/4": statsStationActivity,

  loadZoneUsageChart(): void {
    this.chartLoading = true

    this._api.loadChartByID(4, moment(this.startDate).unix(), moment(this.endDate).unix()).subscribe(data => {
      this.zoneActivityChartData = this._cts.transformChartDataForNgxChartsWithStringLabels(data)
      this.chartLoading = false
      console.log(data)
    })
  }


  loadStatsActivityByDayOfWeek(): void {
    this.chartLoading = true
    this._api.loadChartByID(2, moment(this.startDate).unix(), moment(this.endDate).unix()).subscribe(data => {
      this.usageByDOWChartData = this._cts.transformChartDataForNgxChartsWithStringLabels(data)
      console.log(data)
      this.chartLoading = false
    })
  }
}
