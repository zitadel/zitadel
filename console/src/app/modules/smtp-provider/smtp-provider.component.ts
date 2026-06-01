import { Location } from '@angular/common';
import { Component, computed, effect, inject, linkedSignal, Signal, viewChild } from '@angular/core';
import { FormBuilder, FormControl, Validators } from '@angular/forms';
import { requiredValidator } from '../form-field/validators/validators';
import { ActivatedRoute, Router } from '@angular/router';
import * as SMTPKnownProviders from './known-smtp-providers-settings';
import { toObservable, toSignal } from '@angular/core/rxjs-interop';
import { UserService } from '../../services/user.service';
import { injectMutation, injectQuery, QueryFunction } from '@tanstack/angular-query-experimental';
import { NewAdminService } from '../../services/new-admin.service';
import { MessageInitShape } from '@bufbuild/protobuf';
import { AddEmailProviderSMTPRequestSchema, GetEmailProviderByIdResponse } from '@zitadel/proto/zitadel/admin_pb';
import { EMPTY, map, switchMap } from 'rxjs';
import { filter, startWith } from 'rxjs/operators';
import { EmailProviderState } from '@zitadel/proto/zitadel/settings_pb';
import { ToastService } from '../../services/toast.service';
import { MatStepper } from '@angular/material/stepper';
import { ConnectError } from '@zitadel/client';

type Provider = (typeof SMTPKnownProviders)[keyof typeof SMTPKnownProviders];
type State = SMTPProviderComponent['state'] extends Signal<infer T> ? NonNullable<T> : never;

@Component({
  selector: 'cnsl-smtp-provider',
  templateUrl: './smtp-provider.component.html',
  styleUrls: ['./smtp-provider.scss'],
  standalone: false,
})
export class SMTPProviderComponent {
  protected readonly EmailProviderState = EmailProviderState;

  protected readonly emailQuery = this.buildEmailQuery();
  protected readonly email = linkedSignal(() => this.emailQuery.data() ?? 'test@example.com');

  protected readonly newAdminService = inject(NewAdminService);
  private readonly toast = inject(ToastService);

  protected readonly router = inject(Router);
  protected readonly location = inject(Location);
  private readonly activatedRoute = inject(ActivatedRoute);
  private readonly fb = inject(FormBuilder);

  protected readonly configOrDefaultsQuery: ReturnType<typeof this.buildConfigOrDefaultsQuery>;
  protected readonly state: ReturnType<typeof this.buildState>;
  protected readonly updateDataMutation: ReturnType<typeof this.buildUpdateDataMutation>;
  protected readonly testEmailConfigurationMutation: ReturnType<typeof this.buildTestEmailConfigurationMutation>;

  protected readonly stepper = viewChild(MatStepper);
  protected readonly preselectedStep: ReturnType<typeof this.getPreselectedStep>;

  constructor() {
    this.configOrDefaultsQuery = this.buildConfigOrDefaultsQuery();
    this.state = this.buildState(this.configOrDefaultsQuery);
    this.updateDataMutation = this.buildUpdateDataMutation(this.configOrDefaultsQuery, this.state);
    this.testEmailConfigurationMutation = this.buildTestEmailConfigurationMutation(this.state);

    effect(() => {
      const error = this.configOrDefaultsQuery.error();
      if (error) {
        this.toast.showError(error);
      }
    });

    this.preselectedStep = this.getPreselectedStep(this.activatedRoute);
  }

  private getPreselectedStep(activatedRoute: ActivatedRoute) {
    const paramMapSignal = toSignal(activatedRoute.paramMap, { requireSync: true });

    return computed(() => {
      const paramMap = paramMapSignal();
      const step = paramMap.get('step');
      if (!step) {
        return 0;
      }

      return Number(step);
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
  private buildState(configOrDefaultsQuery: typeof this.configOrDefaultsQuery) {
    const stateSignal = computed(() => {
      const configOrDefaults = configOrDefaultsQuery.data();
      if (!configOrDefaults) {
        return undefined;
      }

      return configOrDefaults.case === 'defaults'
        ? ({ ...this.buildFormFromDefaults(configOrDefaults.defaults) } as const)
        : ({ ...this.buildFormFromConfig(configOrDefaults.config) } as const);
    });

    const authFormSignal = this.buildAuthForm(stateSignal);

    return computed(() => {
      const state = stateSignal();
      const authForm = authFormSignal();

      if (!state || !authForm) {
        return undefined;
      }

      return {
        ...state,
        authForm,
      } as const;
    });
  }

  private readonly hostnameValidator = Validators.pattern(/.+:[0-9]+/);

  private buildFormFromConfig(config: ReturnType<typeof this.getConfig>) {
    const mainForm = this.fb.group({
      description: new FormControl<string>(config.description, {
        nonNullable: true,
        validators: [requiredValidator],
      }),
      user: new FormControl<string>(config.config.value.user, { nonNullable: true, validators: [requiredValidator] }),
      host: new FormControl(config.config.value.host, {
        nonNullable: true,
        validators: [this.hostnameValidator],
      }),
      tls: new FormControl(config.config.value.tls, { nonNullable: true }),
      xoauth2: new FormControl<boolean>(config.config.value.Auth.case === 'xoauth2', { nonNullable: true }),
    });

    mainForm.controls.xoauth2.disable();

    const senderForm = this.buildSenderForm(config.config.value);

    return {
      mainForm,
      senderForm,
      senderEmailPlaceholder: 'sender@example.com',
      config,
    };
  }

  private buildFormFromDefaults(defaults?: Provider):
    | {
        mainForm: typeof mainForm;
        senderForm: typeof senderForm;
        senderEmailPlaceholder: string;
      }
    | {
        mainForm: typeof mainForm;
        senderForm: typeof senderForm;
        senderEmailPlaceholder: string;
        defaults: Provider;
      } {
    const mainForm = this.fb.group({
      description: new FormControl<string>(defaults?.description ?? '', {
        nonNullable: true,
        validators: [requiredValidator],
      }),
      user: new FormControl<string>(defaults?.user.value ?? '', { nonNullable: true, validators: [requiredValidator] }),
      host: new FormControl(defaults?.host ?? '', {
        nonNullable: true,
        validators: [this.hostnameValidator],
      }),
      tls: new FormControl<boolean>(true, { nonNullable: true }),
      xoauth2: new FormControl<boolean>(defaults?.auth.case === 'xoauth2', { nonNullable: true }),
    });

    if (defaults) {
      mainForm.controls.tls.disable();
      mainForm.controls.xoauth2.disable();
    }

    const senderForm = this.buildSenderForm();
    const senderEmailPlaceholder =
      defaults && 'senderEmailPlaceholder' in defaults ? defaults.senderEmailPlaceholder : 'sender@example.com';

    return defaults
      ? {
          mainForm,
          senderForm,
          senderEmailPlaceholder,
          defaults,
        }
      : {
          mainForm,
          senderForm,
          senderEmailPlaceholder,
        };
  }

  private buildSenderForm(config?: { senderAddress: string; senderName: string; replyToAddress: string }) {
    return this.fb.group({
      senderAddress: new FormControl(config?.senderAddress ?? '', { nonNullable: true, validators: [requiredValidator] }),
      senderName: new FormControl(config?.senderName ?? '', { nonNullable: true, validators: [requiredValidator] }),
      replyToAddress: new FormControl(config?.replyToAddress ?? '', { nonNullable: true }),
    });
  }

  private buildAuthForm(
    stateSignal: Signal<
      ReturnType<typeof this.buildFormFromDefaults> | ReturnType<typeof this.buildFormFromConfig> | undefined
    >,
  ) {
    const xoauth2$ = toObservable(stateSignal).pipe(
      switchMap((state) => {
        if (!state) {
          return EMPTY;
        }
        const xoauth2 = state.mainForm.controls.xoauth2;
        return xoauth2.valueChanges.pipe(startWith(xoauth2.value));
      }),
    );

    const xoauth2Signal = toSignal(xoauth2$);

    return computed(() => {
      const state = stateSignal();
      const xoauth2 = xoauth2Signal();
      if (!state || xoauth2 === undefined) {
        return undefined;
      }

      if (!xoauth2) {
        const form = this.fb.group({
          password: new FormControl('', {
            nonNullable: true,
            validators: 'defaults' in state ? [requiredValidator] : [],
          }),
        });

        if ('config' in state) {
          form.controls.password.disable();
        }

        return form;
      }

      const defaultValues =
        'config' in state && state.config.config.value.Auth.case === 'xoauth2'
          ? {
              tokenEndpoint: state.config.config.value.Auth.value.tokenEndpoint,
              scopes: state.config.config.value.Auth.value.scopes.join(','),
              clientId: state.config.config.value.Auth.value.OAuth2Type.value?.clientId ?? '',
            }
          : 'defaults' in state && state.defaults.auth.case === 'xoauth2'
            ? { scopes: state.defaults.auth.scopes }
            : {};

      const form = this.fb.group({
        tokenEndpoint: new FormControl<string>(defaultValues.tokenEndpoint ?? '', {
          nonNullable: true,
          validators: [requiredValidator],
        }),
        scopes: new FormControl<string>(defaultValues.scopes ?? '', { nonNullable: true, validators: [requiredValidator] }),
        clientId: new FormControl<string>(defaultValues.clientId ?? '', {
          nonNullable: true,
          validators: [requiredValidator],
        }),
        clientSecret: new FormControl<string>('', { nonNullable: true, validators: [requiredValidator] }),
      });

      if ('config' in state) {
        form.controls.tokenEndpoint.disable();
        form.controls.scopes.disable();
        form.controls.clientId.disable();
        form.controls.clientSecret.disable();
      }

      return form;
    });
  }

  private buildConfigOrDefaultsQuery() {
    const idOrProvider$ = this.activatedRoute.paramMap.pipe(
      map((params) => params.get('provider')),
      filter(Boolean),
    );

    const idOrProviderSignal = toSignal(idOrProvider$, { requireSync: true });

    return injectQuery(() => {
      const idOrProvider = idOrProviderSignal();

      const query = this.newAdminService.getEmailProviderByIdQueryOptions(idOrProvider);
      const queryKey = query.queryKey as (string | undefined)[];
      const queryFn = query.queryFn as QueryFunction<GetEmailProviderByIdResponse | string>;

      const select = (configOrProvider: GetEmailProviderByIdResponse | string) =>
        typeof configOrProvider === 'string'
          ? ({
              case: 'defaults',
              defaults: SMTPKnownProviders[configOrProvider as keyof typeof SMTPKnownProviders] as
                | (typeof SMTPKnownProviders)[keyof typeof SMTPKnownProviders]
                | undefined,
            } as const)
          : ({ case: 'config', config: this.getConfig(configOrProvider) } as const);

      if (idOrProvider in SMTPKnownProviders || idOrProvider === 'generic') {
        return {
          queryKey,
          queryFn: (async () => idOrProvider) as typeof queryFn,
          gcTime: 0,
          select,
        } as const;
      }

      return {
        queryKey,
        queryFn,
        select,
      } as const;
    });
  }

  private getConfig(resp: GetEmailProviderByIdResponse) {
    if (!resp.config) {
      throw new Error('No SMTP provider config found');
    }

    if (resp.config.config.case !== 'smtp') {
      throw new Error('Email provider config with id ' + resp.config.id + ' is not an SMTP config');
    }

    const config = resp.config.config.value;

    return {
      ...resp.config,
      config: {
        case: 'smtp' as const,
        value: {
          ...config,
          Auth: config.Auth,
        },
      },
    };
  }

  private buildUpdateDataMutation(configOrDefaultsQuery: typeof this.configOrDefaultsQuery, stateSignal: typeof this.state) {
    return injectMutation(() => {
      const state = stateSignal();
      const stepper = this.stepper();

      return {
        mutationFn: () => {
          if (!state) {
            throw new Error('Invalid state');
          }
          return this.updateData(state);
        },
        onSuccess: () => {
          this.toast.showInfo('SETTING.SMTP.SAVED', true);
          stepper?.next();
        },
        onError: (error: ConnectError) => {
          if (!error.message.includes('No changes')) {
            this.toast.showError(error);
            return;
          }

          this.toast.showInfo('SETTING.SMTP.NOCHANGES', true);
          stepper?.next();
        },
        onSettled: () => configOrDefaultsQuery.refetch(),
      };
    });
  }

  private async updateData(state: State) {
    const authValues = state.authForm.getRawValue();

    const { user, tls, host, description } = state.mainForm.getRawValue();
    const { senderAddress, senderName, replyToAddress } = state.senderForm.getRawValue();

    if ('config' in state) {
      return this.newAdminService.updateEmailProviderSMTP({
        id: state.config.id,
        description,
        senderAddress,
        senderName,
        replyToAddress,
        host,
        user,
        tls,
      });
    }

    const Auth: MessageInitShape<typeof AddEmailProviderSMTPRequestSchema>['Auth'] =
      'tokenEndpoint' in authValues
        ? {
            case: 'xoauth2',
            value: {
              tokenEndpoint: authValues.tokenEndpoint,
              scopes: authValues.scopes.replace(/\s/g, '').split(','),
              OAuth2Type: {
                case: 'clientCredentials',
                value: {
                  clientId: authValues.clientId,
                  clientSecret: authValues.clientSecret,
                },
              },
            },
          }
        : authValues.password
          ? { case: 'plain', value: { password: authValues.password } }
          : { case: 'none', value: {} };

    const res = await this.newAdminService.addEmailProviderSMTP({
      senderAddress,
      senderName,
      description,
      replyToAddress,
      host,
      user,
      tls,
      Auth,
    });

    await this.router.navigate(['/instance/smtpprovider', res.id, { step: 3 }], { skipLocationChange: true });

    return res;
  }

  protected async activateSMTPConfig(id: string) {
    try {
      await this.newAdminService.activateSMTPConfig(id);
      this.toast.showInfo('SMTP.LIST.DIALOG.ACTIVATED', true);
      await this.configOrDefaultsQuery.refetch();
    } catch (error) {
      this.toast.showError(error);
    }
  }

  protected async deactivateSMTPConfig(id: string) {
    try {
      await this.newAdminService.deactivateSMTPConfig(id);
      this.toast.showInfo('SMTP.LIST.DIALOG.DEACTIVATED', true);
      await this.configOrDefaultsQuery.refetch();
    } catch (error) {
      this.toast.showError(error);
    }
  }

  protected buildTestEmailConfigurationMutation(stateSignal: typeof this.state) {
    const buildRequest = (state: State, receiverAddress: string) => {
      const authValues = state.authForm.getRawValue();
      const { user, tls, host } = state.mainForm.getRawValue();
      const { senderAddress, senderName } = state.senderForm.getRawValue();

      if ('config' in state) {
        return {
          id: state.config.id,
          senderAddress,
          senderName,
          host,
          user,
          tls,
          receiverAddress,
        };
      }

      const Auth: MessageInitShape<typeof AddEmailProviderSMTPRequestSchema>['Auth'] =
        'tokenEndpoint' in authValues
          ? {
              case: 'xoauth2',
              value: {
                tokenEndpoint: authValues.tokenEndpoint,
                scopes: authValues.scopes.replace(/\s/g, '').split(','),
                OAuth2Type: {
                  case: 'clientCredentials',
                  value: {
                    clientId: authValues.clientId,
                    clientSecret: authValues.clientSecret,
                  },
                },
              },
            }
          : authValues.password
            ? { case: 'plain', value: { password: authValues.password } }
            : { case: 'none', value: {} };

      return {
        senderAddress,
        senderName,
        host,
        user,
        tls,
        Auth,
        receiverAddress,
      };
    };

    return injectMutation(() => {
      const state = stateSignal();
      const email = this.email();

      return {
        mutationFn: () => {
          if (!state) {
            throw new Error('Invalid state');
          }
          return this.newAdminService.testEmailProviderSMTP(buildRequest(state, email));
        },
      };
    });
  }
}
