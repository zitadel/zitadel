import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { TranslateService } from '@ngx-translate/core';
import { take } from 'rxjs/operators';
import { Observable } from 'rxjs';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ToastService } from 'src/app/services/toast.service';

import { InfoSectionType } from 'src/app/modules/info-section/info-section.component';
import { _base64ToArrayBuffer } from '../../u2f-util';
import { _arrayBufferToBase64 } from '../u2f_util';

export enum AuthFactorType {
  OTP,
  U2F,
  OTPSMS,
  OTPEMAIL,
}

export type AddAuthFactorDialogData = {
  otp$: Observable<boolean>;
  u2f$: Observable<boolean>;
  otpSms$: Observable<boolean>;
  otpEmail$: Observable<boolean>;
  otpDisabled$: Observable<boolean>;
  otpSmsDisabled$: Observable<boolean>;
  otpEmailDisabled$: Observable<boolean>;
  phoneVerified: boolean;
};

@Component({
  selector: 'cnsl-auth-factor-dialog',
  templateUrl: './auth-factor-dialog.component.html',
  styleUrls: ['./auth-factor-dialog.component.scss'],
})
export class AuthFactorDialogComponent {
  public otpurl: string = '';
  public otpsecret: string = '';

  public otpcode: string = '';

  public u2fname: string = '';
  public u2fCredentialOptions!: CredentialCreationOptions;
  public u2fLoading: boolean = false;
  public u2fError: string = '';

  public phoneVerified: boolean = false;

  AuthFactorType: any = AuthFactorType;
  selectedType!: AuthFactorType;

  public copied: string = '';
  public InfoSectionType: any = InfoSectionType;
  constructor(
    private authService: GrpcAuthService,
    private toast: ToastService,
    private translate: TranslateService,
    public dialogRef: MatDialogRef<AuthFactorDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: AddAuthFactorDialogData,
  ) {}

  closeDialog(code: string = ''): void {
    this.dialogRef.close(code);
  }

  public selectType(type: AuthFactorType): void {
    this.selectedType = type;

    if (type === AuthFactorType.OTP) {
      this.authService
        .addMyMultiFactorOTP()
        .then((otpresp) => {
          this.otpurl = otpresp.url;
          this.otpsecret = otpresp.secret;
        })
        .catch((error) => {
          this.toast.showError(error);
        });
    } else if (type === AuthFactorType.U2F) {
      this.authService
        .addMyMultiFactorU2F()
        .then((u2fresp) => {
          if (u2fresp.key) {
            const credOptions: CredentialCreationOptions = JSON.parse(atob(u2fresp.key?.publicKey as string));

            if (credOptions.publicKey?.challenge) {
              credOptions.publicKey.challenge = _base64ToArrayBuffer(credOptions.publicKey.challenge as any);
              credOptions.publicKey.user.id = _base64ToArrayBuffer(credOptions.publicKey.user.id as any);
              if (credOptions.publicKey.excludeCredentials) {
                credOptions.publicKey.excludeCredentials.map((cred) => {
                  cred.id = _base64ToArrayBuffer(cred.id as any);
                  return cred;
                });
              }
              this.u2fCredentialOptions = credOptions;
            }
          }
        })
        .catch((error) => {
          this.toast.showError(error);
        });
    } else if (type === AuthFactorType.OTPSMS) {
      this.authService
        .addMyAuthFactorOTPSMS()
        .then(() => {
          this.dialogRef.close(true);
          this.translate
            .get('USER.MFA.OTPSMSSUCCESS')
            .pipe(take(1))
            .subscribe((msg) => {
              this.toast.showInfo(msg);
            });
        })
        .catch((error) => {
          this.dialogRef.close(false);
          this.toast.showError(error);
        });
    } else if (type === AuthFactorType.OTPEMAIL) {
      this.authService
        .addMyAuthFactorOTPEmail()
        .then(() => {
          this.dialogRef.close(true);
          this.translate
            .get('USER.MFA.OTPEMAILSUCCESS')
            .pipe(take(1))
            .subscribe((msg) => {
              this.toast.showInfo(msg);
            });
        })
        .catch((error) => {
          this.dialogRef.close(false);
          this.toast.showError(error);
        });
    }
  }

  public submitAuth(): void {
    if (this.selectedType === AuthFactorType.OTP) {
      this.submitOTP();
    } else if (this.selectedType === AuthFactorType.U2F) {
      this.submitU2F();
    }
  }

  public submitOTP(): void {
    this.authService.verifyMyMultiFactorOTP(this.otpcode).then(
      () => {
        this.dialogRef.close(true);
      },
      (error) => {
        this.toast.showError(error);
        this.dialogRef.close(false);
      },
    );
  }

  public submitU2F(): void {
    if (this.u2fname && this.u2fCredentialOptions.publicKey) {
      // this.data.credOptions.publicKey.rp.id = 'localhost';
      navigator.credentials
        .create(this.u2fCredentialOptions)
        .then((resp) => {
          if (
            resp &&
            (resp as any).response.attestationObject &&
            (resp as any).response.clientDataJSON &&
            (resp as any).rawId
          ) {
            const attestationObject = (resp as any).response.attestationObject;
            const clientDataJSON = (resp as any).response.clientDataJSON;
            const rawId = (resp as any).rawId;

            const data = JSON.stringify({
              id: resp.id,
              rawId: _arrayBufferToBase64(rawId),
              type: resp.type,
              response: {
                attestationObject: _arrayBufferToBase64(attestationObject),
                clientDataJSON: _arrayBufferToBase64(clientDataJSON),
              },
            });

            const base64 = btoa(data);

            this.authService
              .verifyMyMultiFactorU2F(base64, this.u2fname)
              .then(() => {
                this.translate
                  .get('USER.MFA.U2F_SUCCESS')
                  .pipe(take(1))
                  .subscribe((msg) => {
                    this.toast.showInfo(msg);
                  });
                this.dialogRef.close(true);
                this.u2fLoading = false;
              })
              .catch((error) => {
                this.u2fLoading = false;
                this.toast.showError(error);
              });
          } else {
            this.u2fLoading = false;
            this.translate
              .get('USER.MFA.U2F_ERROR')
              .pipe(take(1))
              .subscribe((msg) => {
                this.toast.showInfo(msg);
              });
            this.dialogRef.close(true);
          }
        })
        .catch((error) => {
          this.u2fLoading = false;
          this.u2fError = error;
          this.toast.showInfo(error.message);
        });
    }
  }
}
