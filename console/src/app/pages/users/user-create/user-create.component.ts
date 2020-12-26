import { ChangeDetectorRef, Component, OnDestroy, ViewChild } from '@angular/core';
import { AbstractControl, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { Subscription } from 'rxjs';
import {
    CreateHumanRequest,
    CreateUserRequest,
    Gender,
    OrgDomain,
    UserResponse,
} from 'src/app/proto/generated/zitadel/management_pb';
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
    public user: CreateUserRequest.AsObject = new CreateUserRequest().toObject();
    public genders: Gender[] = [Gender.GENDER_FEMALE, Gender.GENDER_MALE, Gender.GENDER_UNSPECIFIED];
    public languages: string[] = ['de', 'en'];
    public userForm!: FormGroup;
    public envSuffixLabel: string = '';
    private sub: Subscription = new Subscription();

    public userLoginMustBeDomain: boolean = false;
    public loading: boolean = false;

    @ViewChild('suffix') public suffix!: any;
    private primaryDomain!: OrgDomain.AsObject;

    constructor(
        private router: Router,
        private toast: ToastService,
        private fb: FormBuilder,
        private mgmtService: ManagementService,
        private changeDetRef: ChangeDetectorRef,
    ) {
        this.loading = true;
        this.loadOrg();
        this.mgmtService.GetMyOrgIamPolicy().then((iampolicy) => {
            this.userLoginMustBeDomain = iampolicy.toObject().userLoginMustBeDomain;
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
        const domains = (await this.mgmtService.SearchMyOrgDomains().then(doms => doms.toObject()));
        const found = domains.resultList.find(domain => domain.primary);
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

        const humanReq = new CreateHumanRequest();
        humanReq.setFirstName(this.firstName?.value);
        humanReq.setLastName(this.lastName?.value);
        humanReq.setNickName(this.nickName?.value);
        humanReq.setPreferredLanguage(this.preferredLanguage?.value);
        humanReq.setEmail(this.email?.value);
        humanReq.setPhone(this.phone?.value);
        humanReq.setGender(this.gender?.value);
        humanReq.setCountry(this.country?.value);

        this.mgmtService
            .CreateUserHuman(this.userName?.value, humanReq)
            .then((data: UserResponse) => {
                this.loading = false;
                this.toast.showInfo('USER.TOAST.CREATED', true);
                this.router.navigate(['users', data.getId()]);
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
    public get streetAddress(): AbstractControl | null {
        return this.userForm.get('streetAddress');
    }
    public get postalCode(): AbstractControl | null {
        return this.userForm.get('postalCode');
    }
    public get locality(): AbstractControl | null {
        return this.userForm.get('locality');
    }
    public get region(): AbstractControl | null {
        return this.userForm.get('region');
    }
    public get country(): AbstractControl | null {
        return this.userForm.get('country');
    }

    private envSuffix(): string {
        if (this.userLoginMustBeDomain && this.primaryDomain?.domain) {
            return `@${this.primaryDomain.domain}`;
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
