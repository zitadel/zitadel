import { Component, Inject } from '@angular/core';
import { AbstractControl, FormControl, ValidatorFn, Validators } from '@angular/forms';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { Moment } from 'moment';
import { MachineKeyType } from 'src/app/proto/generated/management_pb';

export function afterNowValidator(): ValidatorFn {
    const now = new Date();
    console.log(now);
    return (control: AbstractControl): { [key: string]: any; } | null => {
        const forbidden = control.value.diff(now) > 0;
        return forbidden ? { forbiddenDate: { value: control.value } } : control.value.diff(now);
    };
}

@Component({
    selector: 'app-add-key-dialog',
    templateUrl: './add-key-dialog.component.html',
    styleUrls: ['./add-key-dialog.component.scss'],
})
export class AddKeyDialogComponent {
    public startDate: Date = new Date();
    types: MachineKeyType[] = [
        MachineKeyType.MACHINEKEY_JSON,
    ];
    date!: Moment;
    public type: MachineKeyType = MachineKeyType.MACHINEKEY_JSON;
    public dateControl: FormControl = new FormControl('', [Validators.required, afterNowValidator]);

    constructor(
        public dialogRef: MatDialogRef<AddKeyDialogComponent>,
        @Inject(MAT_DIALOG_DATA) public data: any,
    ) {
        this.dateControl.valueChanges.subscribe(value => {
            console.log(this.dateControl);

        });
    }

    public closeDialog(): void {
        this.dialogRef.close(false);
    }

    public closeDialogWithSuccess(): void {
        this.dialogRef.close({ type: this.type, date: this.date });
    }
}
