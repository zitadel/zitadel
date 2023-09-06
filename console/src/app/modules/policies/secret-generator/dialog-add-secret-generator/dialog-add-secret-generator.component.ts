import { Component, Inject } from '@angular/core';
import { AbstractControl, UntypedFormBuilder, UntypedFormGroup } from '@angular/forms';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { Duration } from 'google-protobuf/google/protobuf/duration_pb';
import { requiredValidator } from 'src/app/modules/form-field/validators/validators';
import { UpdateSecretGeneratorRequest } from 'src/app/proto/generated/zitadel/admin_pb';

@Component({
  selector: 'cnsl-dialog-add-secret-generator',
  templateUrl: './dialog-add-secret-generator.component.html',
  styleUrls: ['./dialog-add-secret-generator.component.scss'],
})
export class DialogAddSecretGeneratorComponent {
  public req: UpdateSecretGeneratorRequest = new UpdateSecretGeneratorRequest();

  public specsForm!: UntypedFormGroup;

  constructor(
    private fb: UntypedFormBuilder,
    public dialogRef: MatDialogRef<DialogAddSecretGeneratorComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any,
  ) {
    let exp = 1;
    if (data.config?.expiry !== undefined) {
      exp = this.durationToHour(data.config?.expiry);
    }
    this.specsForm = this.fb.group({
      generatorType: [data.type, [requiredValidator]],
      expiry: [exp, [requiredValidator]],
      length: [data.config?.length ?? 6, [requiredValidator]],
      includeDigits: [data.config?.includeDigits ?? true, [requiredValidator]],
      includeSymbols: [data.config?.includeSymbols ?? true, [requiredValidator]],
      includeLowerLetters: [data.config?.includeLowerLetters ?? true, [requiredValidator]],
      includeUpperLetters: [data.config?.includeUpperLetters ?? true, [requiredValidator]],
    });
  }

  public closeDialog(): void {
    this.dialogRef.close();
  }

  public closeDialogWithRequest(): void {
    this.req.setGeneratorType(this.generatorType?.value);
    this.req.setExpiry(this.hourToDuration(this.expiry?.value));
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

  private durationToHour(duration: Duration.AsObject): number {
    if (duration.seconds === 0) {
      return 0;
    }
    return (duration.seconds + duration.nanos / 1000000) / 3600;
  }

  private hourToDuration(hour: number): Duration {
    const exp = hour * 60 * 60;
    const sec = Math.floor(exp);
    const nanos = Math.round((exp - sec) * 1000000);
    return new Duration().setSeconds(sec).setNanos(nanos);
  }
}
