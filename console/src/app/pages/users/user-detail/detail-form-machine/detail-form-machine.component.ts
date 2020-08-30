import { Component, EventEmitter, Input, OnDestroy, OnInit, Output } from '@angular/core';
import { AbstractControl, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Subscription } from 'rxjs';

import { UserView } from '../../../../proto/generated/management_pb';

@Component({
    selector: 'app-detail-form-machine',
    templateUrl: './detail-form-machine.component.html',
    styleUrls: ['./detail-form-machine.component.scss'],
})
export class DetailFormMachineComponent implements OnInit, OnDestroy {
    @Input() public username!: string;
    @Input() public user!: UserView;
    @Input() public disabled: boolean = false;
    @Output() public submitData: EventEmitter<any> = new EventEmitter<any>();

    public profileForm!: FormGroup;

    private sub: Subscription = new Subscription();

    constructor(private fb: FormBuilder) {
        this.profileForm = this.fb.group({
            userName: [{ value: '', disabled: true }, [
                Validators.required,
            ]],
            firstName: [{ value: '', disabled: this.disabled }, Validators.required],
            lastName: [{ value: '', disabled: this.disabled }, Validators.required],
            nickName: [{ value: '', disabled: this.disabled }],
            gender: [{ value: 0 }, { disabled: this.disabled }],
            preferredLanguage: [{ value: '', disabled: this.disabled }],
        });
    }

    public ngOnInit(): void {
        this.profileForm.patchValue({ userName: this.username, ...this.user });
    }

    public ngOnDestroy(): void {
        this.sub.unsubscribe();
    }

    public submitForm(): void {
        this.submitData.emit(this.profileForm.value);
    }

    public get userName(): AbstractControl | null {
        return this.profileForm.get('userName');
    }

    public get firstName(): AbstractControl | null {
        return this.profileForm.get('firstName');
    }
    public get lastName(): AbstractControl | null {
        return this.profileForm.get('lastName');
    }
    public get nickName(): AbstractControl | null {
        return this.profileForm.get('nickName');
    }
    public get gender(): AbstractControl | null {
        return this.profileForm.get('gender');
    }
    public get preferredLanguage(): AbstractControl | null {
        return this.profileForm.get('preferredLanguage');
    }

}
