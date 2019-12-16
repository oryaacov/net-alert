import { Component } from '@angular/core';
import { Store } from '@ngrx/store';
import { AppStates } from './root-store/root-state';
import { ToastrService } from 'ngx-toastr';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})


export class AppComponent {

  isLoading$: boolean = false;

  constructor(private store: Store<AppStates>, private toastr: ToastrService) {
    store.select(r => r.netAlert.isLoading).subscribe(isLoading => this.isLoading$ = isLoading);
    store.select(r=> r.netAlert.error).subscribe(err=>this.handleError(err));
   }

  title = 'net-alert-frontend';

  private handleError(err) {
  }
}
