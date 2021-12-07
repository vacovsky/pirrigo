import { Injectable } from '@angular/core';
import { HttpHeaders } from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class GlobalsService {

  public currentPage: string = "stations";
  public uriStem: string = "http://192.168.111.130"

  public headers: HttpHeaders = new HttpHeaders({
    'Accept': 'application/json',
    'Content-Type': 'application/json'
  });

  public statusRefreshRateMs: number = 3000;

  constructor() { }
}
