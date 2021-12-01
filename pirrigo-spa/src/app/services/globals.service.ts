import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class GlobalsService {

  public currentPage: string = "stations";
  public uriStem: string = "http://192.168.111.130"

  constructor() { }
}
