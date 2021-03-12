import { ChangeDetectorRef, Component, OnDestroy, ViewChild } from '@angular/core';
import { AbstractControl, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { Subscription } from 'rxjs';
import { AddHumanUserRequest } from 'src/app/proto/generated/zitadel/management_pb';
import { Domain } from 'src/app/proto/generated/zitadel/org_pb';
import { Gender } from 'src/app/proto/generated/zitadel/user_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

function noEmailValidator(c: AbstractControl): any {
    const EMAIL_REGEXP: RegExp = /^((?!@).)*$/gm;
    if (!c.parent || !c) {
        return;
    }
    const username = c.parent.get('userName');

    if (!username) {
        return;
    }

    return EMAIL_REGEXP.test(username.value) ? null : {
        noEmailValidator: {
            valid: false,
        },
    };
}

@Component({
    selector: 'app-user-create',
    templateUrl: './user-create.component.html',
    styleUrls: ['./user-create.component.scss'],
})
export class UserCreateComponent implements OnDestroy {
    public user: AddHumanUserRequest.AsObject = new AddHumanUserRequest().toObject();
    public genders: Gender[] = [Gender.GENDER_FEMALE, Gender.GENDER_MALE, Gender.GENDER_UNSPECIFIED];
    public languages: string[] = ['de', 'en'];
    public userForm!: FormGroup;
    public envSuffixLabel: string = '';
    private sub: Subscription = new Subscription();

    public userLoginMustBeDomain: boolean = false;
    public loading: boolean = false;

    @ViewChild('suffix') public suffix!: any;
    private primaryDomain!: Domain.AsObject;

    constructor(
        private router: Router,
        private toast: ToastService,
        private fb: FormBuilder,
        private mgmtService: ManagementService,
        private changeDetRef: ChangeDetectorRef,
    ) {
        this.loading = true;
        this.loadOrg();
        this.mgmtService.getOrgIAMPolicy().then((resp) => {
            if (resp.policy?.userLoginMustBeDomain) {
                this.userLoginMustBeDomain = resp.policy.userLoginMustBeDomain;
            }
            this.initForm();
            this.loading = false;
            this.envSuffixLabel = this.envSuffix();
            this.changeDetRef.detectChanges();
        }).catch(error => {
            console.error(error);
            this.initForm();
            this.loading = false;
            this.envSuffixLabel = this.envSuffix();
            this.changeDetRef.detectChanges();
        });
    }

    private async loadOrg(): Promise<void> {
        const domains = (await this.mgmtService.listOrgDomains());
        const found = domains.resultList.find(resp => resp.isPrimary);
        if (found) {
            this.primaryDomain = found;
        }
    }

    private initForm(): void {
        this.userForm = this.fb.group({
            email: ['', [Validators.required, Validators.email]],
            userName: ['',
                [
                    Validators.required,
                    Validators.minLength(2),
                    this.userLoginMustBeDomain ? noEmailValidator : Validators.email,
                ],
            ],
            firstName: ['', Validators.required],
            lastName: ['', Validators.required],
            nickName: [''],
            gender: [Gender.GENDER_UNSPECIFIED],
            preferredLanguage: [''],
            phone: [''],
        });

    }

    public logsuff(): void {
        console.log((this.suffix.nativeElement as HTMLElement), (this.suffix.nativeElement as HTMLElement).offsetWidth);
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

        humanReq.setEmail(new AddHumanUserRequest.Email().setEmail(this.email?.value));
        humanReq.setPhone(new AddHumanUserRequest.Phone().setPhone(this.phone?.value));

        this.mgmtService
            .addHumanUser(humanReq)
            .then((data) => {
                this.loading = false;
                this.toast.showInfo('USER.TOAST.CREATED', true);
                this.router.navigate(['users', data.userId]);
            })
            .catch(error => {
                this.loading = false;
                this.toast.showError(error);
            });
    }

    ngOnDestroy(): void {
        this.sub.unsubscribe();
    }

    public get email(): AbstractControl | null {
        return this.userForm.get('email');
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
        }
    }
}
