import { Component, OnInit } from '@angular/core';
import { ApiClientService } from 'src/app/services/apiclient.service';
import { ChartTransformService } from 'src/app/services/chart-transform.service';

@Component({
  selector: 'app-analytics',
  templateUrl: './analytics.component.html',
  styleUrls: ['./analytics.component.css']
})
export class AnalyticsComponent implements OnInit {

  days = 14

  chartData: Record<string, any> = {}

  colorScheme = 'viridis'

  constructor(
    private _api: ApiClientService,
    private _cts: ChartTransformService
  ) { }

  ngOnInit(): void {
    this.loadAllCharts()
  }

  loadAllCharts(): void {
    // ponytail: load all 5 charts in parallel — each is independent
    for (let id = 1; id <= 5; id++) {
      this._api.loadChartByID(id, this.days).subscribe(data => {
        this.chartData[id] = this._cts.transformChartDataForNgxChartsWithStringLabels(data)
      })
    }
  }
}
