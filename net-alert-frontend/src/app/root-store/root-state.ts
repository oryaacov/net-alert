import { NetAlertState } from './net-alert-store/net-alert.state';
import { createSelector, ActionReducerMap, Action, MetaReducer } from '@ngrx/store';
import * as netAlertRerucers from './net-alert-store/net-alert.reducer';
import { InjectionToken } from '@angular/core';
import { environment } from 'src/environments/environment';
 
export interface AppStates {
  netAlert: NetAlertState;
}
 


const rootReducers: ActionReducerMap<AppStates, Action> = {
    netAlert:netAlertRerucers.netAlertReducer
};
  
  export const reducers = new InjectionToken<ActionReducerMap<AppStates, Action>>(
    'Root reducers token',
    {
      factory: () => rootReducers,
    },
  );
  
  export const metaReducers: MetaReducer<NetAlertState>[] = !environment.production
    ? []
    : [];
  
  