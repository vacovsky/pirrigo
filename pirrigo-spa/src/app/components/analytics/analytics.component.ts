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
  chartLoading: boolean


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
    this.loadAllCharts()
  }


  loadAllCharts(): void {
    this._api.loadChartByID(4, moment(this.startDate).unix(), moment(this.endDate).unix()).subscribe(data => {
      this.overallUsageChartData = this._cts.transformChartDataForNgxChartsWithStringLabels(data)
      console.log(data)
    })
  }
}
