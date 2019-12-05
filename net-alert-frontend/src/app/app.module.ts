import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import {StoreDevtoolsModule} from '@ngrx/store-devtools'
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
import { reducers} from './root-store/root-state';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { NgxLoadingModule } from 'ngx-loading';
import { MatFormFieldModule, MatInputModule, MatCheckboxModule, MatCardModule } from '@angular/material';
import { ToastrModule } from 'ngx-toastr';
import { OwnerComponent } from './components/owner/owner.component';
import { ReactiveFormsModule, FormsModule } from '@angular/forms';
@NgModule({
  declarations: [
    AppComponent,
    ProfileContainerComponent,
    ProfileComponent,
    OwnerComponent
  ],
  imports: [
    HttpClientModule,
    ReactiveFormsModule,
    MatFormFieldModule,
    FormsModule,
    BrowserModule,
    MatCardModule,
    MatFormFieldModule,
    MatCheckboxModule,
    MatInputModule,
    BrowserAnimationsModule,
    ToastrModule.forRoot({timeOut:3000,positionClass:'toast-bottem-left'}),
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
  bootstrap: [AppComponent]
})
export class AppModule { }
