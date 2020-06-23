import { Component, Input, OnDestroy, OnInit } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { MFAState, MfaType, MultiFactor, UserView } from 'src/app/proto/generated/management_pb';
import { MgmtUserService } from 'src/app/services/mgmt-user.service';
import { ToastService } from 'src/app/services/toast.service';


export interface MFAItem {
    name: string;
    verified: boolean;
}

@Component({
    selector: 'app-user-mfa',
    templateUrl: './user-mfa.component.html',
    styleUrls: ['./user-mfa.component.scss'],
})
export class UserMfaComponent implements OnInit, OnDestroy {
    @Input() private user!: UserView.AsObject;
    public mfaSubject: BehaviorSubject<MultiFactor.AsObject[]> = new BehaviorSubject<MultiFactor.AsObject[]>([]);
    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();

    public MfaType: any = MfaType;
    public MFAState: any = MFAState;
    constructor(private mgmtUserService: MgmtUserService,
        private toast: ToastService) { }

    public ngOnInit(): void {
        console.log(this.user);
        this.getOTP();
    }

    public ngOnDestroy(): void {
        this.mfaSubject.complete();
        this.loadingSubject.complete();
    }

    public getOTP(): void {
        console.log('otp', this.user);
        this.mgmtUserService.getUserMfas(this.user.id).then(mfas => {
            this.mfaSubject.next(mfas.toObject().mfasList);
            console.log(mfas.toObject());
        }).catch(error => {
            console.error(error);
            this.toast.showError(error.message);
        });
    }

    // public deleteMFA(type: MfaType): void {
    //     if (type === MfaType.MFATYPE_OTP) {
    //         this.userService.RemoveMfaOTP().then(() => {
    //             this.toast.showInfo('OTP Deleted');

    //             const index = this.mfaSubject.value.findIndex(mfa => mfa.type === type);
    //             if (index > -1) {
    //                 const newValues = this.mfaSubject.value;
    //                 newValues.splice(index, 1);
    //                 this.mfaSubject.next(newValues);
    //             }

    //         }).catch(error => {
    //             this.toast.showError(error.message);
    //         });
    //     }
    // }
}
