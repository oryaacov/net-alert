//dialog-box.component.ts
import { Component, Inject, Optional } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA, _MatTabGroupBase } from '@angular/material';
import { Profile } from '../../root-store/net-alert-store/net-alert.state'
import { NumberValueAccessor } from '@angular/forms';
export interface UsersData {
    name: string;
    id: number;
}


@Component({
    selector: 'app-dialog-box',
    templateUrl: './dialog-box.component.html',
    styleUrls: ['./dialog-box.component.scss']
})
export class DialogBoxComponent {

    action: string;
    local_data: any;
    index: number;
    displayedColumns: string[] = ['mac', 'ip', 'name', 'action'];

    constructor(
        public dialogRef: MatDialogRef<DialogBoxComponent>,
        @Optional() @Inject(MAT_DIALOG_DATA) public data: any) {
        this.local_data = { ...data.data };
        this.action = this.data.action;
        this.index = this.data.index;
    }

    ngOnInit() {
        if (this.action == "Add") {
            this.local_data = {};
            this.local_data.Sites = [{ Mac: "", IP: "" }]
        }
    }
    doAction() {
        this.dialogRef.close({ event: this.action, data: this.local_data, index: this.index });
    }

    addSite() {
        this.local_data.Sites = [...this.local_data.Sites, { Mac: "", IP: "" }]
    }
    deleteItem(row_obj) {
        this.local_data.Sites = this.local_data.Sites.filter((value, key) => {
            return value.IP != row_obj.IP;
        });
    }
    closeDialog() {
        this.dialogRef.close({ event: 'Cancel' });
    }


}
