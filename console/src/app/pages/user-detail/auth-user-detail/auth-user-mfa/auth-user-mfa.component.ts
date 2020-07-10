import { Component, OnDestroy, OnInit } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { BehaviorSubject, Observable } from 'rxjs';
import { MfaOtpResponse, MFAState, MfaType, MultiFactor } from 'src/app/proto/generated/auth_pb';
import { AuthUserService } from 'src/app/services/auth-user.service';
import { ToastService } from 'src/app/services/toast.service';

import { DialogOtpComponent } from '../dialog-otp/dialog-otp.component';

@Component({
    selector: 'app-auth-user-mfa',
    templateUrl: './auth-user-mfa.component.html',
    styleUrls: ['./auth-user-mfa.component.scss'],
})
export class AuthUserMfaComponent implements OnInit, OnDestroy {
    public mfaSubject: BehaviorSubject<MultiFactor.AsObject[]> = new BehaviorSubject<MultiFactor.AsObject[]>([]);
    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();

    public MfaType: any = MfaType;
    public MFAState: any = MFAState;

    public error: string = '';
    constructor(private userService: AuthUserService, private toast: ToastService, private dialog: MatDialog) { }

    public ngOnInit(): void {
        this.getOTP();
    }

    public ngOnDestroy(): void {
        this.mfaSubject.complete();
        this.loadingSubject.complete();
    }

    public addOTP(): void {
        this.userService.AddMfaOTP().then((otpresp) => {
            const otp: MfaOtpResponse.AsObject = otpresp.toObject();
            const dialogRef = this.dialog.open(DialogOtpComponent, {
                data: otp.url,
                width: '400px',
            });

            dialogRef.afterClosed().subscribe((code) => {
                if (code) {
                    this.userService.VerifyMfaOTP(code).then((res) => {
                        // TODO: show state
                    });
                }
            });
        }, error => {
            this.toast.showError(error);
        });
    }

    public getOTP(): void {
        this.userService.GetMyMfas().then(mfas => {
            this.mfaSubject.next(mfas.toObject().mfasList);
        }).catch(error => {
            console.error(error);
            this.error = error.message;
        });
    }

    public deleteMFA(type: MfaType): void {
        if (type === MfaType.MFATYPE_OTP) {
            this.userService.RemoveMfaOTP().then(() => {
                this.toast.showInfo('USER.TOAST.OTPREMOVED', true);

                const index = this.mfaSubject.value.findIndex(mfa => mfa.type === type);
                if (index > -1) {
                    const newValues = this.mfaSubject.value;
                    newValues.splice(index, 1);
                    this.mfaSubject.next(newValues);
                }

            }).catch(error => {
                this.toast.showError(error);
            });
        }
    }
}
