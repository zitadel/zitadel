import { Component, Inject } from '@angular/core';
import { AbstractControl, UntypedFormBuilder, UntypedFormGroup } from '@angular/forms';
import {
  MatLegacyDialogRef as MatDialogRef,
  MAT_LEGACY_DIALOG_DATA as MAT_DIALOG_DATA,
} from '@angular/material/legacy-dialog';
import { Duration } from 'google-protobuf/google/protobuf/duration_pb';
import { requiredValidator } from 'src/app/modules/form-field/validators/validators';
import { UpdateSecretGeneratorRequest } from 'src/app/proto/generated/zitadel/admin_pb';
import { SecretGeneratorType } from 'src/app/proto/generated/zitadel/settings_pb';

@Component({
  selector: 'cnsl-dialog-add-secret-generator',
  templateUrl: './dialog-add-secret-generator.component.html',
  styleUrls: ['./dialog-add-secret-generator.component.scss'],
})
export class DialogAddSecretGeneratorComponent {
  public SecretGeneratorType: any = SecretGeneratorType;
  public availableGenerators: SecretGeneratorType[] = [
    SecretGeneratorType.SECRET_GENERATOR_TYPE_INIT_CODE,
    SecretGeneratorType.SECRET_GENERATOR_TYPE_VERIFY_EMAIL_CODE,
    SecretGeneratorType.SECRET_GENERATOR_TYPE_VERIFY_PHONE_CODE,
    SecretGeneratorType.SECRET_GENERATOR_TYPE_PASSWORD_RESET_CODE,
    SecretGeneratorType.SECRET_GENERATOR_TYPE_PASSWORDLESS_INIT_CODE,
    SecretGeneratorType.SECRET_GENERATOR_TYPE_APP_SECRET,
  ];
  public req: UpdateSecretGeneratorRequest = new UpdateSecretGeneratorRequest();

  public specsForm!: UntypedFormGroup;

  constructor(
    private fb: UntypedFormBuilder,
    public dialogRef: MatDialogRef<DialogAddSecretGeneratorComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any,
  ) {
    this.specsForm = this.fb.group({
      generatorType: [SecretGeneratorType.SECRET_GENERATOR_TYPE_APP_SECRET, [requiredValidator]],
      expiry: [1, [requiredValidator]],
      includeDigits: [true, [requiredValidator]],
      includeLowerLetters: [true, [requiredValidator]],
      includeSymbols: [true, [requiredValidator]],
      includeUpperLetters: [true, [requiredValidator]],
      length: [6, [requiredValidator]],
    });

    this.generatorType?.setValue(data.type);
  }

  public closeDialog(): void {
    this.dialogRef.close();
  }

  public closeDialogWithRequest(): void {
    this.req.setGeneratorType(this.generatorType?.value);

    const expiry = new Duration().setSeconds((this.expiry?.value ?? 1) * 60 * 60);

    this.req.setExpiry(expiry);
    this.req.setIncludeDigits(this.includeDigits?.value);
    this.req.setIncludeLowerLetters(this.includeLowerLetters?.value);
    this.req.setIncludeSymbols(this.includeSymbols?.value);
    this.req.setIncludeUpperLetters(this.includeUpperLetters?.value);
    this.req.setLength(this.length?.value);

    this.dialogRef.close(this.req);
  }

  public get generatorType(): AbstractControl | null {
    return this.specsForm.get('generatorType');
  }

  public get expiry(): AbstractControl | null {
    return this.specsForm.get('expiry');
  }

  public get includeDigits(): AbstractControl | null {
    return this.specsForm.get('includeDigits');
  }

  public get includeLowerLetters(): AbstractControl | null {
    return this.specsForm.get('includeLowerLetters');
  }

  public get includeSymbols(): AbstractControl | null {
    return this.specsForm.get('includeSymbols');
  }

  public get includeUpperLetters(): AbstractControl | null {
    return this.specsForm.get('includeUpperLetters');
  }

  public get length(): AbstractControl | null {
    return this.specsForm.get('length');
  }
}
