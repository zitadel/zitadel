import { Component, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatSort } from '@angular/material/sort';
import { MatTable, MatTableDataSource } from '@angular/material/table';
import { TranslateService } from '@ngx-translate/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import { MfaOtpResponse, MFAState, MfaType, MultiFactor, WebAuthNResponse } from 'src/app/proto/generated/auth_pb';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ToastService } from 'src/app/services/toast.service';

import { DialogOtpComponent } from '../dialog-otp/dialog-otp.component';
import { DialogU2FComponent } from '../dialog-u2f/dialog-u2f.component';

export interface WebAuthNOptions {
  challenge: string;
  rp: {name: string, id: string};
  user: {name: string, id: string, displayName: string};
  pubKeyCredParams: any;
  authenticatorSelection: {userVerification: string};
  timeout: number;
  attestation: string;
}

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
    constructor(private service: GrpcAuthService,
      private toast: ToastService,
      private dialog: MatDialog,
      private translate: TranslateService) { }

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

            console.log(credOptions);
            if (credOptions.publicKey?.challenge) {
              credOptions.publicKey.challenge = this._base64ToArrayBuffer(credOptions.publicKey.challenge as any);
              credOptions.publicKey.user.id = this._base64ToArrayBuffer(credOptions.publicKey.user.id as any);
              console.log(credOptions);
              const dialogRef = this.dialog.open(DialogU2FComponent, {
                width: '400px',
              });

              dialogRef.afterClosed().subscribe(tokenname => {
                if (tokenname && credOptions.publicKey) {
                  navigator.credentials.create(credOptions).then((resp) => {
                    console.log(resp);

                      if (resp &&
                        (resp as any).response.attestationObject &&
                        (resp as any).response.clientDataJSON &&
                        (resp as any).rawId) {

                        const attestationObject = new Uint8Array((resp as any).response.attestationObject);
                        const clientDataJSON = new Uint8Array((resp as any).response.clientDataJSON);
                        const rawId = new Uint8Array((resp as any).rawId);

                        const data = JSON.stringify({
                            id: resp.id,
                            rawId: this._arrayBufferToBase64(rawId),
                            type: resp.type,
                            response: {
                                attestationObject: this._arrayBufferToBase64(attestationObject),
                                clientDataJSON: this._arrayBufferToBase64(clientDataJSON),
                            },
                        });

                        console.log(data);

                        const base64 = btoa(data);
                        console.log(base64);
                        this.service.VerifyMyMfaU2F(base64, tokenname).then(() => {
                          this.translate.get('USER.MFA.U2F_SUCCESS').subscribe(msg => {
                            this.toast.showInfo(msg);
                          });
                        }).catch(error => {
                          this.toast.showError(error);
                        });
                      } else {
                        this.translate.get('USER.MFA.U2F_ERROR').subscribe(msg => {
                          this.toast.showInfo(msg);
                        });
                      }
                  });
                }
              });
            }

        }, error => {
            this.toast.showError(error);
        });
    }

    public getMFAs(): void {
        this.service.GetMyMfas().then(mfas => {
          console.log(mfas.toObject().mfasList);
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

    private _base64ToArrayBuffer(base64: string): any {
      const binaryString = window.atob(base64);
      const len = binaryString.length;
      const bytes = new Uint8Array(len);
      for (let i = 0; i < len; i++) {
          bytes[i] = binaryString.charCodeAt(i);
      }
      return bytes.buffer;
    }

    private _arrayBufferToBase64( buffer: any): string {
      let binary = '';
      const bytes = new Uint8Array( buffer );
      const len = bytes.byteLength;
      for (let i = 0; i < len; i++) {
          binary += String.fromCharCode( bytes[ i ] );
      }
      return window.btoa( binary );
  }
}
