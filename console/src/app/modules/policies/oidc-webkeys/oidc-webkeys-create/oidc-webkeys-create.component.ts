import { ChangeDetectionStrategy, Component, EventEmitter, Input, Output, signal, WritableSignal } from '@angular/core';
import { FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { WebKey } from '@zitadel/proto/zitadel/webkey/v2beta/key_pb';
import { ReplaySubject } from 'rxjs';
import { RSAHasher, RSABits, ECDSACurve } from '@zitadel/proto/zitadel/webkey/v2beta/key_pb';

type RawValue<T extends FormGroup> = ReturnType<T['getRawValue']>;

@Component({
  selector: 'cnsl-oidc-webkeys-create',
  templateUrl: './oidc-webkeys-create.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class OidcWebKeysCreateComponent {
  protected readonly keyType: WritableSignal<NonNullable<WebKey['key']['case']>> = signal('rsa');
  protected readonly RSAHasher = RSAHasher;
  protected readonly RSABits = RSABits;
  protected readonly ECDSACurve = ECDSACurve;
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
      bits: new FormControl<RSABits>(RSABits.RSA_BITS_2048, {
        nonNullable: true,
        validators: [Validators.required],
      }),
      hasher: new FormControl<RSAHasher>(RSAHasher.RSA_HASHER_SHA256, {
        nonNullable: true,
        validators: [Validators.required],
      }),
    });
  }

  private buildEcdsaForm() {
    return this.fb.group({
      curve: new FormControl<ECDSACurve>(ECDSACurve.ECDSA_CURVE_P256, {
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
