import { animate, style, transition, trigger } from '@angular/animations';
import { Location } from '@angular/common';
import { Component } from '@angular/core';
import { AbstractControl, UntypedFormBuilder, UntypedFormGroup, ValidatorFn, Validators } from '@angular/forms';
import { MatLegacySlideToggleChange as MatSlideToggleChange } from '@angular/material/legacy-slide-toggle';
import { Router } from '@angular/router';
import {
  containsLowerCaseValidator,
  containsNumberValidator,
  containsSymbolValidator,
  containsUpperCaseValidator,
  minLengthValidator,
  passwordConfirmValidator,
  requiredValidator,
} from 'src/app/modules/form-field/validators/validators';
import { SetUpOrgRequest } from 'src/app/proto/generated/zitadel/admin_pb';
import { PasswordComplexityPolicy } from 'src/app/proto/generated/zitadel/policy_pb';
import { Gender } from 'src/app/proto/generated/zitadel/user_pb';
import { AdminService } from 'src/app/services/admin.service';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
  selector: 'cnsl-org-create',
  templateUrl: './org-create.component.html',
  styleUrls: ['./org-create.component.scss'],
  animations: [
    trigger('openClose', [
      transition(':enter', [
        style({ height: '0', opacity: 0 }),
        animate('150ms ease-in-out', style({ height: '*', opacity: 1 })),
      ]),
      transition(':leave', [animate('150ms ease-in-out', style({ height: '0', opacity: 0 }))]),
    ]),
  ],
})
export class OrgCreateComponent {
  public orgForm: UntypedFormGroup = this.fb.group({
    name: ['', [requiredValidator]],
    domain: [''],
  });

  public userForm?: UntypedFormGroup;
  public pwdForm?: UntypedFormGroup;

  public genders: Gender[] = [Gender.GENDER_FEMALE, Gender.GENDER_MALE, Gender.GENDER_UNSPECIFIED];
  public languages: string[] = ['de', 'en', 'it', 'fr', 'pl', 'zh'];

  public policy?: PasswordComplexityPolicy.AsObject;
  public usePassword: boolean = false;

  public forSelf: boolean = true;

  constructor(
    private router: Router,
    private toast: ToastService,
    private adminService: AdminService,
    private _location: Location,
    private fb: UntypedFormBuilder,
    private mgmtService: ManagementService,
    breadcrumbService: BreadcrumbService,
  ) {
    const instanceBread = new Breadcrumb({
      type: BreadcrumbType.INSTANCE,
      name: 'Instance',
      routerLink: ['/instance'],
    });

    breadcrumbService.setBreadcrumb([instanceBread]);
    this.initForm();

    this.adminService.getSupportedLanguages().then((supportedResp) => {
      this.languages = supportedResp.languagesList;
    });
  }

  public createSteps: number = 2;
  public currentCreateStep: number = 1;

  public finish(): void {
    const createOrgRequest: SetUpOrgRequest.Org = new SetUpOrgRequest.Org();
    createOrgRequest.setName(this.name?.value);
    createOrgRequest.setDomain(this.domain?.value);

    const humanRequest: SetUpOrgRequest.Human = new SetUpOrgRequest.Human();
    humanRequest.setEmail(
      new SetUpOrgRequest.Human.Email().setEmail(this.email?.value).setIsEmailVerified(this.isVerified?.value),
    );
    humanRequest.setUserName(this.userName?.value);

    const profile: SetUpOrgRequest.Human.Profile = new SetUpOrgRequest.Human.Profile();
    profile.setFirstName(this.firstName?.value);
    profile.setLastName(this.lastName?.value);
    profile.setNickName(this.nickName?.value);
    profile.setGender(this.gender?.value);
    profile.setPreferredLanguage(this.preferredLanguage?.value);

    humanRequest.setProfile(profile);
    if (this.usePassword && this.password) {
      humanRequest.setPassword(this.password?.value);
    }

    this.adminService
      .SetUpOrg(createOrgRequest, humanRequest)
      .then(() => {
        this.router.navigate(['/orgs']);
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public next(): void {
    this.currentCreateStep++;
  }

  public previous(): void {
    this.currentCreateStep--;
  }

  private initForm(): void {
    this.userForm = this.fb.group({
      userName: ['', [requiredValidator]],
      firstName: ['', [requiredValidator]],
      lastName: ['', [requiredValidator]],
      email: ['', [requiredValidator]],
      isVerified: [false, []],
      gender: [''],
      nickName: [''],
      preferredLanguage: [''],
    });
  }

  public initPwdValidators(): void {
    const validators: Validators[] = [requiredValidator];

    if (this.usePassword) {
      this.mgmtService.getDefaultPasswordComplexityPolicy().then((data) => {
        if (data.policy) {
          this.policy = data.policy;

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

          const pwdValidators = [...validators] as ValidatorFn[];
          const confirmPwdValidators = [requiredValidator, passwordConfirmValidator()] as ValidatorFn[];
          this.pwdForm = this.fb.group({
            password: ['', pwdValidators],
            confirmPassword: ['', confirmPwdValidators],
          });
        }
      });
    } else {
      this.pwdForm = this.fb.group({
        password: ['', []],
        confirmPassword: ['', []],
      });
    }
  }

  public changeSelf(change: MatSlideToggleChange): void {
    if (change.checked) {
      this.createSteps = 1;

      this.orgForm = this.fb.group({
        name: ['', [requiredValidator]],
      });
    } else {
      this.createSteps = 2;

      this.orgForm = this.fb.group({
        name: ['', [requiredValidator]],
        domain: [''],
      });
    }
  }

  public createOrgForSelf(): void {
    if (this.name && this.name.value) {
      this.mgmtService
        .addOrg(this.name.value)
        .then(() => {
          this.router.navigate(['/orgs']);
        })
        .catch((error) => {
          this.toast.showError(error);
        });
    }
  }

  public get name(): AbstractControl | null {
    return this.orgForm.get('name');
  }

  public get domain(): AbstractControl | null {
    return this.orgForm.get('domain');
  }

  public get userName(): AbstractControl | null {
    return this.userForm?.get('userName') ?? null;
  }

  public get firstName(): AbstractControl | null {
    return this.userForm?.get('firstName') ?? null;
  }

  public get lastName(): AbstractControl | null {
    return this.userForm?.get('lastName') ?? null;
  }

  public get email(): AbstractControl | null {
    return this.userForm?.get('email') ?? null;
  }

  public get isVerified(): AbstractControl | null {
    return this.userForm?.get('isVerified') ?? null;
  }

  public get nickName(): AbstractControl | null {
    return this.userForm?.get('nickName') ?? null;
  }

  public get preferredLanguage(): AbstractControl | null {
    return this.userForm?.get('preferredLanguage') ?? null;
  }

  public get gender(): AbstractControl | null {
    return this.userForm?.get('gender') ?? null;
  }

  public get password(): AbstractControl | null {
    return this.pwdForm?.get('password') ?? null;
  }

  public get confirmPassword(): AbstractControl | null {
    return this.pwdForm?.get('confirmPassword') ?? null;
  }

  public close(): void {
    this._location.back();
  }
}
