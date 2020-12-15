import { Component, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatSort } from '@angular/material/sort';
import { MatTable, MatTableDataSource } from '@angular/material/table';
import { BehaviorSubject, Observable } from 'rxjs';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import { MfaOtpResponse, MFAState, MfaType, MultiFactor, WebAuthNResponse } from 'src/app/proto/generated/auth_pb';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ToastService } from 'src/app/services/toast.service';

import { DialogOtpComponent } from '../dialog-otp/dialog-otp.component';
import { DialogU2FComponent } from '../dialog-u2f/dialog-u2f.component';

export function _base64ToArrayBuffer(base64: string): any {
    const binaryString = atob(base64);
    const len = binaryString.length;
    const bytes = new Uint8Array(len);
    for (let i = 0; i < len; i++) {
        bytes[i] = binaryString.charCodeAt(i);
    }
    return bytes.buffer;
}

export interface WebAuthNOptions {
    challenge: string;
    rp: { name: string, id: string; };
    user: { name: string, id: string, displayName: string; };
    pubKeyCredParams: any;
    authenticatorSelection: { userVerification: string; };
    timeout: number;
    attestation: string;
}

@Component({
    selector: 'app-auth-user-mfa',
    templateUrl: './auth-user-mfa.component.html',
    styleUrls: ['./auth-user-mfa.component.scss'],
})
export class AuthUserMfaComponent implements OnInit, OnDestroy {
    public displayedColumns: string[] = ['type', 'attr', 'state', 'actions'];
    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();

    @ViewChild(MatTable) public table!: MatTable<MultiFactor.AsObject>;
    @ViewChild(MatSort) public sort!: MatSort;
    public dataSource!: MatTableDataSource<MultiFactor.AsObject>;

    public MfaType: any = MfaType;
    public MFAState: any = MFAState;

    public error: string = '';
    public otpAvailable: boolean = false;

    constructor(private service: GrpcAuthService,
        private toast: ToastService,
        private dialog: MatDialog) { }

    public ngOnInit(): void {
        this.getMFAs();
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
                        this.getMFAs();
                    });
                }
            });
        }, error => {
            this.toast.showError(error);
        });
    }

    public verifyU2f(): void {

    }

    public addU2F(): void {
        this.service.AddMyMfaU2F().then((u2fresp) => {
            const webauthn: WebAuthNResponse.AsObject = u2fresp.toObject();
            const credOptions: CredentialCreationOptions = JSON.parse(atob(webauthn.publicKey as string));

            if (credOptions.publicKey?.challenge) {
                credOptions.publicKey.challenge = _base64ToArrayBuffer(credOptions.publicKey.challenge as any);
                credOptions.publicKey.user.id = _base64ToArrayBuffer(credOptions.publicKey.user.id as any);
                const dialogRef = this.dialog.open(DialogU2FComponent, {
                    width: '400px',
                    data: {
                        credOptions,
                    },
                });

                dialogRef.afterClosed().subscribe(done => {
                    if (done) {
                        this.getMFAs();
                    } else {
                        this.getMFAs();
                    }
                });
            }

        }, error => {
            this.toast.showError(error);
        });
    }

    public getMFAs(): void {
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

    public deleteMFA(type: MfaType, id?: string): void {
        const dialogRef = this.dialog.open(WarnDialogComponent, {
            data: {
                confirmKey: 'ACTIONS.DELETE',
                cancelKey: 'ACTIONS.CANCEL',
                titleKey: 'USER.MFA.DIALOG.MFA_DELETE_TITLE',
                descriptionKey: 'USER.MFA.DIALOG.MFA_DELETE_DESCRIPTION',
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
                        this.getMFAs();
                    }).catch(error => {
                        this.toast.showError(error);
                    });
                } else if (type === MfaType.MFATYPE_U2F && id) {
                    this.service.RemoveMyMfaU2F(id).then(() => {
                        this.toast.showInfo('USER.TOAST.U2FREMOVED', true);

                        const index = this.dataSource.data.findIndex(mfa => mfa.type === type);
                        if (index > -1) {
                            this.dataSource.data.splice(index, 1);
                        }
                        this.getMFAs();
                    }).catch(error => {
                        this.toast.showError(error);
                    });
                }
            }
        });
    }


}
