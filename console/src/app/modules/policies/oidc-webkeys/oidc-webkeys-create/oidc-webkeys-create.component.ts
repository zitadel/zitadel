import { ChangeDetectionStrategy, Component, EventEmitter, Input, Output, signal, WritableSignal } from '@angular/core';
import {
  WebKeyECDSAConfig_ECDSACurve,
  WebKeyRSAConfig_RSABits,
  WebKeyRSAConfig_RSAHasher,
} from '@zitadel/proto/zitadel/resources/webkey/v3alpha/config_pb';
import { WebKey } from '@zitadel/proto/zitadel/resources/webkey/v3alpha/key_pb';
import { FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { ReplaySubject } from 'rxjs';

type RawValue<T extends FormGroup> = ReturnType<T['getRawValue']>;

@Component({
  selector: 'cnsl-oidc-webkeys-create',
  templateUrl: './oidc-webkeys-create.component.html',
  styleUrls: ['./oidc-webkeys-create.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class OidcWebKeysCreateComponent {
  protected readonly keyType: WritableSignal<NonNullable<WebKey['config']['case']>> = signal('rsa');
  protected readonly RSAHasher = WebKeyRSAConfig_RSAHasher;
  protected readonly RSABits = WebKeyRSAConfig_RSABits;
  protected readonly ECDSACurve = WebKeyECDSAConfig_ECDSACurve;
  protected readonly Number = Number;
  protected readonly rsaForm = this.buildRsaForm();
  protected readonly ecdsaForm = this.buildEcdsaForm();
  protected readonly loading$ = new ReplaySubject<boolean>();

  @Output()
  public readonly ngSubmit = new EventEmitter<RawValue<typeof this.rsaForm> | RawValue<typeof this.ecdsaForm> | void>();

  @Input()
  public set loading(loading: boolean) {
    this.loading$.next(loading);
  }

  constructor(private readonly fb: FormBuilder) {}

  private buildRsaForm() {
    return this.fb.group({
      bits: new FormControl<WebKeyRSAConfig_RSABits>(WebKeyRSAConfig_RSABits.RSA_BITS_2048, {
        nonNullable: true,
        validators: [Validators.required],
      }),
      hasher: new FormControl<WebKeyRSAConfig_RSAHasher>(WebKeyRSAConfig_RSAHasher.RSA_HASHER_SHA256, {
        nonNullable: true,
        validators: [Validators.required],
      }),
    });
  }

  private buildEcdsaForm() {
    return this.fb.group({
      curve: new FormControl<WebKeyECDSAConfig_ECDSACurve>(WebKeyECDSAConfig_ECDSACurve.ECDSA_CURVE_P256, {
        nonNullable: true,
        validators: [Validators.required],
      }),
    });
  }

  protected emitEd25519(event: SubmitEvent) {
    event.preventDefault();
    this.ngSubmit.emit();
  }
}
