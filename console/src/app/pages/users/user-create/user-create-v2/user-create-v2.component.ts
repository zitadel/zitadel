import { ChangeDetectionStrategy, Component, signal } from '@angular/core';
import { Router } from '@angular/router';
import { ToastService } from 'src/app/services/toast.service';
import { FormBuilder, FormControl } from '@angular/forms';
import { UserService } from 'src/app/services/user.service';
import { LanguagesService } from 'src/app/services/languages.service';
import { Location } from '@angular/common';
import {
  emailValidator,
  minLengthValidator,
  passwordConfirmValidator,
  requiredValidator,
} from 'src/app/modules/form-field/validators/validators';
import { NewMgmtService } from 'src/app/services/new-mgmt.service';
import { defaultIfEmpty, defer, EMPTY, Observable, shareReplay } from 'rxjs';
import { catchError, filter, map, startWith, tap } from 'rxjs/operators';
import { PasswordComplexityPolicy } from '@zitadel/proto/zitadel/policy_pb';
import { MessageInitShape } from '@bufbuild/protobuf';
import { AddHumanUserRequestSchema } from '@zitadel/proto/zitadel/user/v2/user_service_pb';
import { withLatestFromSynchronousFix } from 'src/app/utils/withLatestFromSynchronousFix';
import { PasswordComplexityValidatorFactoryService } from 'src/app/services/password-complexity-validator-factory.service';

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
export class UserCreateV2Component {
  protected readonly loading = signal(false);

  protected readonly userForm: ReturnType<typeof this.buildUserForm>;

  private readonly passwordComplexityPolicy$: Observable<PasswordComplexityPolicy>;
  protected readonly authenticationFactor$: Observable<AuthenticationFactor>;

  constructor(
    private readonly router: Router,
    private readonly toast: ToastService,
    private readonly fb: FormBuilder,
    private readonly userService: UserService,
    private readonly newMgmtService: NewMgmtService,
    private readonly passwordComplexityValidatorFactory: PasswordComplexityValidatorFactoryService,
    protected readonly location: Location,
    public readonly langSvc: LanguagesService,
  ) {
    this.userForm = this.buildUserForm();

    this.passwordComplexityPolicy$ = this.getPasswordComplexityPolicy().pipe(shareReplay({ refCount: true, bufferSize: 1 }));
    this.authenticationFactor$ = this.getAuthenticationFactor(this.userForm, this.passwordComplexityPolicy$);
  }

  public buildUserForm() {
    return this.fb.group({
      email: new FormControl('', { nonNullable: true, validators: [requiredValidator, emailValidator] }),
      username: new FormControl('', { nonNullable: true, validators: [requiredValidator, minLengthValidator(2)] }),
      givenName: new FormControl('', { nonNullable: true, validators: [requiredValidator] }),
      familyName: new FormControl('', { nonNullable: true, validators: [requiredValidator] }),
      nickName: new FormControl('', { nonNullable: true }),
      emailVerified: new FormControl(false, { nonNullable: true }),
      authenticationFactor: new FormControl<AuthenticationFactor['factor']>('none', { nonNullable: true }),
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
      tap((factor) => console.log('factor', factor)),
      withLatestFromSynchronousFix(pwdForm$, passwordComplexityPolicy$),
      map(([factor, form, policy]) => {
        if (factor === 'initialPassword') {
          return { factor, form, policy };
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

    const userValues = this.userForm.getRawValue();

    const humanReq: MessageInitShape<typeof AddHumanUserRequestSchema> = {
      username: userValues.username,
      profile: {
        givenName: userValues.givenName,
        familyName: userValues.familyName,
        nickName: userValues.nickName,
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
      humanReq.passwordType = {
        case: 'password',
        value: {
          password,
        },
      };
    }

    const resp = await this.userService.addHumanUser(humanReq);
    if (authenticationFactor.factor === 'invitation') {
      await this.userService.createInviteCode({
        userId: resp.userId,
        verification: {
          case: 'sendCode',
          value: {},
        },
      });
    }

    this.toast.showInfo('USER.TOAST.CREATED', true);
    await this.router.navigate(['users', resp.userId], { queryParams: { new: true } });
  }
}
