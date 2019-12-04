import { createAction, props } from '@ngrx/store';
import { Profile, NetworkInfo } from './net-alert.state';
import { Observable } from 'rxjs';

export const getAllProfiles = createAction('[Profile Component] getAllProfiles');
export const getNetworkInfo = createAction('[Profile Component] getNetworkInfo');
// export const loadProfilesSuccess = createAction('[Profile Component] loadProfilesSuccess',props<{payload:Profile[]}>());
export const loadNetworkInfoSuccess = createAction('[Profile Component] loadNetworkInfoSuccess', props<NetworkInfo>());
export const loadProfilesFailure = createAction('[Profile Component] loadProfilesFailure',props<String>());
export const loadNetworkInfoFailure = createAction('[Profile Component] loadNetworkInfoFailure',props<String>());
export const loadProfilesSuccess = createAction('[Profile Component] loadProfilesSuccess', ( profiles: Profile[]) => ({
    payload: profiles,
  }));