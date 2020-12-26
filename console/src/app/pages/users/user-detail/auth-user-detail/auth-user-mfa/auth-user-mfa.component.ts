import { Component, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatSort } from '@angular/material/sort';
import { MatTable, MatTableDataSource } from '@angular/material/table';
import { BehaviorSubject, Observable } from 'rxjs';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import { MfaOtpResponse, MFAState, MfaType, MultiFactor } from 'src/app/proto/generated/zitadel/auth_pb';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ToastService } from 'src/app/services/toast.service';

import { DialogOtpComponent } from '../dialog-otp/dialog-otp.component';

@Component({
    selector: 'app-auth-user-mfa',
    templateUrl: './auth-user-mfa.component.html',
    styleUrls: ['./auth-user-mfa.component.scss'],
})
export class AuthUserMfaComponent implements OnInit, OnDestroy {
    public displayedColumns: string[] = ['type', 'state', 'actions'];
    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();

    @ViewChild(MatTable) public table!: MatTable<MultiFactor.AsObject>;
    @ViewChild(MatSort) public sort!: MatSort;
    public dataSource!: MatTableDataSource<MultiFactor.AsObject>;

    public MfaType: any = MfaType;
    public MFAState: any = MFAState;

    public error: string = '';
    public otpAvailable: boolean = false;
    constructor(private service: GrpcAuthService, private toast: ToastService, private dialog: MatDialog) { }

    public ngOnInit(): void {
        this.getOTP();
    }

    public ngOnDestroy(): void {
        this.loadingSubject.complete();
    }

    public addOTP(): void {
        this.service.AddMfaOTP().then((otpresp) => {
            const otp: MfaOtpResponse.AsObject = otpresp.toObject();
            const dialogRef = this.dialog.open(DialogOtpComponent, {
                data: otp.url,
                width: '400px',
            });

            dialogRef.afterClosed().subscribe((code) => {
                if (code) {
                    this.service.VerifyMfaOTP(code).then(() => {
                        this.getOTP();
                    });
                }
            });
        }, error => {
            this.toast.showError(error);
        });
    }

    public getOTP(): void {
        this.service.GetMyMfas().then(mfas => {
            this.dataSource = new MatTableDataSource(mfas.toObject().mfasList);
            this.dataSource.sort = this.sort;

            const index = mfas.toObject().mfasList.findIndex(mfa => mfa.type === MfaType.MFATYPE_OTP);
            if (index === -1) {
                this.otpAvailable = true;
            }
        }).catch(error => {
            this.error = error.message;
        });
    }

    public deleteMFA(type: MfaType): void {
        const dialogRef = this.dialog.open(WarnDialogComponent, {
            data: {
                confirmKey: 'ACTIONS.DELETE',
                cancelKey: 'ACTIONS.CANCEL',
                titleKey: 'USER.MFA.DIALOG.OTP_DELETE_TITLE',
                descriptionKey: 'USER.MFA.DIALOG.OTP_DELETE_DESCRIPTION',
            },
            width: '400px',
        });

        dialogRef.afterClosed().subscribe(resp => {
            if (resp) {
                if (type === MfaType.MFATYPE_OTP) {
                    this.service.RemoveMfaOTP().then(() => {
                        this.toast.showInfo('USER.TOAST.OTPREMOVED', true);

                        const index = this.dataSource.data.findIndex(mfa => mfa.type === type);
                        if (index > -1) {
                            this.dataSource.data.splice(index, 1);
                        }
                        this.getOTP();
                    }).catch(error => {
                        this.toast.showError(error);
                    });
                }
            }
        });
    }
}
