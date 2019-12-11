import { createAction, props } from '@ngrx/store';
import { Profile, NetworkInfo, Owner } from './net-alert.state';
import { Observable } from 'rxjs';


//*******getAllProfiles********/
export const getAllProfiles = createAction('[Profile Component] getAllProfiles');
export const loadProfilesSuccess = createAction('[Profile Component] loadProfilesSuccess', (profiles: Profile[]) =>
    ({ payload: profiles }));
export const loadProfilesFailure = createAction('[Profile Component] loadProfilesFailure', (err: string) => ({ error: err }));

//*******getNetworkInfo********/
export const getNetworkInfo = createAction('[Profile Component] getNetworkInfo');
export const loadNetworkInfoSuccess = createAction('[Profile Component] loadNetworkInfoSuccess', (networkInfo: NetworkInfo) =>
    ({ payload: networkInfo }));
export const loadNetworkInfoFailure = createAction('[Profile Component] loadNetworkInfoFailure', (err: string) => ({ error: err }));

//*******getOwnerInfo********/
export const getOwnerInfo = createAction('[Owner Component] getOwnerInfo');
export const loadOwnerSuccess = createAction('[Owner Component] loadOwnerSuccess', (owner: Owner) =>
    ({ payload: owner }));
export const loadOwnerFailure = createAction('[Owner Component] loadOwnerFailure', (err: string) => ({ error: err }));

//*******updateOwnerInfo********/
export const updateOwnerInfo = createAction('[Owner Component] updateOwnerInfo', (owner: Owner) =>
    ({ payload: owner }));
export const updateOwnerInfoSuccess = createAction('[Owner Component] updateOwnerInfoSuccess', (owner: Owner) =>
    ({ payload: owner }));
export const updateOwnerFailure = createAction('[Owner Component] updateOwnerFailure', (err: string) => ({ error: err }));


//*******updateProfiles********/
export const updateProfiles = createAction('[Profile Component] updateProfiles', (profiles: Profile[]) =>
    ({ payload: profiles }));
export const updateProfilesSuccess = createAction('[Profile Component] updateProfilesSuccess');
export const updateProfileFailure = createAction('[Profile Component] updateProfileFailure', (err: string) => ({ error: err }));

