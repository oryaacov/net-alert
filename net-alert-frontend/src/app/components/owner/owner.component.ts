import { Component, OnInit } from '@angular/core';
import { FormGroup, FormControl } from '@angular/forms';
import { AppStates } from 'src/app/root-store/root-state';
import { Store } from '@ngrx/store';
import { Owner } from 'src/app/root-store/net-alert-store/net-alert.state';
import {updateOwnerInfo} from '../../root-store/net-alert-store/net-alert.actions'
import { Subscription } from 'rxjs';
@Component({
  selector: 'app-owner',
  templateUrl: './owner.component.html',
  styleUrls: ['./owner.component.scss']
})
export class OwnerComponent implements OnInit {

  formGroup: FormGroup;
  owner: Owner;
  constructor(private store: Store<AppStates>) { }
  ngOnInit() {
    this.buildFormGroup();
    setTimeout(() => {
      this.store.dispatch({ type: '[Owner Component] getOwnerInfo' });
    }, 0);
    this.store.select(r => r.netAlert).subscribe(r => {
      this.owner = r.Owner;
      if (r.Owner) {
        this.formGroup.patchValue(r.Owner)
      }
    });


  }

 
  


  private buildFormGroup() {
    this.formGroup = new FormGroup({});
    this.formGroup.addControl("Nickname", new FormControl());
    this.formGroup.addControl("Email", new FormControl());
    this.formGroup.addControl("Phone", new FormControl());
    this.formGroup.addControl("GetEmailAlerts", new FormControl());
    this.formGroup.addControl("GetSMSAlerts", new FormControl());
  }
  onSubmit(value) {
    this.store.dispatch(updateOwnerInfo(value));
  }
}
