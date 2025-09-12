import { Component, Input, OnInit } from '@angular/core';
import { AbstractControl, UntypedFormBuilder, UntypedFormGroup } from '@angular/forms';
import { Duration } from 'google-protobuf/google/protobuf/duration_pb';
import { requiredValidator } from 'src/app/modules/form-field/validators/validators';
import { GetSecretGeneratorRequest, UpdateSecretGeneratorRequest } from 'src/app/proto/generated/zitadel/admin_pb';
import { SecretGeneratorType } from 'src/app/proto/generated/zitadel/settings_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ToastService } from 'src/app/services/toast.service';

const MIN_EXPIRATION_IN_MINUTES = 5;
const MIN_LENGTH = 6;

@Component({
  selector: 'cnsl-secret-generator-card',
  templateUrl: './secret-generator-card.component.html',
  styleUrls: ['./secret-generator-card.component.scss'],
})
export class SecretGeneratorCardComponent implements OnInit {
  @Input({ required: true }) generatorType!: SecretGeneratorType;

  public specsForm!: UntypedFormGroup;
  public loading: boolean = false;

  ngOnInit() {
    this.fetchData();
  }

  constructor(
    private fb: UntypedFormBuilder,
    private service: AdminService,
    private toast: ToastService,
  ) {
    this.specsForm = this.fb.group({
      expiry: [MIN_EXPIRATION_IN_MINUTES, [requiredValidator]],
      length: [MIN_LENGTH, [requiredValidator]],
      includeDigits: [false, [requiredValidator]],
      includeSymbols: [false, [requiredValidator]],
      includeLowerLetters: [false, [requiredValidator]],
      includeUpperLetters: [false, [requiredValidator]],
    });
  }

  private fetchData(): void {
    this.loading = true;
    const req = new GetSecretGeneratorRequest();
    req.setGeneratorType(this.generatorType);

    this.service
      .getSecretGenerator(req)
      .then((resp) => {
        let generator = resp.secretGenerator;
        if (generator) {
          this.specsForm.patchValue({
            length: generator.length,
            includeDigits: generator.includeDigits,
            includeSymbols: generator.includeSymbols,
            includeLowerLetters: generator.includeLowerLetters,
            includeUpperLetters: generator.includeUpperLetters,
          });

          if (generator.expiry !== undefined) {
            this.specsForm.patchValue({ expiry: this.durationToMinutes(generator.expiry) });
          }
          this.specsForm.markAsPristine();
          this.loading = false;
        }
      })
      .catch((error) => {
        this.toast.showError(error);
        this.loading = false;
      });
  }

  public saveSecretGenerator() {
    const req = new UpdateSecretGeneratorRequest();

    req.setGeneratorType(this.generatorType);
    req.setExpiry(this.minutesToDuration(this.expiry?.value));
    req.setIncludeDigits(this.includeDigits?.value);
    req.setIncludeLowerLetters(this.includeLowerLetters?.value);
    req.setIncludeSymbols(this.includeSymbols?.value);
    req.setIncludeUpperLetters(this.includeUpperLetters?.value);
    req.setLength(this.length?.value);

    this.loading = true;
    this.service
      .updateSecretGenerator(req)
      .then(() => {
        this.toast.showInfo('SETTING.SECRETS.UPDATED', true);
        this.fetchData();
      })
      .catch((error) => {
        this.toast.showError(error);
        this.loading = false;
      });
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

  private durationToMinutes(duration: Duration.AsObject): number {
    if (duration.seconds === 0) {
      return 0;
    }
    return (duration.seconds + duration.nanos / 1000000000) / 60;
  }

  private minutesToDuration(minutes: number): Duration {
    const exp = minutes * 60;
    const sec = Math.floor(exp);
    const nanos = Math.round((exp - sec) * 1000000000);
    return new Duration().setSeconds(sec).setNanos(nanos);
  }
}
