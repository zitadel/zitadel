import { Component, OnDestroy } from '@angular/core';
import { AbstractControl, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { Subscription } from 'rxjs';
import { CreateUserRequest, Gender, User } from 'src/app/proto/generated/management_pb';
import { MgmtUserService } from 'src/app/services/mgmt-user.service';
import { ToastService } from 'src/app/services/toast.service';

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

    private sub: Subscription = new Subscription();

    constructor(private router: Router, private toast: ToastService, public userService: MgmtUserService,
        private fb: FormBuilder) {

        this.userForm = this.fb.group({
            email: ['', [Validators.required, Validators.email]],
            userName: ['', [Validators.required, Validators.minLength(2)]],
            firstName: ['', Validators.required],
            lastName: ['', Validators.required],
            nickName: [''],
            gender: [Gender.GENDER_UNSPECIFIED],
            preferredLanguage: [''],
            phone: [''],
            streetAddress: [''],
            postalCode: [''],
            locality: [''],
            region: [''],
            country: [''],
        });
    }

    public createUser(): void {
        this.user = this.userForm.value;
        console.log(this.user);

        this.userService
            .CreateUser(this.user)
            .then((data: User) => {
                this.toast.showInfo('User created');
                this.router.navigate(['users', data.getId()]);
            })
            .catch(data => {
                this.toast.showError(data.message);
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
}
