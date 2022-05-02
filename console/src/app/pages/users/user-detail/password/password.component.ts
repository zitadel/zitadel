import { Component, OnDestroy } from '@angular/core';
import { AbstractControl, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute } from '@angular/router';
import { Subscription, take } from 'rxjs';
import { lowerCaseValidator, numberValidator, symbolValidator, upperCaseValidator } from 'src/app/pages/validators';
import { PasswordComplexityPolicy } from 'src/app/proto/generated/zitadel/policy_pb';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

function passwordConfirmValidator(c: AbstractControl): any {
  if (!c.parent || !c) {
    return;
  }
  const pwd = c.parent.get('password');
  const cpwd = c.parent.get('confirmPassword');

  if (!pwd || !cpwd) {
    return;
  }
  if (pwd.value !== cpwd.value) {
    return { invalid: true, notequal: 'Password is not equal' };
  }
}

@Component({
  selector: 'cnsl-password',
  templateUrl: './password.component.html',
  styleUrls: ['./password.component.scss'],
})
export class PasswordComponent implements OnDestroy {
  userId: string = '';

  public policy!: PasswordComplexityPolicy.AsObject;
  public passwordForm!: FormGroup;

  private formSub: Subscription = new Subscription();

  constructor(
    activatedRoute: ActivatedRoute,
    private fb: FormBuilder,
    private authService: GrpcAuthService,
    private mgmtUserService: ManagementService,
    private toast: ToastService,
    private breadcrumbService: BreadcrumbService,
  ) {
    activatedRoute.params.subscribe((data) => {
      const { id } = data;
      if (id) {
        this.userId = id;

        breadcrumbService.setBreadcrumb([
          new Breadcrumb({
            type: BreadcrumbType.IAM,
            name: 'Instance',
            routerLink: ['/instance'],
          }),
          new Breadcrumb({
            type: BreadcrumbType.ORG,
            routerLink: ['/org'],
          }),
        ]);
      } else {
        this.authService.user.pipe(take(1)).subscribe((user) => {
          if (user) {
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

      const validators: Validators[] = [Validators.required];
      this.authService
        .getMyPasswordComplexityPolicy()
        .then((resp) => {
          if (resp.policy) {
            this.policy = resp.policy;
          }
          if (this.policy.minLength) {
            validators.push(Validators.minLength(this.policy.minLength));
          }
          if (this.policy.hasLowercase) {
            validators.push(lowerCaseValidator);
          }
          if (this.policy.hasUppercase) {
            validators.push(upperCaseValidator);
          }
          if (this.policy.hasNumber) {
            validators.push(numberValidator);
          }
          if (this.policy.hasSymbol) {
            validators.push(symbolValidator);
          }

          this.setupForm(validators);
        })
        .catch((error) => {
          this.setupForm(validators);
        });
    });
  }

  ngOnDestroy(): void {
    this.formSub.unsubscribe();
  }

  setupForm(validators: Validators[]): void {
    if (this.userId) {
      this.passwordForm = this.fb.group({
        password: ['', validators],
        confirmPassword: ['', [...validators, passwordConfirmValidator]],
      });
    } else {
      this.passwordForm = this.fb.group({
        currentPassword: ['', Validators.required],
        newPassword: ['', validators],
        confirmPassword: ['', [...validators, passwordConfirmValidator]],
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
