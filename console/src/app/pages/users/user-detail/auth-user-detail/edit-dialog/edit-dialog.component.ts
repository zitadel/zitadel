import { Component, Inject, OnInit } from '@angular/core';
import { UntypedFormControl, Validators } from '@angular/forms';
import {
  MAT_LEGACY_DIALOG_DATA as MAT_DIALOG_DATA,
  MatLegacyDialogRef as MatDialogRef,
} from '@angular/material/legacy-dialog';
import { getCountryForTimezone } from 'countries-and-timezones';
import { CountryCode, parsePhoneNumber } from 'libphonenumber-js';
import { CountryCallingCodesService, CountryPhoneCode } from 'src/app/services/country-calling-codes.service';

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
  public phoneCountry: string =
    (getCountryForTimezone(Intl.DateTimeFormat().resolvedOptions().timeZone)?.id as string) ?? 'CH';
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
      // Get country phones
      this.countryPhoneCodes = this.countryCallingCodesService.getCountryCallingCodes();

      // Guess user's country from Intl.DateTimeFormat
      const defaultCountryCallingCode =
        (getCountryForTimezone(Intl.DateTimeFormat().resolvedOptions().timeZone)?.id as CountryCode) ?? 'CH';
      this.selected = this.countryPhoneCodes.find((code) => code.countryCode === defaultCountryCallingCode);

      // Set current calling country code and format phone if possible
      if (this.valueControl?.value) {
        try {
          const phoneNumber = parsePhoneNumber(this.valueControl?.value ?? '', defaultCountryCallingCode);
          if (phoneNumber) {
            const formatted = phoneNumber.formatInternational();
            this.selected = this.countryPhoneCodes.find((code) => code.countryCode === phoneNumber.country);
            if (formatted !== this.valueControl.value) {
              this.valueControl.setValue(formatted);
            }
          }
        } catch (error) {
          console.error(error);
        }
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
