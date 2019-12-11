import { Injectable } from '@angular/core';
import { Actions, createEffect, ofType } from '@ngrx/effects';
import { map, mergeMap, catchError,tap } from 'rxjs/operators';
import { DataService } from '../../services/data.service';
import { of } from 'rxjs';
import * as NetAlertActions from './net-alert.actions'
import { selectProfiles } from './net-alert.state';

@Injectable()
export class NetAlertEffects {

  loadProfiles$ = createEffect(() =>
    this.actions$.pipe(
      ofType(NetAlertActions.getAllProfiles),
      mergeMap(() => this.dataService.getAllProfiles()
        .pipe(
          map(NetAlertActions.loadProfilesSuccess),
          catchError(() => of(NetAlertActions.loadProfilesFailure))
        )
      )
    )
  );

  updateProfiles$ = createEffect(() =>
    this.actions$.pipe(
      ofType(NetAlertActions.updateProfiles),
      mergeMap(action => this.dataService.updateProfiles(action.payload)
        .pipe(
          map(NetAlertActions.updateProfilesSuccess),
          catchError(() => of(NetAlertActions.updateProfileFailure))
        )
      )
    )
  );
  
  loadOwner$ = createEffect(() =>
    this.actions$.pipe(
      ofType(NetAlertActions.getOwnerInfo),
      mergeMap(() => this.dataService.getOwner()
        .pipe(
          map(NetAlertActions.loadOwnerSuccess),
          catchError(() => of(NetAlertActions.loadOwnerFailure))
        )
      )
    )
  );

  updateOwner$ = createEffect(() =>
    this.actions$.pipe(
      ofType(NetAlertActions.updateOwnerInfo),
      mergeMap(action => this.dataService.updateOwner(action.payload)
        .pipe(
          map(NetAlertActions.updateOwnerInfoSuccess),
          catchError(() => of(NetAlertActions.updateOwnerFailure))
        )
      )
    )
  );

  loadNetworkInfo$ = createEffect(() =>
    this.actions$.pipe(
      ofType(NetAlertActions.getNetworkInfo),
      mergeMap(() => this.dataService.getNetworkInfo()
        .pipe(
          map(NetAlertActions.loadNetworkInfoSuccess),
          catchError(() => of(NetAlertActions.loadNetworkInfoFailure))
        )
      )
    )
  );
  
  constructor(
    private actions$: Actions,
    private dataService: DataService
  ) { }
}