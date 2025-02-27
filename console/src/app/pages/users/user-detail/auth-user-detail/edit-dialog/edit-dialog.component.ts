import { Component, Inject, OnInit } from '@angular/core';
import { FormGroup, UntypedFormControl, UntypedFormGroup, ValidatorFn } from '@angular/forms';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { debounceTime } from 'rxjs';
import { requiredValidator } from 'src/app/modules/form-field/validators/validators';
import { CountryCallingCodesService, CountryPhoneCode } from 'src/app/services/country-calling-codes.service';
import { formatPhone } from 'src/app/utils/formatPhone';

export enum EditDialogType {
  PHONE = 1,
  EMAIL = 2,
}

export type EditDialogData = {
  confirmKey: 'ACTIONS.SAVE' | 'ACTIONS.CHANGE';
  cancelKey: 'ACTIONS.CANCEL';
  labelKey: 'ACTIONS.NEWVALUE';
  titleKey: 'USER.LOGINMETHODS.EMAIL.EDITTITLE';
  descriptionKey: 'USER.LOGINMETHODS.EMAIL.EDITDESC';
  isVerifiedTextKey?: 'USER.LOGINMETHODS.EMAIL.ISVERIFIED';
  isVerifiedTextDescKey?: 'USER.LOGINMETHODS.EMAIL.ISVERIFIEDDESC';
  value: string | undefined;
  type: EditDialogType;
  validator?: ValidatorFn;
};

export type EditDialogResult = {
  value?: string;
  isVerified: boolean;
};

@Component({
  selector: 'cnsl-edit-dialog',
  templateUrl: './edit-dialog.component.html',
  styleUrls: ['./edit-dialog.component.scss'],
})
export class EditDialogComponent implements OnInit {
  public controlKey = 'editingField';
  public isPhone: boolean = false;
  public isVerified: boolean = false;
  public phoneCountry: string = 'US';
  public dialogForm!: UntypedFormGroup;
  public EditDialogType: any = EditDialogType;
  public selected: CountryPhoneCode | undefined = {
    countryCallingCode: '1',
    countryCode: 'US',
    countryName: 'United States of America',
  };
  public countryPhoneCodes: CountryPhoneCode[] = [];
  constructor(
    public dialogRef: MatDialogRef<EditDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: EditDialogData,
    private countryCallingCodesService: CountryCallingCodesService,
  ) {
    if (data.type === EditDialogType.PHONE) {
      this.isPhone = true;
    }
    this.dialogForm = new FormGroup({
      [this.controlKey]: new UntypedFormControl(data.value, data.validator || requiredValidator),
    });

    if (this.isPhone) {
      this.ctrl?.valueChanges.pipe(debounceTime(200)).subscribe((value: string) => {
        const phoneNumber = formatPhone(value);
        if (phoneNumber) {
          this.selected = this.countryPhoneCodes.find((code) => code.countryCode === phoneNumber.country);
          this.ctrl?.setValue(phoneNumber.phone);
        }
      });
    }
  }

  public setCountryCallingCode(): void {
    let value = (this.dialogForm.controls[this.controlKey]?.value as string) || '';
    this.countryPhoneCodes.forEach((code) => (value = value.replace(`+${code.countryCallingCode}`, '')));
    value = value.trim();
    this.dialogForm.controls[this.controlKey]?.setValue('+' + this.selected?.countryCallingCode + ' ' + value);
  }

  ngOnInit(): void {
    if (this.isPhone) {
      // Get country phone codes and set selected flag to guessed country or default country
      this.countryPhoneCodes = this.countryCallingCodesService.getCountryCallingCodes();
      const phoneNumber = formatPhone(this.dialogForm.controls[this.controlKey]?.value);
      if (phoneNumber) {
        this.selected = this.countryPhoneCodes.find((code) => code.countryCode === phoneNumber.country);
        this.dialogForm.controls[this.controlKey].setValue(phoneNumber.phone);
      }
    }
  }

  closeDialog(): void {
    this.dialogRef.close();
  }

  closeDialogWithValue(): void {
    this.dialogRef.close({ value: this.dialogForm.controls[this.controlKey].value, isVerified: this.isVerified });
  }

  public get ctrl() {
    return this.dialogForm.get(this.controlKey);
  }

  public compareCountries(i1: CountryPhoneCode, i2: CountryPhoneCode) {
    return (
      i1 &&
      i2 &&
      i1.countryCallingCode === i2.countryCallingCode &&
      i1.countryCode == i2.countryCode &&
      i1.countryName == i2.countryName
    );
  }
}
