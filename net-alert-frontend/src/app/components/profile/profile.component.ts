import { Component, OnInit } from '@angular/core';
import { Observable } from 'rxjs';
import { Profile, NetAlertState } from 'src/app/root-store/net-alert-store/net-alert.state';
import { select, Store } from '@ngrx/store';
import { map, flatMap } from 'rxjs/operators';
import { selectNetworkInfo ,selectProfiles} from 'src/app/root-store/net-alert-store/net-alert.selector';
import { AppStates } from 'src/app/root-store/root-state';

@Component({
  selector: 'app-profile',
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.scss']
})
export class ProfileComponent implements OnInit {
  network$: any;
  profiles$: any;
  constructor(private store: Store<AppStates>) { 
  }

  ngOnInit() {
    console.log(this.store)
    this.store.dispatch({ type: '[Profile Component] getAllProfiles' });
    this.store.dispatch({ type: '[Profile Component] getNetworkInfo' });
    this.network$ = this.store.select(r=>r.netAlert.NetworkInfo);
    this.profiles$=this.store.select(r=>r.netAlert.Profiles)
   }
}