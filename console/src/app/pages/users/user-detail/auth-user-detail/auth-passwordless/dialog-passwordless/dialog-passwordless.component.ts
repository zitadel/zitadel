import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { TranslateService } from '@ngx-translate/core';
import { take } from 'rxjs/operators';

import { EnvironmentService } from '../../../../../../services/environment.service';
import { ToastService } from '../../../../../../services/toast.service';
import { UserService } from '../../../../../../services/user.service';
import { _arrayBufferToBase64 } from '../../u2f_util';

@Component({
  selector: 'cnsl-dialog-passwordless',
  templateUrl: './dialog-passwordless.component.html',
  styleUrls: ['./dialog-passwordless.component.scss'],
  standalone: false,
})
export class DialogPasswordlessComponent {
  public name: string = '';
  public error: string = '';
  public loading: boolean = false;

  public showSent: boolean = false;
  public showQR: boolean = false;
  public qrcodeLink: string = '';

  private loginBaseUrl: string = window.location.origin;

  constructor(
    public dialogRef: MatDialogRef<DialogPasswordlessComponent>,
    @Inject(MAT_DIALOG_DATA) public data: { credOptions: any; passkeyId: string },
    private userService: UserService,
    private envService: EnvironmentService,
    private translate: TranslateService,
    private toast: ToastService,
  ) {
    this.envService.env.subscribe((env) => {
      if (env.login_v2_base_url) {
        this.loginBaseUrl = env.login_v2_base_url;
      }
    });
  }

  public closeDialog(): void {
    this.dialogRef.close();
  }

  public closeDialogWithCode(): void {
    this.error = '';
    this.loading = true;

    const userId = this.userService.userId();
    if (!userId) {
      this.loading = false;
      this.toast.showError('USER.PASSWORDLESS.U2F_ERROR', false, true);
      return;
    }

    if (this.name && this.data.credOptions.publicKey) {
      navigator.credentials
        .create(this.data.credOptions)
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

            // v2 expects the credential as a structured object (google.protobuf.Struct), not a
            // base64-encoded JSON string like the old v1 API.
            const publicKeyCredential = {
              id: resp.id,
              rawId: _arrayBufferToBase64(rawId),
              type: resp.type,
              response: {
                attestationObject: _arrayBufferToBase64(attestationObject),
                clientDataJSON: _arrayBufferToBase64(clientDataJSON),
              },
            };

            this.userService
              .verifyPasskeyRegistration({
                userId,
                passkeyId: this.data.passkeyId,
                publicKeyCredential,
                passkeyName: this.name,
              })
              .then(() => {
                this.translate
                  .get('USER.PASSWORDLESS.U2F_SUCCESS')
                  .pipe(take(1))
                  .subscribe((msg) => {
                    this.toast.showInfo(msg);
                  });
                this.dialogRef.close(true);
                this.loading = false;
              })
              .catch((error) => {
                this.loading = false;
                this.toast.showError(error);
              });
          } else {
            this.loading = false;
            this.translate
              .get('USER.PASSWORDLESS.U2F_ERROR')
              .pipe(take(1))
              .subscribe((msg) => {
                this.toast.showInfo(msg);
              });
            this.dialogRef.close(true);
          }
        })
        .catch((error) => {
          this.loading = false;
          this.error = error;
          this.toast.showInfo(error.message);
        });
    }
  }

  public sendMyPasswordlessLink(): void {
    const userId = this.userService.userId();
    if (!userId) {
      this.toast.showError('USER.PASSWORDLESS.U2F_ERROR', false, true);
      return;
    }

    this.userService
      .createPasskeyRegistrationLink({
        userId,
        medium: { case: 'sendLink', value: {} },
      })
      .then(() => {
        this.toast.showInfo('USER.TOAST.PASSWORDLESSREGISTRATIONSENT', true);
        this.showSent = true;
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public addMyPasswordlessLink(): void {
    const userId = this.userService.userId();
    if (!userId) {
      this.toast.showError('USER.PASSWORDLESS.U2F_ERROR', false, true);
      return;
    }

    this.userService
      .createPasskeyRegistrationLink({
        userId,
        medium: { case: 'returnCode', value: {} },
      })
      .then((resp) => {
        if (!resp.code) {
          this.toast.showError('USER.PASSWORDLESS.U2F_ERROR', false, true);
          return;
        }

        const params = new URLSearchParams({
          userId,
          code: resp.code.code,
          codeId: resp.code.id,
        });

        this.showQR = true;
        this.qrcodeLink = `${this.loginBaseUrl}/passkey/set?${params.toString()}`;
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }
}
