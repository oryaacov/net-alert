import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from 'src/environments/environment';

@Injectable({
  providedIn: 'root'
})
export class DataService {

  constructor (private http: HttpClient) {}

  getAllProfiles() {
    return this.http.get(environment.baseURL+'/profiles');
  }

  getOwner() {
    return this.http.get(environment.baseURL+'/master');
  }

  getNetworkInfo() {
    return this.http.get(environment.baseURL+'/network');
  }

  getIsAlive() {
    return this.http.get(environment.baseURL+'/alive');
  }

}
