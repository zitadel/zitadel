import { Component, DestroyRef, OnInit } from '@angular/core';
import { AbstractControl, UntypedFormBuilder, UntypedFormGroup } from '@angular/forms';
import { ActivatedRoute } from '@angular/router';
import {
  map,
  switchMap,
  firstValueFrom,
  mergeWith,
  Observable,
  defer,
  of,
  shareReplay,
  combineLatestWith,
  EMPTY,
} from 'rxjs';
import { passwordConfirmValidator, requiredValidator } from 'src/app/modules/form-field/validators/validators';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ToastService } from 'src/app/services/toast.service';
import { catchError, filter } from 'rxjs/operators';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { UserService } from 'src/app/services/user.service';
import { User } from '@zitadel/proto/zitadel/user/v2/user_pb';
import { NewAuthService } from 'src/app/services/new-auth.service';
import { PasswordComplexityPolicy } from '@zitadel/proto/zitadel/policy_pb';
import { PasswordComplexityValidatorFactoryService } from 'src/app/services/password-complexity-validator-factory.service';

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
  protected readonly passwordPolicy$: Observable<PasswordComplexityPolicy | undefined>;
  protected readonly user$: Observable<User>;

  constructor(
    private readonly activatedRoute: ActivatedRoute,
    private readonly fb: UntypedFormBuilder,
    private readonly userService: UserService,
    private readonly newAuthService: NewAuthService,
    private readonly toast: ToastService,
    private readonly breadcrumbService: BreadcrumbService,
    private readonly destroyRef: DestroyRef,
    private readonly passwordComplexityValidatorFactory: PasswordComplexityValidatorFactoryService,
  ) {
    this.id$ = activatedRoute.paramMap.pipe(map((params) => params.get('id') ?? undefined));
    this.user$ = this.getUser().pipe(shareReplay({ refCount: true, bufferSize: 1 }));
    this.username$ = this.getUsername(this.user$);
    this.breadcrumb$ = this.getBreadcrumb$(this.id$, this.user$);
    this.passwordPolicy$ = this.getPasswordPolicy$().pipe(shareReplay({ refCount: true, bufferSize: 1 }));
    this.form$ = this.getForm$(this.id$, this.passwordPolicy$);
  }

  ngOnInit() {
    this.breadcrumb$.pipe(takeUntilDestroyed(this.destroyRef)).subscribe((breadcrumbs) => {
      this.breadcrumbService.setBreadcrumb(breadcrumbs);
    });
  }

  private getUser() {
    return this.userService.user$.pipe(
      catchError((err) => {
        this.toast.showError(err);
        return EMPTY;
      }),
    );
  }

  private getUsername(user$: Observable<User>) {
    const prefferedLoginName$ = user$.pipe(map((user) => user.preferredLoginName));

    return this.activatedRoute.queryParamMap.pipe(
      map((params) => params.get('username')),
      filter(Boolean),
      mergeWith(prefferedLoginName$),
    );
  }

  private getBreadcrumb$(id$: Observable<string | undefined>, user$: Observable<User>): Observable<Breadcrumb[]> {
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
            name: (user.type.case === 'human' && user.type.value.profile?.displayName) || undefined,
            routerLink: ['/users', 'me'],
          }),
        ];
      }),
    );
  }

  private getPasswordPolicy$(): Observable<PasswordComplexityPolicy | undefined> {
    return defer(() => this.newAuthService.getMyPasswordComplexityPolicy()).pipe(
      map((resp) => resp.policy),
      catchError((err) => {
        this.toast.showError(err);
        return of(undefined);
      }),
    );
  }

  private getForm$(
    id$: Observable<string | undefined>,
    policy$: Observable<PasswordComplexityPolicy | undefined>,
  ): Observable<UntypedFormGroup> {
    const validators$ = policy$.pipe(map((policy) => this.passwordComplexityValidatorFactory.buildValidators(policy)));

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

  public async setPassword(form: UntypedFormGroup, user: User): Promise<void> {
    const currentPassword = this.currentPassword(form);
    const newPassword = this.newPassword(form);

    if (form.invalid || !currentPassword?.value || !newPassword?.value || newPassword?.invalid) {
      return;
    }

    try {
      await this.userService.setPassword({
        userId: user.userId,
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
