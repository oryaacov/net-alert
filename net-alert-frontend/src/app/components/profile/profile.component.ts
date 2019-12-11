import { Component, OnInit, ViewChild } from '@angular/core';
import { Profile } from 'src/app/root-store/net-alert-store/net-alert.state';
import {  Store } from '@ngrx/store';
import { AppStates } from 'src/app/root-store/root-state';
import { MatTable, MatDialog } from '@angular/material';
import { DialogBoxComponent } from 'src/app/dialogs/dialog-box/dialog-box.component';
import {updateProfiles} from '../../root-store/net-alert-store/net-alert.actions'

@Component({
  selector: 'app-profile',
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.scss']
})
export class ProfileComponent implements OnInit {
  @ViewChild(MatTable, {static:true}) table: MatTable<any>;
 
  dataSource: Profile[]=[];
  displayedColumns: string[] = ['mac', 'name', 'sites','created_at', 'action'];
  backup:any;
  constructor(private store: Store<AppStates>,public dialog: MatDialog) {
  }
  ngOnInit() {
    setTimeout(() => {
      this.store.dispatch({ type: '[Profile Component] getAllProfiles' });
    }, 0);
    this.store.select(r => r.netAlert).subscribe(state => {
      this.dataSource = state.Profiles;
      this.table.dataSource=this.dataSource;
    });

  }
  openDialog(action, obj,i) {
    obj.action = action;
    this.backup=JSON.stringify(this.dataSource)
    const dialogRef = this.dialog.open(DialogBoxComponent, {
      data: {action:action,data:obj,index:i}
    });

    dialogRef.afterClosed().subscribe(result => {
      let changed:boolean=true
      if (result.event == 'Add') {
        this.dataSource=[...this.dataSource,{...result.data}]
      } else if (result.event == 'Update') {
        this.dataSource[result.index]={...result.data}
      } else if (result.event == 'Delete') {
        this.deleteRowData(result.data);
      }else{
        this.dataSource= [...JSON.parse(this.backup)]
        changed=false;
      }
      this.table.dataSource=[...this.dataSource]
      this.table.renderRows()
      if (changed){
        this.store.dispatch(updateProfiles(this.dataSource))
      }  
    });
  }

  addRowData(row_obj) {
    var d = new Date();
    let profile :Profile= {Mac:"",NickName:"",CreateDate:new Date(),Sites:[]}
    this.dataSource.push(
      profile
    );
    this.table.renderRows();

  }
  updateRowData(row_obj) {
    this.dataSource = this.dataSource.filter((value, key) => {
      if (value.Mac == row_obj.Mac) {
        value.NickName = row_obj.NickName;
        value.Sites=row_obj.Sites
      }
      return true;
    });
  }
  deleteRowData(row_obj) {
    this.dataSource = this.dataSource.filter((value, key) => {
      return value.Mac != row_obj.Mac;
    });
  }
}