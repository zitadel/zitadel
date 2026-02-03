import { Location } from '@angular/common';
import { Component, computed, inject, linkedSignal, Signal, signal } from '@angular/core';
import { FormBuilder, FormControl, Validators } from '@angular/forms';
import { map, switchMap } from 'rxjs';
import { StepperSelectionEvent } from '@angular/cdk/stepper';
import { Options } from 'src/app/proto/generated/zitadel/idp_pb';
import { requiredValidator } from '../form-field/validators/validators';

import { AdminService } from 'src/app/services/admin.service';
import { ToastService } from 'src/app/services/toast.service';
import { ActivatedRoute } from '@angular/router';
import * as SMTPKnownProviders from './known-smtp-providers-settings';
import { toObservable, toSignal } from '@angular/core/rxjs-interop';
import { filter, startWith } from 'rxjs/operators';
import { UserService } from '../../services/user.service';
import { injectMutation, injectQuery } from '@tanstack/angular-query-experimental';
import { MatCheckboxChange } from '@angular/material/checkbox';
import { NewAdminService } from '../../services/new-admin.service';
import { MessageInitShape } from '@bufbuild/protobuf';
import { TestEmailProviderSMTPRequestSchema } from '@zitadel/proto/zitadel/admin_pb';
import { ConnectError } from '@connectrpc/connect';

@Component({
  selector: 'cnsl-smtp-provider',
  templateUrl: './smtp-provider.component.html',
  styleUrls: ['./smtp-provider.scss'],
  standalone: false,
})
export class SMTPProviderComponent {
  protected readonly emailQuery = this.buildEmailQuery();
  protected readonly email = linkedSignal(() => this.emailQuery.data() ?? 'test@example.com');

  protected readonly state = this.buildState();
  protected readonly auth = this.buildAuth();

  protected readonly secondFormGroup;

  protected testEmailConfiguration = this.testEmailConfigurationMutation();

  public options: Options = new Options().setIsCreationAllowed(true).setIsLinkingAllowed(true);

  protected readonly currentCreateStep = signal(1);
  protected readonly location = inject(Location);
  protected readonly newAdminService = inject(NewAdminService);

  public senderEmailPlaceholder: Signal<string>;

  constructor(
    private service: AdminService,
    private toast: ToastService,
  ) {
    this.senderEmailPlaceholder = computed(() => {
      const state = this.state();

      if ('senderEmailPlaceholder' in state.providerDefaults) {
        return state.providerDefaults.senderEmailPlaceholder;
      }
      return 'sender@example.com';
    });

    this.secondFormGroup = inject(FormBuilder).group({
      senderAddress: new FormControl('', { nonNullable: true, validators: [requiredValidator] }),
      senderName: new FormControl('', { nonNullable: true, validators: [requiredValidator] }),
      replyToAddress: new FormControl('', { nonNullable: true }),
    });
  }

  private buildEmailQuery() {
    const userQueryOptions = inject(UserService).userQueryOptions();
    return injectQuery(() => ({
      ...userQueryOptions,
      select: (user) => {
        if (user?.type.case !== 'human') {
          return '';
        }
        return user.type.value.email?.email ?? '';
      },
    }));
  }

  private getProviderDefaults() {
    const providerKey$ = inject(ActivatedRoute).paramMap.pipe(
      map((params) => params.get('provider')),
      filter(Boolean),
    );

    const providerKeySignal = toSignal(providerKey$, { requireSync: true });

    return computed(() => {
      const providerKey = providerKeySignal();
      if (providerKey in SMTPKnownProviders) {
        return SMTPKnownProviders[providerKey as keyof typeof SMTPKnownProviders];
      }
      throw new Error('Unknown SMTP provider key: ' + providerKey);
    });
  }

  private buildState() {
    const fb = inject(FormBuilder);
    const providerDefaultsSignal = this.getProviderDefaults();

    const hostnameValidator = Validators.pattern(/.+:[0-9]+/);

    return computed(() => {
      const providerDefaults = providerDefaultsSignal();

      const baseControls = {
        description: new FormControl<string>(providerDefaults.name, { nonNullable: true, validators: [requiredValidator] }),
        xoauth2: new FormControl<boolean>('auth' in providerDefaults && providerDefaults.auth.case === 'xoauth2', {
          nonNullable: true,
        }),
        tls: new FormControl<true>(true, { nonNullable: true }),
      };

      if ('regions' in providerDefaults) {
        const form = fb.group({
          ...baseControls,
          region: new FormControl<string>('', { nonNullable: true, validators: [requiredValidator] }),
        });

        form.controls.tls.disable();
        form.controls.xoauth2.disable();

        return {
          case: 'region',
          form,
          providerDefaults,
        } as const;
      }

      const host = 'host' in providerDefaults ? providerDefaults.host : '';
      const port = 'ports' in providerDefaults ? (`${providerDefaults.ports.encryptedPort}` as const) : '';

      const form = fb.group({
        ...baseControls,
        hostAndPort: new FormControl(`${host}:${port}`, {
          nonNullable: true,
          validators: [hostnameValidator],
        }),
        tls: new FormControl(true, { nonNullable: true }),
      });

      if ('ports' in providerDefaults && !('unencryptedPort' in providerDefaults.ports)) {
        form.controls.tls.disable();
      }

      if ('auth' in providerDefaults && providerDefaults.auth.case === 'xoauth2') {
        form.controls.xoauth2.disable();
      }

      return {
        case: 'host',
        form,
        providerDefaults,
      } as const;
    });
  }

  private buildAuth() {
    const fb = inject(FormBuilder);

    const xoauth2$ = toObservable(this.state).pipe(
      switchMap(({ form }) => form.controls.xoauth2.valueChanges.pipe(startWith(form.controls.xoauth2.value))),
    );
    const xoauth2Signal = toSignal(xoauth2$, { initialValue: this.state().form.controls.xoauth2.value });

    return computed(() => {
      const state = this.state();
      const xoauth2 = xoauth2Signal();

      if (!xoauth2) {
        return fb.group({
          user: new FormControl('', { nonNullable: true, validators: [requiredValidator] }),
          password: new FormControl('', { nonNullable: true, validators: [requiredValidator] }),
        });
      }

      const scopes =
        'auth' in state.providerDefaults && state.providerDefaults.auth.case === 'xoauth2'
          ? state.providerDefaults.auth.scopes
          : '';

      return fb.group({
        tokenEndpoint: new FormControl('', { nonNullable: true, validators: [requiredValidator] }),
        scopes: new FormControl(scopes, { nonNullable: true, validators: [requiredValidator] }),
        clientId: new FormControl('', { nonNullable: true, validators: [requiredValidator] }),
        clientSecret: new FormControl('', { nonNullable: true, validators: [requiredValidator] }),
      });
    });
  }

  protected toggleTls(event: MatCheckboxChange) {
    const state = this.state();

    if (state.case !== 'host') {
      return;
    }

    if (!('ports' in state.providerDefaults) || !('unencryptedPort' in state.providerDefaults.ports)) {
      return;
    }

    const port = event.checked ? state.providerDefaults.ports.encryptedPort : state.providerDefaults.ports.unencryptedPort;
    state.form.controls.hostAndPort.setValue(`${state.providerDefaults.host}:${port}`);
  }

  public changeStep(event: StepperSelectionEvent): void {
    this.currentCreateStep.set(event.selectedIndex + 1);
  }

  // public toggleTLS(event: MatCheckboxChange) {
  //   const
  //   if (!('host' in this.providerDefaultSetting)) {
  //     return;
  //   }
  //   this.hostAndPort?.setValue(
  //     `${this.providerDefaultSetting?.host}:${
  //       event.checked ? this.providerDefaultSetting?.encryptedPort : this.providerDefaultSetting?.unencryptedPort
  //     }`,
  //   );
  // }

  // private fetchData(id: string): void {
  //   this.smtpLoading = true;
  //   this.service
  //     .getSMTPConfigById(id)
  //     .then((data) => {
  //       this.smtpLoading = false;
  //       if (data.smtpConfig) {
  //         this.isActive = data.smtpConfig.state === SMTPConfigState.SMTP_CONFIG_ACTIVE;
  //         this.hasSMTPConfig = true;
  //         this.firstFormGroup.patchValue({
  //           ['description']: data.smtpConfig.description,
  //           ['tls']: data.smtpConfig.tls,
  //           ['hostAndPort']: data.smtpConfig.host,
  //           ['user']: data.smtpConfig.user,
  //         });
  //         this.secondFormGroup.patchValue({
  //           ['senderAddress']: data.smtpConfig.senderAddress,
  //           ['senderName']: data.smtpConfig.senderName,
  //           ['replyToAddress']: data.smtpConfig.replyToAddress,
  //         });
  //       }
  //     })
  //     .catch((error) => {
  //       this.smtpLoading = false;
  //       if (error && error.code === 5) {
  //         this.hasSMTPConfig = false;
  //       }
  //     });
  // }

  // private updateData(): Promise<UpdateSMTPConfigResponse.AsObject | AddSMTPConfigResponse.AsObject> {
  //   if (this.hasSMTPConfig) {
  //     const req = new UpdateSMTPConfigRequest();
  //     req.setId(this.id);
  //     req.setDescription(this.description?.value || '');
  //     req.setTls(this.tls?.value ?? false);
  //
  //     if (this.hostAndPort && this.hostAndPort.value) {
  //       req.setHost(this.hostAndPort.value);
  //     }
  //     if (this.user && this.user.value) {
  //       req.setUser(this.user.value);
  //     }
  //     if (this.password && this.password.value) {
  //       req.setPassword(this.password.value);
  //     }
  //     if (this.senderAddress && this.senderAddress.value) {
  //       req.setSenderAddress(this.senderAddress.value);
  //     }
  //     if (this.senderName && this.senderName.value) {
  //       req.setSenderName(this.senderName.value);
  //     }
  //     if (this.replyToAddress && this.replyToAddress.value) {
  //       req.setReplyToAddress(this.replyToAddress.value);
  //     }
  //     return this.service.updateSMTPConfig(req);
  //   } else {
  //     const req = new AddSMTPConfigRequest();
  //     req.setDescription(this.description?.value ?? '');
  //     req.setHost(this.hostAndPort?.value ?? '');
  //     req.setSenderAddress(this.senderAddress?.value ?? '');
  //     req.setSenderName(this.senderName?.value ?? '');
  //     req.setReplyToAddress(this.replyToAddress?.value ?? '');
  //     req.setTls(this.tls?.value ?? false);
  //     req.setUser(this.user?.value ?? '');
  //     req.setPassword(this.password?.value ?? '');
  //     return this.service.addSMTPConfig(req);
  //   }
  // }

  protected readonly id = '';
  protected isActive = false;
  protected hasSMTPConfig = false;
  protected smtpLoading = false;
  protected resultClass = '';
  protected isLoading = signal(false);

  public activateSMTPConfig() {
    this.service
      .activateSMTPConfig(this.id)
      .then(() => {
        this.toast.showInfo('SMTP.LIST.DIALOG.ACTIVATED', true);
        this.isActive = true;
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public deactivateSMTPConfig() {
    this.service
      .deactivateSMTPConfig(this.id)
      .then(() => {
        this.toast.showInfo('SMTP.LIST.DIALOG.DEACTIVATED', true);
        this.isActive = false;
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  // public savePolicy(stepper: MatStepper): void {
  //   this.updateData()
  //     .then((resp) => {
  //       if (!this.id) {
  //         // This is a new SMTP provider let's get the ID from the addSMTPConfig response
  //         let createResponse = resp as AddSMTPConfigResponse.AsObject;
  //         this.id = createResponse.id;
  //       }
  //
  //       this.toast.showInfo('SETTING.SMTP.SAVED', true);
  //       setTimeout(() => {
  //         stepper.next();
  //       }, 2000);
  //     })
  //     .catch((error: unknown) => {
  //       if (`${error}`.includes('No changes')) {
  //         this.toast.showInfo('SETTING.SMTP.NOCHANGES', true);
  //         setTimeout(() => {
  //           stepper.next();
  //         }, 2000);
  //       } else {
  //         this.toast.showError(error);
  //       }
  //     });
  // }
  //
  private testEmailConfigurationMutation() {
    const req = computed(() => {
      const state = this.state();
      const auth = this.auth().getRawValue();

      const { tls } = state.form.getRawValue();
      const { senderAddress, senderName } = this.secondFormGroup.getRawValue();

      const host =
        state.case === 'host'
          ? state.form.controls.hostAndPort.value
          : `${state.form.controls.region.value}:${state.providerDefaults.ports.encryptedPort}`;

      const Auth: MessageInitShape<typeof TestEmailProviderSMTPRequestSchema>['Auth'] =
        'tokenEndpoint' in auth
          ? {
              case: 'xoauth2',
              value: {
                tokenEndpoint: auth.tokenEndpoint,
                scopes: auth.scopes.split(','),
                OAuth2Type: {
                  case: 'clientCredentials',
                  value: {
                    clientId: auth.clientId,
                    clientSecret: auth.clientSecret,
                  },
                },
              },
            }
          : auth.password
            ? { case: 'plain', value: { password: auth.password } }
            : { case: 'none', value: {} };

      return {
        senderAddress,
        senderName,
        host,
        user: 'user' in auth ? auth.user : undefined,
        tls,
        receiverAddress: this.email(),
        Auth,
      };
    });

    return injectMutation<unknown, ConnectError>(() => ({
      mutationFn: async () => {
        return this.newAdminService.testEmailProviderSMTP(req());
      },
    }));
  }
  //   this.isLoading.set(true);
  //
  //   const req = new TestSMTPConfigRequest();
  //   req.setSenderAddress(this.senderAddress?.value ?? '');
  //   req.setSenderName(this.senderName?.value ?? '');
  //   req.setHost(this.hostAndPort?.value ?? '');
  //   req.setUser(this.user?.value);
  //   req.setPassword(this.password?.value ?? '');
  //   req.setTls(this.tls?.value ?? false);
  //   req.setId(this.id ?? '');
  //   req.setReceiverAddress(this.email ?? '');
  //
  //   this.service
  //     .testSMTPConfig(req)
  //     .then(() => {
  //       this.resultClass = 'test-success';
  //       this.isLoading.set(false);
  //       this.translate
  //         .get('SMTP.CREATE.STEPS.TEST.RESULT')
  //         .pipe(take(1))
  //         .subscribe((msg) => {
  //           this.testResult = msg;
  //         });
  //     })
  //     .catch((error) => {
  //       this.resultClass = 'test-error';
  //       this.isLoading.set(false);
  //       this.testResult = error;
  //     });
  // }
}
