import { Component, Inject, OnInit } from '@angular/core';
import { UntypedFormControl, Validators } from '@angular/forms';
import {
  MatLegacyDialogRef as MatDialogRef,
  MAT_LEGACY_DIALOG_DATA as MAT_DIALOG_DATA,
} from '@angular/material/legacy-dialog';
import { CountryCallingCodesService, CountryPhoneCode } from 'src/app/services/country-calling-codes.service';
import { formatPhone } from 'src/app/utils/formatPhone';

export enum EditDialogType {
  PHONE = 1,
  EMAIL = 2,
}

@Component({
  selector: 'cnsl-edit-email-dialog',
  templateUrl: './edit-dialog.component.html',
  styleUrls: ['./edit-dialog.component.scss'],
})
export class EditDialogComponent implements OnInit {
  public isPhone: boolean = false;
  public isVerified: boolean = false;
  public phoneCountry: string = 'CH';
  public valueControl: UntypedFormControl = new UntypedFormControl(['', [Validators.required]]);
  public EditDialogType: any = EditDialogType;
  public selected: CountryPhoneCode | undefined;
  public countryPhoneCodes: CountryPhoneCode[] = [];
  constructor(
    public dialogRef: MatDialogRef<EditDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any,
    private countryCallingCodesService: CountryCallingCodesService,
  ) {
    this.valueControl.setValue(data.value);
    if (data.type === EditDialogType.PHONE) {
      this.isPhone = true;
    }
  }

  public setCountryCallingCode(): void {
    let value = (this.valueControl?.value as string) || '';
    this.valueControl?.setValue('+' + this.selected?.countryCallingCode + ' ' + value.replace(/\+[0-9]*\s/, ''));
  }

  ngOnInit(): void {
    if (this.isPhone) {
      // Get country phone codes and set selected flag to guessed country or default country
      this.countryPhoneCodes = this.countryCallingCodesService.getCountryCallingCodes();
      const phoneNumber = formatPhone(this.valueControl?.value);
      this.selected = this.countryPhoneCodes.find((code) => code.countryCode === phoneNumber.country);
      this.valueControl.setValue(phoneNumber.phone);
    }
  }

  closeDialog(): void {
    this.dialogRef.close();
  }

  closeDialogWithValue(): void {
    this.dialogRef.close({ value: this.valueControl.value, isVerified: this.isVerified });
  }
}
