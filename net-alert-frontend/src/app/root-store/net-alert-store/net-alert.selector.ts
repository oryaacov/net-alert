import { createSelector,  } from '@ngrx/store';
import * as fromNetAlertState from './net-alert.state';

export const netAlertState = (state: fromNetAlertState.NetAlertState) => state;

export const selectNetworkInfo = createSelector(
  netAlertState,
  (state: fromNetAlertState.NetAlertState) => state.NetworkInfo
);

export const selectProfiles = createSelector(
  netAlertState,
  (state: fromNetAlertState.NetAlertState) => state.Profiles
);