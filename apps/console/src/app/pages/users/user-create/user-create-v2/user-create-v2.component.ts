import { ChangeDetectionStrategy, Component, DestroyRef, OnInit, signal } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { ToastService } from 'src/app/services/toast.service';
import { FormBuilder, FormControl } from '@angular/forms';
import { UserService } from 'src/app/services/user.service';
import { Location } from '@angular/common';
import {
  emailValidator,
  minLengthValidator,
  passwordConfirmValidator,
  requiredValidator,
} from 'src/app/modules/form-field/validators/validators';
import { NewMgmtService } from 'src/app/services/new-mgmt.service';
import {
  defaultIfEmpty,
  defer,
  EMPTY,
  firstValueFrom,
  mergeWith,
  NEVER,
  Observable,
  of,
  shareReplay,
  TimeoutError,
} from 'rxjs';
import { catchError, filter, map, startWith, timeout } from 'rxjs/operators';
import { PasswordComplexityPolicy } from '@zitadel/proto/zitadel/policy_pb';
import { MessageInitShape } from '@bufbuild/protobuf';
import { AddHumanUserRequestSchema } from '@zitadel/proto/zitadel/user/v2/user_service_pb';
import { LoginV2FeatureFlag } from '@zitadel/proto/zitadel/feature/v2/feature_pb';
import { withLatestFromSynchronousFix } from 'src/app/utils/withLatestFromSynchronousFix';
import { PasswordComplexityValidatorFactoryService } from 'src/app/services/password-complexity-validator-factory.service';
import { NewFeatureService } from 'src/app/services/new-feature.service';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';

type PwdForm = ReturnType<UserCreateV2Component['buildPwdForm']>;
type AuthenticationFactor =
  | { factor: 'none' }
  | { factor: 'initialPassword'; form: PwdForm; policy: PasswordComplexityPolicy }
  | { factor: 'invitation' };

@Component({
  selector: 'cnsl-user-create-v2',
  templateUrl: './user-create-v2.component.html',
  styleUrls: ['./user-create-v2.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class UserCreateV2Component implements OnInit {
  protected readonly loading = signal(false);

  protected readonly userForm: ReturnType<typeof this.buildUserForm>;

  private readonly passwordComplexityPolicy$: Observable<PasswordComplexityPolicy>;
  protected readonly authenticationFactor$: Observable<AuthenticationFactor>;
  private readonly useLoginV2$: Observable<LoginV2FeatureFlag | undefined>;

  constructor(
    private readonly router: Router,
    private readonly toast: ToastService,
    private readonly fb: FormBuilder,
    private readonly userService: UserService,
    private readonly newMgmtService: NewMgmtService,
    private readonly passwordComplexityValidatorFactory: PasswordComplexityValidatorFactoryService,
    private readonly featureService: NewFeatureService,
    private readonly destroyRef: DestroyRef,
    private readonly route: ActivatedRoute,
    protected readonly location: Location,
    private readonly authService: GrpcAuthService,
  ) {
    this.userForm = this.buildUserForm();

    this.passwordComplexityPolicy$ = this.getPasswordComplexityPolicy().pipe(shareReplay({ refCount: true, bufferSize: 1 }));
    this.authenticationFactor$ = this.getAuthenticationFactor(this.userForm, this.passwordComplexityPolicy$);
    this.useLoginV2$ = this.getUseLoginV2().pipe(shareReplay({ refCount: true, bufferSize: 1 }));
  }

  ngOnInit(): void {
    this.useLoginV2$.pipe(takeUntilDestroyed(this.destroyRef)).subscribe();
    this.authenticationFactor$.pipe(takeUntilDestroyed(this.destroyRef)).subscribe(async ({ factor }) => {
      // preserve current factor choice when reloading helpful while developing
      await this.router.navigate([], {
        relativeTo: this.route,
        queryParams: {
          factor,
        },
        queryParamsHandling: 'merge',
      });
    });
  }

  public buildUserForm() {
    const param = this.route.snapshot.queryParamMap.get('factor');
    const authenticationFactor =
      param === 'none' ? param : param === 'initialPassword' ? param : param === 'invitation' ? param : 'none';

    return this.fb.group({
      email: new FormControl('', { nonNullable: true, validators: [requiredValidator, emailValidator] }),
      username: new FormControl('', { nonNullable: true, validators: [requiredValidator, minLengthValidator(2)] }),
      givenName: new FormControl('', { nonNullable: true, validators: [requiredValidator] }),
      familyName: new FormControl('', { nonNullable: true, validators: [requiredValidator] }),
      emailVerified: new FormControl(false, { nonNullable: true }),
      authenticationFactor: new FormControl<AuthenticationFactor['factor']>(authenticationFactor, {
        nonNullable: true,
      }),
    });
  }

  private getPasswordComplexityPolicy() {
    return defer(() => this.newMgmtService.getPasswordComplexityPolicy()).pipe(
      map(({ policy }) => policy),
      filter(Boolean),
      catchError((error) => {
        this.toast.showError(error);
        return EMPTY;
      }),
    );
  }

  private getAuthenticationFactor(
    userForm: typeof this.userForm,
    passwordComplexityPolicy$: Observable<PasswordComplexityPolicy>,
  ): Observable<AuthenticationFactor> {
    const pwdForm$ = passwordComplexityPolicy$.pipe(
      defaultIfEmpty(undefined),
      map((policy) => this.buildPwdForm(policy)),
    );

    return userForm.controls.authenticationFactor.valueChanges.pipe(
      startWith(userForm.controls.authenticationFactor.value),
      withLatestFromSynchronousFix(pwdForm$, passwordComplexityPolicy$),
      map(([factor, form, policy]) => {
        if (factor === 'initialPassword') {
          return { factor, form, policy };
        }
        // reset emailVerified when we switch to invitation
        if (factor === 'invitation') {
          userForm.controls.emailVerified.setValue(false);
        }
        return { factor };
      }),
    );
  }

  private buildPwdForm(policy: PasswordComplexityPolicy | undefined) {
    return this.fb.group({
      password: new FormControl('', {
        nonNullable: true,
        validators: this.passwordComplexityValidatorFactory.buildValidators(policy),
      }),
      confirmPassword: new FormControl('', {
        nonNullable: true,
        validators: [requiredValidator, passwordConfirmValidator()],
      }),
    });
  }

  private getUseLoginV2() {
    return defer(() => this.featureService.getInstanceFeatures()).pipe(
      map(({ loginV2 }) => loginV2),
      timeout(1000),
      catchError((err) => {
        if (!(err instanceof TimeoutError)) {
          this.toast.showError(err);
        }
        return of(undefined);
      }),
      mergeWith(NEVER),
    );
  }

  protected async createUserV2(authenticationFactor: AuthenticationFactor) {
    try {
      await this.createUserV2Try(authenticationFactor);
    } catch (error) {
      this.toast.showError(error);
    } finally {
      this.loading.set(false);
    }
  }

  private async createUserV2Try(authenticationFactor: AuthenticationFactor) {
    this.loading.set(true);

    const org = await this.authService.getActiveOrg();

    const userValues = this.userForm.getRawValue();

    const humanReq: MessageInitShape<typeof AddHumanUserRequestSchema> = {
      organization: { org: { case: 'orgId', value: org.id } },
      username: userValues.username,
      profile: {
        givenName: userValues.givenName,
        familyName: userValues.familyName,
      },
      email: {
        email: userValues.email,
        verification: {
          case: 'isVerified',
          value: userValues.emailVerified,
        },
      },
    };

    if (authenticationFactor.factor === 'initialPassword') {
      const { password } = authenticationFactor.form.getRawValue();
      humanReq['passwordType'] = {
        case: 'password',
        value: {
          password,
        },
      };
    }

    const resp = await this.userService.addHumanUser(humanReq);
    if (authenticationFactor.factor === 'invitation') {
      const url = await this.getUrlTemplate();
      await this.userService.createInviteCode({
        userId: resp.userId,
        verification: {
          case: 'sendCode',
          value: url
            ? {
                urlTemplate: `${url}verify?code={{.Code}}&userId={{.UserID}}&organization={{.OrgID}}&invite=true`,
              }
            : {},
        },
      });
    }

    this.toast.showInfo('USER.TOAST.CREATED', true);
    await this.router.navigate(['users', resp.userId], { queryParams: { new: true } });
  }

  private async getUrlTemplate() {
    const useLoginV2 = await firstValueFrom(this.useLoginV2$);
    if (!useLoginV2?.required) {
      // loginV2 is not enabled
      return undefined;
    }

    const { baseUri } = useLoginV2;
    // if base uri is not set, we use the default for the cloud hosted login v2
    if (!baseUri) {
      return new URL(location.origin + '/ui/v2/login/');
    }

    const baseUriWithTrailingSlash = baseUri.endsWith('/') ? baseUri : `${baseUri}/`;
    try {
      // first we try to create a URL directly from the baseUri
      return new URL(baseUriWithTrailingSlash);
    } catch (_) {
      // if this does not work we assume that the baseUri is relative,
      // and we need to add the location.origin
      // we make sure the relative url has a slash at the beginning and end
      const baseUriWithSlashes = baseUriWithTrailingSlash.startsWith('/')
        ? baseUriWithTrailingSlash
        : `/${baseUriWithTrailingSlash}`;
      return new URL(location.origin + baseUriWithSlashes);
    }
  }
}
