import { Injectable } from '@angular/core';
import { ChartData } from '../structs/chart-data';
import { NgxChartData, NgxChartSeries } from '../structs/ngx-chart-data';
import * as moment from 'moment';

@Injectable({
  providedIn: 'root'
})
export class ChartTransformService {

  constructor() { }

  transformChartDataForNgxCharts(input: ChartData): NgxChartData[] {
    let result: NgxChartData[] = [];
    if (input != undefined) {
      let counter = 0
      for (var series of input.Series) {
        let dataCounter = 0;
        let data: NgxChartData = new NgxChartData();

        data.series = [];
        data.name = series;

        while (dataCounter < input.Data[counter].length) {
          let seriesData = new NgxChartSeries();
          seriesData.name = moment.unix(
            Number(input.Labels[dataCounter]) / 1000
          ).format()


          seriesData.value = Number(input.Data[counter][dataCounter])
          data.series.push(seriesData)
          dataCounter++;
        }
        result.push(
          data
        )
        counter++;
      }
    }

    console.log(input, result)
    return result
  }



  transformChartDataForNgxChartsWithStringLabels(input: ChartData): NgxChartData[] {
    let result: NgxChartData[] = [];
    if (input != undefined) {
      let counter = 0
      for (var series of input.Series) {
        let dataCounter = 0;
        let data: NgxChartData = new NgxChartData();

        data.series = [];
        data.name = series;

        while (dataCounter < input.Data[counter].length) {
          let seriesData = new NgxChartSeries();
          seriesData.name = input.Labels[dataCounter]
          seriesData.value = Number(input.Data[counter][dataCounter])
          data.series.push(seriesData)
          dataCounter++;
        }
        result.push(
          data
        )
        counter++;
      }
    }

    console.log(input, result)
    return result
  }


  private convertDate(epoch: number): number {
    return Math.round(epoch / 1000);
  }
}