import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { StoreDevtoolsModule } from '@ngrx/store-devtools'
import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { StoreModule } from '@ngrx/store';
import { HttpClientModule } from '@angular/common/http';
import { EffectsModule } from '@ngrx/effects';
import { NetAlertEffects } from './root-store/net-alert-store/net-alert.effects';
import { ProfileComponent } from './components/profile/profile.component';
import { ProfileContainerComponent } from './containers/profile-container/profile-container.component';
import { RootStoreModule } from './root-store/root-store.module';
import { environment } from 'src/environments/environment';
import { reducers } from './root-store/root-state';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { NgxLoadingModule } from 'ngx-loading';
import { MatFormFieldModule, MatInputModule, MatCheckboxModule, MatCardModule, MatListModule, MatTableModule, MatDialogModule, MatTabsModule, MatIconModule, MatToolbarModule, MatGridListModule } from '@angular/material';
import { ToastrModule } from 'ngx-toastr';
import { OwnerComponent } from './components/owner/owner.component';
import { ReactiveFormsModule, FormsModule } from '@angular/forms';
import { MyNetworkComponent } from './components/my-network/my-network.component';
import { DialogBoxComponent } from './dialogs/dialog-box/dialog-box.component';
import { DragDropModule } from '@angular/cdk/drag-drop'
@NgModule({
  declarations: [
    AppComponent,
    ProfileContainerComponent,
    ProfileComponent,
    OwnerComponent,
    MyNetworkComponent,
    DialogBoxComponent
  ],
  imports: [
    HttpClientModule,
    ReactiveFormsModule,
    MatGridListModule,
    MatToolbarModule,
    MatFormFieldModule,
    MatListModule,
    DragDropModule,
    FormsModule,
    MatTabsModule,
    MatIconModule,
    MatTableModule,
    MatDialogModule,
    BrowserModule,
    MatCardModule,
    MatFormFieldModule,
    MatCheckboxModule,
    MatInputModule,
    BrowserAnimationsModule,
    ToastrModule.forRoot({ timeOut: 3000, positionClass: 'toast-bottem-left' }),
    NgxLoadingModule,
    AppRoutingModule,
    StoreModule.forRoot(reducers),
    EffectsModule.forRoot([NetAlertEffects]),
    StoreDevtoolsModule.instrument({
      maxAge: 25, // Retains last 25 states
      logOnly: environment.production, // Restrict extension to log-only mode
    }),
    RootStoreModule
  ],
  providers: [],
  entryComponents: [DialogBoxComponent],
  bootstrap: [AppComponent]
})
export class AppModule { }
