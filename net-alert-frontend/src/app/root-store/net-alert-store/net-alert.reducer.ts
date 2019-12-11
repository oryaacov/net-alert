import { createReducer, on, Action, State, ActionReducerMap, MetaReducer } from '@ngrx/store';
import { NetAlertState, Profile } from './net-alert.state'
import * as NetAlertActions from './net-alert.actions'
import { InjectionToken } from '@angular/core';
import { environment } from 'src/environments/environment';
import { state } from '@angular/animations';

export const initialState: NetAlertState = {
  Profiles: null,
  NetworkInfo: null,
  isLoading: false,
  error: null,
  Owner:null
}

export const netAlertReducer = createReducer<NetAlertState | undefined>(
  initialState,
  on(NetAlertActions.getAllProfiles, (state) => { return { ...state, isLoading: true, error: null } }),
  on(NetAlertActions.getNetworkInfo, (state) => { return { ...state, isLoading: true, error: null } }),
  on(NetAlertActions.getOwnerInfo, (state) => { return { ...state, isLoading: true, error: null } }),
  on(NetAlertActions.loadProfilesSuccess, (state, { payload }) => { return { ...state, Profiles: payload, isLoading: false, error: null } }),
  on(NetAlertActions.loadNetworkInfoSuccess, (state, { payload }) => { return { ...state, NetworkInfo: payload, isLoading: false, error: null } }),
  on(NetAlertActions.loadOwnerSuccess, (state, { payload }) => {return { ...state, Owner: payload, isLoading: false, error: null } }),
  on(NetAlertActions.loadProfilesFailure, (state, { error }) => { return  { ...state, isLoading: false, error: error } }),
  on(NetAlertActions.loadNetworkInfoFailure, (state, { error }) => { return { ...state, isLoading: false, error: error } }),
  on(NetAlertActions.loadOwnerFailure, (state, { error }) => { return { ...state, isLoading: false, error: error } }),
  on(NetAlertActions.updateProfiles,(state)=>{return {...state,isLoading:true,error:null}}),
  on(NetAlertActions.updateProfilesSuccess,(state)=>{return {...state,isLoading:false}}),
  on(NetAlertActions.updateProfileFailure,(state,error)=>{return {...state,isLoading:false,error:error}}),
  on(NetAlertActions.updateOwnerInfo,(state)=>{return {...state,isLoading:true,error:null}}),
  on(NetAlertActions.updateOwnerInfoSuccess,(state)=>{return {...state,isLoading:false}}),
  on(NetAlertActions.updateOwnerFailure,(state,error)=>{return {...state,isLoading:false,error:error}})
  );

