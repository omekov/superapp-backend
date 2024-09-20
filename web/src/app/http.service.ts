import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

export interface IMark {
  Name: string;
}
export interface IModel {
  Name: string;
}
export interface IVolume {
  Value: number;
}
export interface ISpecification {
  Amount: number;
  Year: number;
}

@Injectable({
  providedIn: 'root'
})
export class HttpService {
  private apiUrl = '/transport'; // URL для запроса
  // private apiUrl = '/transport'; // URL для запроса

  constructor(private http: HttpClient) {}

  getMarks(): Observable<any> {
    return this.http.get(this.apiUrl);
  }
  getModels(mark: string): Observable<any> {
    return this.http.get(this.apiUrl + "?mark=" + mark);
  }
  getVolumes(mark: string, model: string): Observable<any> {
    return this.http.get(this.apiUrl + "?mark=" + mark + "&model=" + model);
  }
  getSpecifications(mark: string, model: string, volume: number): Observable<any> {
    return this.http.get(this.apiUrl + "?mark=" + mark + "&model=" + model+ "&volume=" + volume);
  }
}
