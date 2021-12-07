import { Component, OnInit } from '@angular/core';
import { ApiClientService } from 'src/app/services/apiclient.service';
import { Gpio } from 'src/app/structs/gpio';

@Component({
  selector: 'app-settings',
  templateUrl: './settings.component.html',
  styleUrls: ['./settings.component.css']
})
export class SettingsComponent implements OnInit {


  gpioList: Gpio[]
  gpioCurrentCommon: Gpio

  constructor(
    private _api: ApiClientService
  ) { }


  ngOnInit(): void {
    this._api.getAllGPIOs().subscribe((data) => {
      this.gpioList = data.gpios
      for (let g of data.gpios) {
        if (g.Common == true) {
          this.gpioCurrentCommon = g
        }
      }
    })
  }


  setCommon(gpio: number): void {
    this._api.setCommonWireGpio(gpio).subscribe(() => {
      this.ngOnInit()
    })
  }
}
