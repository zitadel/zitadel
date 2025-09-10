import { Location } from '@angular/common';
import { Component } from '@angular/core';
import { AbstractControl, UntypedFormBuilder, UntypedFormGroup } from '@angular/forms';
import { MatSlideToggleChange } from '@angular/material/slide-toggle';
import { Router } from '@angular/router';
import { passwordConfirmValidator, requiredValidator } from 'src/app/modules/form-field/validators/validators';
import { Gender } from 'src/app/proto/generated/zitadel/user_pb';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ToastService } from 'src/app/services/toast.service';
import { LanguagesService } from 'src/app/services/languages.service';
import { PasswordComplexityPolicy } from '@zitadel/proto/zitadel/policy_pb';
import { NewMgmtService } from 'src/app/services/new-mgmt.service';
import { PasswordComplexityValidatorFactoryService } from 'src/app/services/password-complexity-validator-factory.service';
import { injectMutation } from '@tanstack/angular-query-experimental';
import { NewOrganizationService } from '../../services/new-organization.service';
import { MessageInitShape } from '@bufbuild/protobuf';
import { SetUpOrgRequestSchema } from '@zitadel/proto/zitadel/admin_pb';

@Component({
  selector: 'cnsl-org-create',
  templateUrl: './org-create.component.html',
  styleUrls: ['./org-create.component.scss'],
})
export class OrgCreateComponent {
  protected orgForm = this.fb.group({
    name: ['', [requiredValidator]],
    domain: [''],
  });

  protected userForm?: UntypedFormGroup;
  protected pwdForm?: UntypedFormGroup;

  protected readonly genders: Gender[] = [Gender.GENDER_FEMALE, Gender.GENDER_MALE, Gender.GENDER_UNSPECIFIED];

  protected policy?: PasswordComplexityPolicy;
  protected usePassword: boolean = false;

  protected forSelf: boolean = true;

  protected readonly setupOrgMutation = injectMutation(this.newOrganizationService.setupOrgMutationOptions);
  protected readonly addOrgMutation = injectMutation(this.newOrganizationService.addOrgMutationOptions);

  constructor(
    private readonly router: Router,
    private readonly toast: ToastService,
    private readonly location: Location,
    private readonly fb: UntypedFormBuilder,
    private readonly newMgmtService: NewMgmtService,
    private readonly passwordComplexityValidatorFactory: PasswordComplexityValidatorFactoryService,
    public readonly langSvc: LanguagesService,
    private readonly newOrganizationService: NewOrganizationService,
    breadcrumbService: BreadcrumbService,
  ) {
    const instanceBread = new Breadcrumb({
      type: BreadcrumbType.INSTANCE,
      name: 'Instance',
      routerLink: ['/instance'],
    });

    breadcrumbService.setBreadcrumb([instanceBread]);
    this.initForm();
  }

  public createSteps: number = 2;
  public currentCreateStep: number = 1;

  public async finish(): Promise<void> {
    const req: MessageInitShape<typeof SetUpOrgRequestSchema> = {
      org: {
        name: this.name?.value,
        domain: this.domain?.value,
      },
      user: {
        case: 'human',
        value: {
          email: {
            email: this.email?.value,
            isEmailVerified: this.isVerified?.value,
          },
          userName: this.userName?.value,
          profile: {
            firstName: this.firstName?.value,
            lastName: this.lastName?.value,
            nickName: this.nickName?.value,
            gender: this.gender?.value,
            preferredLanguage: this.preferredLanguage?.value,
          },
          password: this.usePassword && this.password ? this.password.value : undefined,
        },
      },
    };

    try {
      await this.setupOrgMutation.mutateAsync(req);
      await this.router.navigate(['/orgs']);
    } catch (error) {
      this.toast.showError(error);
    }
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
    if (this.usePassword) {
      this.newMgmtService.getDefaultPasswordComplexityPolicy().then((data) => {
        this.pwdForm = this.fb.group({
          password: ['', this.passwordComplexityValidatorFactory.buildValidators(data.policy)],
          confirmPassword: ['', [requiredValidator, passwordConfirmValidator()]],
        });
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

  public async createOrgForSelf() {
    if (!this.name?.value) {
      return;
    }
    try {
      await this.addOrgMutation.mutateAsync(this.name.value);
      await this.router.navigate(['/orgs']);
    } catch (error) {
      this.toast.showError(error);
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
    this.location.back();
  }
}
