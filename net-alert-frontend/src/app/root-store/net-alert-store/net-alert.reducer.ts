import { createReducer, on, Action, State, ActionReducerMap, MetaReducer } from '@ngrx/store';
import {NetAlertState, Profile} from './net-alert.state'
import * as NetAlertActions from './net-alert.actions'
import { InjectionToken } from '@angular/core';
import { environment } from 'src/environments/environment';

export const initialState: NetAlertState = {
    Profiles:null,
    NetworkInfo:null,
    isLoading: false,
    error: null
}

export const netAlertReducer = createReducer<NetAlertState | undefined>(
  initialState,
  on(NetAlertActions.loadProfilesSuccess, (state, { payload }) =>{return {...state,Profiles:payload}})
);

  
// export const reducer = createReducer(
//   initialState,
//   on(NetAlertActions.getAllProfiles,state => ({...state,isLoading:true})),
//   on(NetAlertActions.loadProfilesSuccess,(state,{payload})=>{console.log(payload);return {...state,Profiles:payload,isLoading:false}})
  
// );

// export function NetAlertReducer(state: NetAlertState | undefined, action: Action) {
//   console.log("ASdfasdf")
//   return reducer(state, action);
//   }