import { createAction, props } from '@ngrx/store';
import { Profile, NetworkInfo } from './net-alert.state';
import { Observable } from 'rxjs';

export const getAllProfiles = createAction('[Profile Component] getAllProfiles');
export const getNetworkInfo = createAction('[Profile Component] getNetworkInfo');
export const loadNetworkInfoSuccess = createAction('[Profile Component] loadNetworkInfoSuccess', (networkInfo: NetworkInfo) =>
    ({ payload: networkInfo }));
export const loadProfilesSuccess = createAction('[Profile Component] loadProfilesSuccess', (profiles: Profile[]) =>
    ({ payload: profiles }));
export const loadProfilesFailure = createAction('[Profile Component] loadProfilesFailure', (err: string) => ({ error: err }));
export const loadNetworkInfoFailure = createAction('[Profile Component] loadNetworkInfoFailure', (err: string) => ({ error: err }));
