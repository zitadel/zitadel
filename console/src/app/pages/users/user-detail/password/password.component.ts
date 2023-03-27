import { Component, OnDestroy } from '@angular/core';
import { AbstractControl, UntypedFormBuilder, UntypedFormGroup, Validators } from '@angular/forms';
import { ActivatedRoute } from '@angular/router';
import { Subject, Subscription, take, takeUntil } from 'rxjs';
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
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
  selector: 'cnsl-password',
  templateUrl: './password.component.html',
  styleUrls: ['./password.component.scss'],
})
export class PasswordComponent implements OnDestroy {
  userId: string = '';
  public username: string = '';

  public policy!: PasswordComplexityPolicy.AsObject;
  public passwordForm!: UntypedFormGroup;

  private formSub: Subscription = new Subscription();
  private destroy$: Subject<void> = new Subject();

  constructor(
    activatedRoute: ActivatedRoute,
    private fb: UntypedFormBuilder,
    private authService: GrpcAuthService,
    private mgmtUserService: ManagementService,
    private toast: ToastService,
    private breadcrumbService: BreadcrumbService,
  ) {
    activatedRoute.queryParams.pipe(takeUntil(this.destroy$)).subscribe((data) => {
      const { username } = data;
      this.username = username;
    });
    activatedRoute.params.pipe(takeUntil(this.destroy$)).subscribe((data) => {
      const { id } = data;
      if (id) {
        this.userId = id;
        breadcrumbService.setBreadcrumb([
          new Breadcrumb({
            type: BreadcrumbType.ORG,
            routerLink: ['/org'],
          }),
        ]);
      } else {
        this.authService.user.pipe(take(1)).subscribe((user) => {
          if (user) {
            this.username = user.preferredLoginName;
            this.breadcrumbService.setBreadcrumb([
              new Breadcrumb({
                type: BreadcrumbType.AUTHUSER,
                name: user.human?.profile?.displayName,
                routerLink: ['/users', 'me'],
              }),
            ]);
          }
        });
      }

      const validators: Validators[] = [requiredValidator];
      this.authService
        .getMyPasswordComplexityPolicy()
        .then((resp) => {
          if (resp.policy) {
            this.policy = resp.policy;
          }
          if (this.policy.minLength) {
            validators.push(minLengthValidator(this.policy.minLength));
          }
          if (this.policy.hasLowercase) {
            validators.push(containsLowerCaseValidator);
          }
          if (this.policy.hasUppercase) {
            validators.push(containsUpperCaseValidator);
          }
          if (this.policy.hasNumber) {
            validators.push(containsNumberValidator);
          }
          if (this.policy.hasSymbol) {
            validators.push(containsSymbolValidator);
          }

          this.setupForm(validators);
        })
        .catch((error) => {
          this.setupForm(validators);
        });
    });
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
    this.formSub.unsubscribe();
  }

  setupForm(validators: Validators[]): void {
    if (this.userId) {
      this.passwordForm = this.fb.group({
        password: ['', validators],
        confirmPassword: ['', [requiredValidator, passwordConfirmValidator()]],
      });
    } else {
      this.passwordForm = this.fb.group({
        currentPassword: ['', requiredValidator],
        newPassword: ['', validators],
        confirmPassword: ['', [requiredValidator, passwordConfirmValidator()]],
      });
    }
  }

  public setInitialPassword(userId: string): void {
    if (this.passwordForm.valid && this.password && this.password.value) {
      this.mgmtUserService
        .setHumanInitialPassword(userId, this.password.value)
        .then((data: any) => {
          this.toast.showInfo('USER.TOAST.INITIALPASSWORDSET', true);
          window.history.back();
        })
        .catch((error) => {
          this.toast.showError(error);
        });
    }
  }

  public setPassword(): void {
    if (
      this.passwordForm.valid &&
      this.currentPassword &&
      this.currentPassword.value &&
      this.newPassword &&
      this.newPassword.value &&
      this.newPassword.valid
    ) {
      this.authService
        .updateMyPassword(this.currentPassword.value, this.newPassword.value)
        .then((data: any) => {
          this.toast.showInfo('USER.TOAST.PASSWORDCHANGED', true);
          window.history.back();
        })
        .catch((error) => {
          this.toast.showError(error);
        });
    }
  }

  public get password(): AbstractControl | null {
    return this.passwordForm.get('password');
  }

  public get newPassword(): AbstractControl | null {
    return this.passwordForm.get('newPassword');
  }

  public get currentPassword(): AbstractControl | null {
    return this.passwordForm.get('currentPassword');
  }

  public get confirmPassword(): AbstractControl | null {
    return this.passwordForm.get('confirmPassword');
  }
}
