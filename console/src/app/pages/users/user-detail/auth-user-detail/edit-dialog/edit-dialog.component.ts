import { Component, Inject } from '@angular/core';
import { UntypedFormControl, Validators } from '@angular/forms';
import { MAT_LEGACY_DIALOG_DATA as MAT_DIALOG_DATA, MatLegacyDialogRef as MatDialogRef } from '@angular/material/legacy-dialog';
import { parsePhoneNumber } from 'libphonenumber-js';

export enum EditDialogType {
  PHONE = 1,
  EMAIL = 2,
}

@Component({
  selector: 'cnsl-edit-email-dialog',
  templateUrl: './edit-dialog.component.html',
  styleUrls: ['./edit-dialog.component.scss'],
})
export class EditDialogComponent {
  public isPhone: boolean = false;
  public isVerified: boolean = false;
  public phoneCountry: string = 'CH';
  public valueControl: UntypedFormControl = new UntypedFormControl(['', [Validators.required]]);
  public EditDialogType: any = EditDialogType;
  constructor(public dialogRef: MatDialogRef<EditDialogComponent>, @Inject(MAT_DIALOG_DATA) public data: any) {
    this.valueControl.setValue(data.value);
    if (data.type === EditDialogType.PHONE) {
      this.isPhone = true;
    }

    this.valueControl.valueChanges.subscribe((value) => {
      if (value && value.length > 1) {
        this.changeValue(value);
      }
    });
  }

  private changeValue(changedValue: string): void {
    if (this.isPhone && changedValue) {
      try {
        const phoneNumber = parsePhoneNumber(changedValue ?? '', 'CH');
        if (phoneNumber) {
          const formmatted = phoneNumber.formatInternational();
          this.phoneCountry = phoneNumber.country || '';
          if (formmatted !== this.valueControl.value) {
            this.valueControl.setValue(formmatted);
          }
        }
      } catch (error) {
        console.error(error);
      }
    }
  }

  closeDialog(): void {
    this.dialogRef.close();
  }

  closeDialogWithValue(): void {
    this.dialogRef.close({ value: this.valueControl.value, isVerified: this.isVerified });
  }
}
