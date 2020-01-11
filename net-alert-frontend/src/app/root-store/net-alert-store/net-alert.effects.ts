import { Injectable } from '@angular/core';
import { Actions, createEffect, ofType } from '@ngrx/effects';
import { map, mergeMap, catchError, tap } from 'rxjs/operators';
import { DataService } from '../../services/data.service';
import { of, pipe } from 'rxjs';
import * as NetAlertActions from './net-alert.actions'
import { selectProfiles } from './net-alert.state';
import { ToastrService } from 'ngx-toastr';

@Injectable()
export class NetAlertEffects {

  loadProfiles$ = createEffect(() =>
    this.actions$.pipe(
      ofType(NetAlertActions.getAllProfiles),
      mergeMap(() => this.dataService.getAllProfiles()
        .pipe(
          map(NetAlertActions.loadProfilesSuccess),
          catchError(err => {
            this.toasterService.error(err.error)
            return of(NetAlertActions.loadProfilesFailure(err))})
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
          tap(()=>this.toasterService.success("Profiles updated")),
          catchError(err => {
            console.log(err);
            this.toasterService.error(err.error);
            return of(NetAlertActions.updateProfileFailure(err))
          })
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
          catchError(err => {
            this.toasterService.error(err.error)
            return of(NetAlertActions.loadOwnerFailure(err))
          })
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
          tap(()=>this.toasterService.success("Owner updated")),
          catchError(err => {
            this.toasterService.error(err.error)
            return of(NetAlertActions.updateOwnerFailure(err))})
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
          catchError(err => {
            this.toasterService.error(err.error);
            return of(NetAlertActions.loadNetworkInfoFailure(err))})
        )
      )
    )
  );
  startRequest$ = createEffect(() =>
  this.actions$.pipe(
    ofType(NetAlertActions.startRequest),
    mergeMap(() => this.dataService.startRequest()
      .pipe(tap(res=>this.toasterService.success(res)),
        map(NetAlertActions.startRequestSuccess),
        catchError(err => {
          this.toasterService.error(err.error)
          return of(NetAlertActions.startRequestFailure(err))})
      )
    )
  )
);
  constructor(
    private actions$: Actions,
    private dataService: DataService,
    private toasterService: ToastrService
  ) { }
}