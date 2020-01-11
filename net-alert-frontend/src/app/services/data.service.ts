import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from 'src/environments/environment';
import { Owner, Profile } from '../root-store/net-alert-store/net-alert.state';

@Injectable({
  providedIn: 'root'
})
export class DataService {

  constructor (private http: HttpClient) {}

  getAllProfiles() {
    return this.http.get(environment.baseURL+'/profiles');
  }

  startRequest() {
    return this.http.get(environment.baseURL+'/start',{responseType:"text"});
  }

  updateProfiles(profiles:Profile[]) {
    console.log(profiles);
    return this.http.post(environment.baseURL+'/profiles',profiles);
  }

  getOwner() {
    return this.http.get(environment.baseURL+'/master');
  }
  
  updateOwner(owner:Owner) {
    return this.http.post(environment.baseURL+'/master',owner);
  }

  getNetworkInfo() {
    return this.http.get(environment.baseURL+'/network');
  }

  getIsAlive() {
    return this.http.get(environment.baseURL+'/alive');
  }

}
