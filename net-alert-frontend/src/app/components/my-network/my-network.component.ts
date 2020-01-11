import { Component, OnInit } from '@angular/core';
import { Store } from '@ngrx/store';
import { AppStates } from 'src/app/root-store/root-state';
import { NetworkInfo } from 'src/app/root-store/net-alert-store/net-alert.state';
import { Observable } from 'rxjs';

@Component({
  selector: 'app-my-network',
  templateUrl: './my-network.component.html',
  styleUrls: ['./my-network.component.scss']
})
export class MyNetworkComponent implements OnInit {

  public networkInfo:NetworkInfo
  constructor(private store :Store<AppStates>) { 
    
  }

  ngOnInit() {
    setTimeout(() => {
      this.store.dispatch({ type: '[Profile Component] getNetworkInfo' });
    },0);
    this.store.select(r=>r.netAlert).subscribe(state=>this.networkInfo=state.NetworkInfo);
  }

}
