import { createAction, props } from '@ngrx/store';
import { Profile, NetworkInfo, Owner } from './net-alert.state';
import { Observable } from 'rxjs';

export const getAllProfiles = createAction('[Profile Component] getAllProfiles');
export const getNetworkInfo = createAction('[Profile Component] getNetworkInfo');
export const getOwnerInfo = createAction('[Owner Component] getOwnerInfo');
export const loadNetworkInfoSuccess = createAction('[Profile Component] loadNetworkInfoSuccess', (networkInfo: NetworkInfo) =>
    ({ payload: networkInfo }));
export const loadProfilesSuccess = createAction('[Profile Component] loadProfilesSuccess', (profiles: Profile[]) =>
    ({ payload: profiles }));
export const loadOwnerSuccess = createAction('[Owner Component] loadOwnerSuccess', (owner: Owner) =>
    ({ payload: owner }));
export const loadProfilesFailure = createAction('[Profile Component] loadProfilesFailure', (err: string) => ({ error: err }));
export const loadOwnerFailure = createAction('[Owner Component] loadOwnerFailure', (err: string) => ({ error: err }));
export const loadNetworkInfoFailure = createAction('[Profile Component] loadNetworkInfoFailure', (err: string) => ({ error: err }));
