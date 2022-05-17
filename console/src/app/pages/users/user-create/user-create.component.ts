import { ChangeDetectorRef, Component, OnDestroy, ViewChild } from '@angular/core';
import { AbstractControl, FormBuilder, FormGroup, ValidatorFn, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import parsePhoneNumber from 'libphonenumber-js';
import { Subject } from 'rxjs';
import { debounceTime, takeUntil } from 'rxjs/operators';
import { AddHumanUserRequest } from 'src/app/proto/generated/zitadel/management_pb';
import { Domain } from 'src/app/proto/generated/zitadel/org_pb';
import { PasswordComplexityPolicy } from 'src/app/proto/generated/zitadel/policy_pb';
import { Gender } from 'src/app/proto/generated/zitadel/user_pb';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { lowerCaseValidator, numberValidator, symbolValidator, upperCaseValidator } from '../../validators';

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
    return {
      invalid: true,
      notequal: {
        valid: false,
      },
    };
  }
}

@Component({
  selector: 'cnsl-user-create',
  templateUrl: './user-create.component.html',
  styleUrls: ['./user-create.component.scss'],
})
export class UserCreateComponent implements OnDestroy {
  public user: AddHumanUserRequest.AsObject = new AddHumanUserRequest().toObject();
  public genders: Gender[] = [Gender.GENDER_FEMALE, Gender.GENDER_MALE, Gender.GENDER_UNSPECIFIED];
  public languages: string[] = ['de', 'en'];
  public userForm!: FormGroup;
  public pwdForm!: FormGroup;

  public envSuffixLabel: string = '';
  private destroyed$: Subject<void> = new Subject();

  public userLoginMustBeDomain: boolean = false;
  public loading: boolean = false;

  @ViewChild('suffix') public suffix!: any;
  private primaryDomain!: Domain.AsObject;
  public usePassword: boolean = false;
  public policy!: PasswordComplexityPolicy.AsObject;

  constructor(
    private router: Router,
    private toast: ToastService,
    private fb: FormBuilder,
    private mgmtService: ManagementService,
    private changeDetRef: ChangeDetectorRef,
    breadcrumbService: BreadcrumbService,
  ) {
    breadcrumbService.setBreadcrumb([
      new Breadcrumb({
        type: BreadcrumbType.ORG,
        routerLink: ['/org'],
      }),
    ]);

    this.loading = true;
    this.loadOrg();
    this.mgmtService
      .getDomainPolicy()
      .then((resp) => {
        if (resp.policy?.userLoginMustBeDomain) {
          this.userLoginMustBeDomain = resp.policy.userLoginMustBeDomain;
        }
        this.initForm();
        this.loading = false;
        this.envSuffixLabel = this.envSuffix();
        this.changeDetRef.detectChanges();
      })
      .catch((error) => {
        console.error(error);
        this.initForm();
        this.loading = false;
        this.envSuffixLabel = this.envSuffix();
        this.changeDetRef.detectChanges();
      });
  }

  private async loadOrg(): Promise<void> {
    const domains = await this.mgmtService.listOrgDomains();
    const found = domains.resultList.find((resp) => resp.isPrimary);
    if (found) {
      this.primaryDomain = found;
    }
  }

  private initForm(): void {
    this.userForm = this.fb.group({
      email: ['', [Validators.required, Validators.email]],
      userName: ['', [Validators.required, Validators.minLength(2)]],
      firstName: ['', Validators.required],
      lastName: ['', Validators.required],
      nickName: [''],
      gender: [],
      preferredLanguage: [''],
      phone: [''],
      isVerified: [false, []],
    });

    const validators: Validators[] = [Validators.required];

    this.mgmtService.getPasswordComplexityPolicy().then((data) => {
      if (data.policy) {
        this.policy = data.policy;

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
        const pwdValidators = [...validators] as ValidatorFn[];
        const confirmPwdValidators = [...validators, passwordConfirmValidator] as ValidatorFn[];

        this.pwdForm = this.fb.group({
          password: ['', pwdValidators],
          confirmPassword: ['', confirmPwdValidators],
        });
      }
    });

    this.userForm.controls['phone'].valueChanges.pipe(takeUntil(this.destroyed$), debounceTime(300)).subscribe((value) => {
      const phoneNumber = parsePhoneNumber(value ?? '', 'CH');
      if (phoneNumber) {
        const formmatted = phoneNumber.formatInternational();
        const country = phoneNumber.country;
        if (this.phone && country && this.phone.value && this.phone.value !== formmatted) {
          this.phone.setValue(formmatted);
        }
      }
    });
  }

  public createUser(): void {
    this.user = this.userForm.value;

    this.loading = true;

    const profileReq = new AddHumanUserRequest.Profile();
    profileReq.setFirstName(this.firstName?.value);
    profileReq.setLastName(this.lastName?.value);
    profileReq.setNickName(this.nickName?.value);
    profileReq.setPreferredLanguage(this.preferredLanguage?.value);
    profileReq.setGender(this.gender?.value);

    const humanReq = new AddHumanUserRequest();
    humanReq.setUserName(this.userName?.value);
    humanReq.setProfile(profileReq);

    const emailreq = new AddHumanUserRequest.Email();
    emailreq.setEmail(this.email?.value);
    emailreq.setIsEmailVerified(this.isVerified?.value);
    humanReq.setEmail(emailreq);

    if (this.usePassword && this.password?.value) {
      humanReq.setInitialPassword(this.password.value);
    }

    if (this.phone && this.phone.value) {
      humanReq.setPhone(new AddHumanUserRequest.Phone().setPhone(this.phone.value));
    }

    this.mgmtService
      .addHumanUser(humanReq)
      .then((data) => {
        this.loading = false;
        this.toast.showInfo('USER.TOAST.CREATED', true);
        this.router.navigate(['users', data.userId]);
      })
      .catch((error) => {
        this.loading = false;
        this.toast.showError(error);
      });
  }

  ngOnDestroy(): void {
    this.destroyed$.next();
    this.destroyed$.complete();
  }

  public get email(): AbstractControl | null {
    return this.userForm.get('email');
  }
  public get isVerified(): AbstractControl | null {
    return this.userForm.get('isVerified');
  }
  public get userName(): AbstractControl | null {
    return this.userForm.get('userName');
  }
  public get firstName(): AbstractControl | null {
    return this.userForm.get('firstName');
  }
  public get lastName(): AbstractControl | null {
    return this.userForm.get('lastName');
  }
  public get nickName(): AbstractControl | null {
    return this.userForm.get('nickName');
  }
  public get gender(): AbstractControl | null {
    return this.userForm.get('gender');
  }
  public get preferredLanguage(): AbstractControl | null {
    return this.userForm.get('preferredLanguage');
  }
  public get phone(): AbstractControl | null {
    return this.userForm.get('phone');
  }

  public get password(): AbstractControl | null {
    return this.pwdForm.get('password');
  }
  public get confirmPassword(): AbstractControl | null {
    return this.pwdForm.get('confirmPassword');
  }
  private envSuffix(): string {
    if (this.userLoginMustBeDomain && this.primaryDomain?.domainName) {
      return `@${this.primaryDomain.domainName}`;
    } else {
      return '';
    }
  }

  public get suffixPadding(): string | undefined {
    if (this.suffix?.nativeElement.offsetWidth) {
      return `${(this.suffix.nativeElement as HTMLElement).offsetWidth + 10}px`;
    } else {
      return;
    }
  }
}
