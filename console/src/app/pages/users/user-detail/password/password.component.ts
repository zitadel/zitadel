import { Component, DestroyRef, OnInit } from '@angular/core';
import { AbstractControl, UntypedFormBuilder, UntypedFormGroup, Validators } from '@angular/forms';
import { ActivatedRoute } from '@angular/router';
import {
  take,
  map,
  switchMap,
  firstValueFrom,
  mergeWith,
  Observable,
  defer,
  of,
  shareReplay,
  combineLatestWith,
} from 'rxjs';
import {
  containsLowerCaseValidator,
  containsNumberValidator,
  containsSymbolValidator,
  containsUpperCaseValidator,
  minLengthValidator,
  passwordConfirmValidator,
  requiredValidator,
} from 'src/app/modules/form-field/validators/validators';
import { PasswordComplexityPolicy } from 'src/app/proto/generated/zitadel/policy_pb';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ToastService } from 'src/app/services/toast.service';
import { catchError, filter } from 'rxjs/operators';
import { User } from 'src/app/proto/generated/zitadel/user_pb';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { UserService } from '../../../../services/user.service';

@Component({
  selector: 'cnsl-password',
  templateUrl: './password.component.html',
  styleUrls: ['./password.component.scss'],
})
export class PasswordComponent implements OnInit {
  private readonly breadcrumb$: Observable<Breadcrumb[]>;
  protected readonly username$: Observable<string>;
  protected readonly id$: Observable<string | undefined>;
  protected readonly form$: Observable<UntypedFormGroup>;
  protected readonly passwordPolicy$: Observable<PasswordComplexityPolicy.AsObject | undefined>;
  protected readonly user$: Observable<User.AsObject>;

  constructor(
    activatedRoute: ActivatedRoute,
    private readonly fb: UntypedFormBuilder,
    private readonly authService: GrpcAuthService,
    private readonly userService: UserService,
    private readonly toast: ToastService,
    private readonly breadcrumbService: BreadcrumbService,
    private readonly destroyRef: DestroyRef,
  ) {
    const usernameParam$ = activatedRoute.queryParamMap.pipe(
      map((params) => params.get('username')),
      filter(Boolean),
    );
    this.id$ = activatedRoute.paramMap.pipe(map((params) => params.get('id') ?? undefined));

    this.user$ = this.authService.user.pipe(take(1), filter(Boolean));
    this.username$ = usernameParam$.pipe(mergeWith(this.user$.pipe(map((user) => user.preferredLoginName))));

    this.breadcrumb$ = this.getBreadcrumb$(this.id$, this.user$);
    this.passwordPolicy$ = this.getPasswordPolicy$().pipe(shareReplay({ refCount: true, bufferSize: 1 }));
    const validators$ = this.getValidators$(this.passwordPolicy$);
    this.form$ = this.getForm$(this.id$, validators$);
  }

  private getBreadcrumb$(id$: Observable<string | undefined>, user$: Observable<User.AsObject>): Observable<Breadcrumb[]> {
    return id$.pipe(
      switchMap(async (id) => {
        if (id) {
          return [
            new Breadcrumb({
              type: BreadcrumbType.ORG,
              routerLink: ['/org'],
            }),
          ];
        }
        const user = await firstValueFrom(user$);
        if (!user) {
          return [];
        }
        return [
          new Breadcrumb({
            type: BreadcrumbType.AUTHUSER,
            name: user.human?.profile?.displayName,
            routerLink: ['/users', 'me'],
          }),
        ];
      }),
    );
  }

  private getValidators$(
    passwordPolicy$: Observable<PasswordComplexityPolicy.AsObject | undefined>,
  ): Observable<Validators[]> {
    return passwordPolicy$.pipe(
      map((policy) => {
        const validators: Validators[] = [requiredValidator];
        if (!policy) {
          return validators;
        }
        if (policy.minLength) {
          validators.push(minLengthValidator(policy.minLength));
        }
        if (policy.hasLowercase) {
          validators.push(containsLowerCaseValidator);
        }
        if (policy.hasUppercase) {
          validators.push(containsUpperCaseValidator);
        }
        if (policy.hasNumber) {
          validators.push(containsNumberValidator);
        }
        if (policy.hasSymbol) {
          validators.push(containsSymbolValidator);
        }
        return validators;
      }),
    );
  }

  private getForm$(
    id$: Observable<string | undefined>,
    validators$: Observable<Validators[]>,
  ): Observable<UntypedFormGroup> {
    return id$.pipe(
      combineLatestWith(validators$),
      map(([id, validators]) => {
        if (id) {
          return this.fb.group({
            password: ['', validators],
            confirmPassword: ['', [requiredValidator, passwordConfirmValidator()]],
          });
        } else {
          return this.fb.group({
            currentPassword: ['', requiredValidator],
            newPassword: ['', validators],
            confirmPassword: ['', [requiredValidator, passwordConfirmValidator()]],
          });
        }
      }),
    );
  }

  private getPasswordPolicy$(): Observable<PasswordComplexityPolicy.AsObject | undefined> {
    return defer(() => this.authService.getMyPasswordComplexityPolicy()).pipe(
      map((resp) => resp.policy),
      catchError(() => of(undefined)),
    );
  }

  ngOnInit() {
    this.breadcrumb$.pipe(takeUntilDestroyed(this.destroyRef)).subscribe((breadcrumbs) => {
      this.breadcrumbService.setBreadcrumb(breadcrumbs);
    });
  }

  public async setInitialPassword(userId: string, form: UntypedFormGroup): Promise<void> {
    const password = this.password(form)?.value;

    if (form.invalid || !password) {
      return;
    }

    try {
      await this.userService.setPassword({
        userId,
        newPassword: {
          password,
          changeRequired: false,
        },
      });
    } catch (error) {
      this.toast.showError(error);
      return;
    }
    this.toast.showInfo('USER.TOAST.INITIALPASSWORDSET', true);
    window.history.back();
  }

  public async setPassword(form: UntypedFormGroup, user: User.AsObject): Promise<void> {
    const currentPassword = this.currentPassword(form);
    const newPassword = this.newPassword(form);

    if (form.invalid || !currentPassword?.value || !newPassword?.value || newPassword?.invalid) {
      return;
    }

    try {
      await this.userService.setPassword({
        userId: user.id,
        newPassword: {
          password: newPassword.value,
          changeRequired: false,
        },
        verification: {
          case: 'currentPassword',
          value: currentPassword.value,
        },
      });
    } catch (error) {
      this.toast.showError(error);
      return;
    }
    this.toast.showInfo('USER.TOAST.PASSWORDCHANGED', true);
    window.history.back();
  }

  public password(form: UntypedFormGroup): AbstractControl | null {
    return form.get('password');
  }

  public newPassword(form: UntypedFormGroup): AbstractControl | null {
    return form.get('newPassword');
  }

  public currentPassword(form: UntypedFormGroup): AbstractControl | null {
    return form.get('currentPassword');
  }
}
