import { Location } from '@angular/common';
import { Component, DestroyRef, ElementRef, OnInit, ViewChild } from '@angular/core';
import { FormBuilder, FormControl, ValidatorFn } from '@angular/forms';
import { Router } from '@angular/router';
import { debounceTime, defer, of, Observable, shareReplay, forkJoin, ObservedValueOf, EMPTY, ReplaySubject } from 'rxjs';
import { PasswordComplexityPolicy } from 'src/app/proto/generated/zitadel/policy_pb';
import { Gender } from 'src/app/proto/generated/zitadel/user_pb';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';
import { CountryCallingCodesService, CountryPhoneCode } from 'src/app/services/country-calling-codes.service';
import { formatPhone } from 'src/app/utils/formatPhone';
import {
  containsLowerCaseValidator,
  containsNumberValidator,
  containsSymbolValidator,
  containsUpperCaseValidator,
  emailValidator,
  minLengthValidator,
  passwordConfirmValidator,
  phoneValidator,
  requiredValidator,
} from 'src/app/modules/form-field/validators/validators';
import { LanguagesService } from 'src/app/services/languages.service';
import { AddHumanUserRequest } from 'src/app/proto/generated/zitadel/management_pb';
import { catchError, map, startWith } from 'rxjs/operators';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';

@Component({
  selector: 'cnsl-user-create',
  templateUrl: './user-create.component.html',
  styleUrls: ['./user-create.component.scss'],
})
export class UserCreateComponent implements OnInit {
  protected readonly genders: Gender[] = [Gender.GENDER_FEMALE, Gender.GENDER_MALE, Gender.GENDER_UNSPECIFIED];
  protected selected: CountryPhoneCode | undefined = {
    countryCallingCode: '1',
    countryCode: 'US',
    countryName: 'United States of America',
  };
  protected readonly countryPhoneCodes: CountryPhoneCode[];

  protected loading = false;

  private readonly suffix$ = new ReplaySubject<HTMLSpanElement>(1);
  @ViewChild('suffix') public set suffix(suffix: ElementRef<HTMLSpanElement>) {
    this.suffix$.next(suffix.nativeElement);
  }

  protected usePassword: boolean = false;
  protected readonly envSuffix$: Observable<string>;
  protected readonly userForm: ReturnType<typeof this.buildUserForm>;
  protected readonly pwdForm$: ReturnType<typeof this.buildPwdForm>;
  protected readonly passwordComplexityPolicy$: Observable<PasswordComplexityPolicy.AsObject | undefined>;
  protected readonly suffixPadding$: Observable<string>;

  constructor(
    private readonly router: Router,
    private readonly toast: ToastService,
    private readonly fb: FormBuilder,
    private readonly mgmtService: ManagementService,
    private readonly destroyRef: DestroyRef,
    private readonly breadcrumbService: BreadcrumbService,
    protected readonly location: Location,
    protected readonly langSvc: LanguagesService,
    countryCallingCodesService: CountryCallingCodesService,
  ) {
    this.envSuffix$ = this.getEnvSuffix();
    this.suffixPadding$ = this.getSuffixPadding();
    this.passwordComplexityPolicy$ = this.getPasswordComplexityPolicy().pipe(shareReplay({ refCount: true, bufferSize: 1 }));

    this.userForm = this.buildUserForm();
    this.pwdForm$ = this.buildPwdForm(this.passwordComplexityPolicy$);

    this.countryPhoneCodes = countryCallingCodesService.getCountryCallingCodes();
  }

  ngOnInit(): void {
    this.watchPhoneChanges();

    this.breadcrumbService.setBreadcrumb([
      new Breadcrumb({
        type: BreadcrumbType.ORG,
        routerLink: ['/org'],
      }),
    ]);
  }

  private getEnvSuffix() {
    const domainPolicy$ = defer(() => this.mgmtService.getDomainPolicy());
    const orgDomains$ = defer(() => this.mgmtService.listOrgDomains());

    return forkJoin([domainPolicy$, orgDomains$]).pipe(
      map(([policy, domains]) => {
        const userLoginMustBeDomain = policy.policy?.userLoginMustBeDomain;
        const primaryDomain = domains.resultList.find((resp) => resp.isPrimary);
        if (userLoginMustBeDomain && primaryDomain) {
          return `@${primaryDomain.domainName}`;
        } else {
          return '';
        }
      }),
      catchError(() => of('')),
    );
  }

  private getSuffixPadding() {
    return this.suffix$.pipe(
      map((suffix) => `${suffix.offsetWidth + 10}px`),
      startWith('10px'),
    );
  }

  private getPasswordComplexityPolicy() {
    return defer(() => this.mgmtService.getPasswordComplexityPolicy()).pipe(
      map(({ policy }) => policy),
      catchError((error) => {
        this.toast.showError(error);
        return EMPTY;
      }),
    );
  }

  private buildUserForm() {
    return this.fb.group({
      email: new FormControl('', { nonNullable: true, validators: [requiredValidator, emailValidator] }),
      userName: new FormControl('', { nonNullable: true, validators: [requiredValidator, minLengthValidator(2)] }),
      firstName: new FormControl('', { nonNullable: true, validators: [requiredValidator] }),
      lastName: new FormControl('', { nonNullable: true, validators: [requiredValidator] }),
      nickName: new FormControl('', { nonNullable: true }),
      gender: new FormControl(Gender.GENDER_UNSPECIFIED, { nonNullable: true, validators: [requiredValidator] }),
      preferredLanguage: new FormControl('', { nonNullable: true }),
      phone: new FormControl('', { nonNullable: true, validators: [phoneValidator] }),
      emailVerified: new FormControl(false, { nonNullable: true }),
    });
  }

  private buildPwdForm(passwordComplexityPolicy$: Observable<PasswordComplexityPolicy.AsObject | undefined>) {
    return passwordComplexityPolicy$.pipe(
      map((policy) => {
        const validators: ValidatorFn[] = [requiredValidator];
        if (policy?.minLength) {
          validators.push(minLengthValidator(policy.minLength));
        }
        if (policy?.hasLowercase) {
          validators.push(containsLowerCaseValidator);
        }
        if (policy?.hasUppercase) {
          validators.push(containsUpperCaseValidator);
        }
        if (policy?.hasNumber) {
          validators.push(containsNumberValidator);
        }
        if (policy?.hasSymbol) {
          validators.push(containsSymbolValidator);
        }
        return this.fb.group({
          password: new FormControl('', { nonNullable: true, validators }),
          confirmPassword: new FormControl('', {
            nonNullable: true,
            validators: [requiredValidator, passwordConfirmValidator()],
          }),
        });
      }),
    );
  }

  private watchPhoneChanges(): void {
    const phone = this.userForm.controls.phone;

    phone.valueChanges.pipe(debounceTime(200), takeUntilDestroyed(this.destroyRef)).subscribe((value) => {
      const phoneNumber = formatPhone(value);
      if (phoneNumber) {
        this.selected = this.countryPhoneCodes.find((code) => code.countryCode === phoneNumber.country);
        phone.setValue(phoneNumber.phone);
      }
    });
  }

  protected async createUser(pwdForm: ObservedValueOf<typeof this.pwdForm$>): Promise<void> {
    this.loading = true;

    const controls = this.userForm.controls;
    const profileReq = new AddHumanUserRequest.Profile();
    profileReq.setFirstName(controls.firstName.value);
    profileReq.setLastName(controls.lastName.value);
    profileReq.setNickName(controls.nickName.value);
    profileReq.setPreferredLanguage(controls.preferredLanguage.value);
    profileReq.setGender(controls.gender.value);

    const humanReq = new AddHumanUserRequest();
    humanReq.setUserName(controls.userName.value);
    humanReq.setProfile(profileReq);

    const emailreq = new AddHumanUserRequest.Email();
    emailreq.setEmail(controls.email.value);
    emailreq.setIsEmailVerified(controls.emailVerified.value);
    humanReq.setEmail(emailreq);

    if (this.usePassword) {
      humanReq.setInitialPassword(pwdForm.controls.password.value);
    }

    if (controls.phone.value) {
      // Try to parse number and format it according to country
      const phoneNumber = formatPhone(controls.phone.value);
      if (phoneNumber) {
        this.selected = this.countryPhoneCodes.find((code) => code.countryCode === phoneNumber.country);
        humanReq.setPhone(new AddHumanUserRequest.Phone().setPhone(phoneNumber.phone));
      }
    }

    try {
      const data = await this.mgmtService.addHumanUser(humanReq);
      this.toast.showInfo('USER.TOAST.CREATED', true);
      await this.router.navigate(['users', data.userId], { queryParams: { new: true } });
    } catch (error) {
      this.toast.showError(error);
    } finally {
      this.loading = false;
    }
  }

  protected setCountryCallingCode(): void {
    let value = this.userForm.controls.phone.value;
    this.countryPhoneCodes.forEach((code) => (value = value.replace(`+${code.countryCallingCode}`, '')));
    value = value.trim();

    this.userForm.controls.phone.setValue('+' + this.selected?.countryCallingCode + ' ' + value);
  }

  protected compareCountries(i1: CountryPhoneCode, i2: CountryPhoneCode) {
    return (
      i1 &&
      i2 &&
      i1.countryCallingCode === i2.countryCallingCode &&
      i1.countryCode == i2.countryCode &&
      i1.countryName == i2.countryName
    );
  }
}
