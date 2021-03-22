import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { parsePhoneNumber } from 'libphonenumber-js';

export enum EditDialogType {
    PHONE = 1,
    EMAIL = 2,
}

@Component({
    selector: 'app-edit-email-dialog',
    templateUrl: './edit-dialog.component.html',
    styleUrls: ['./edit-dialog.component.scss'],
})
export class EditDialogComponent {
    public value: string = '';
    public isPhone: boolean = false;
    public phoneCountry: string = 'CH';
    constructor(public dialogRef: MatDialogRef<EditDialogComponent>,
        @Inject(MAT_DIALOG_DATA) public data: any) {
        this.value = data.value;
        if (data.type == EditDialogType.PHONE) {
            this.isPhone = true;

            if (this.value) {
                const phoneNumber = parsePhoneNumber(this.value ?? '', 'CH');
                if (phoneNumber) {
                    const formmatted = phoneNumber.formatInternational();
                    this.phoneCountry = phoneNumber.country || '';
                    this.value = formmatted;
                }
            }
        }
    }

    changeValue(change: any) {
        const value = change.target.value;
        if (this.isPhone && value) {
            const phoneNumber = parsePhoneNumber(value ?? '', 'CH');
            if (phoneNumber) {
                const formmatted = phoneNumber.formatInternational();
                this.phoneCountry = phoneNumber.country || '';
                this.value = formmatted;
            }
        }
    }

    closeDialog(email: string = ''): void {
        this.dialogRef.close(email);
    }

    closeDialogWithValue(value: string = ''): void {
        this.dialogRef.close(value);
    }
}
