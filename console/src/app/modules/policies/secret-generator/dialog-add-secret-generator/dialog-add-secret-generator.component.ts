import { Component, Inject } from '@angular/core';
import { AbstractControl, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { AddSMSProviderTwilioRequest } from 'src/app/proto/generated/zitadel/admin_pb';
import { SecretGeneratorType } from 'src/app/proto/generated/zitadel/settings_pb';

@Component({
  selector: 'cnsl-dialog-add-secret-generator',
  templateUrl: './dialog-add-secret-generator.component.html',
  styleUrls: ['./dialog-add-secret-generator.component.scss'],
})
export class DialogAddSecretGeneratorComponent {
  public SecretGeneratorType: any = SecretGeneratorType;
  public availableGenerators: SecretGeneratorType[] = [
    SecretGeneratorType.SECRET_GENERATOR_TYPE_APP_SECRET,
    SecretGeneratorType.SECRET_GENERATOR_TYPE_APP_SECRET,
    SecretGeneratorType.SECRET_GENERATOR_TYPE_APP_SECRET,
    SecretGeneratorType.SECRET_GENERATOR_TYPE_APP_SECRET,
    SecretGeneratorType.SECRET_GENERATOR_TYPE_APP_SECRET,
  ];
  public generatorType: SecretGeneratorType = SecretGeneratorType.SECRET_GENERATOR_TYPE_APP_SECRET;
  public req: AddSMSProviderTwilioRequest = new AddSMSProviderTwilioRequest();

  public specsForm!: FormGroup;

  constructor(
    private fb: FormBuilder,
    public dialogRef: MatDialogRef<DialogAddSecretGeneratorComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any,
  ) {
    this.specsForm = this.fb.group({
      sid: ['', [Validators.required]],
      token: ['', [Validators.required]],
      senderNumber: ['', [Validators.required]],
    });
  }

  public closeDialog(): void {
    this.dialogRef.close();
  }

  public closeDialogWithRequest(): void {
    this.req.setSid(this.sid?.value);
    this.req.setToken(this.token?.value);
    this.req.setSenderNumber(this.senderNumber?.value);

    this.dialogRef.close(this.req);
  }

  public get senderNumber(): AbstractControl | null {
    return this.specsForm.get('senderNumber');
  }

  public get token(): AbstractControl | null {
    return this.specsForm.get('token');
  }

  public get sid(): AbstractControl | null {
    return this.specsForm.get('sid');
  }
}
