import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { TranslateService } from '@ngx-translate/core';
import { take } from 'rxjs/operators';

import { GrpcAuthService } from '../../../../../../services/grpc-auth.service';
import { ToastService } from '../../../../../../services/toast.service';
import { _arrayBufferToBase64 } from '../../u2f_util';

export enum U2FComponentDestination {
  MFA = 'mfa',
  PASSWORDLESS = 'passwordless',
}

@Component({
  selector: 'app-dialog-passwordless',
  templateUrl: './dialog-passwordless.component.html',
  styleUrls: ['./dialog-passwordless.component.scss'],
})
export class DialogPasswordlessComponent {
  private type!: U2FComponentDestination;
  public name: string = '';
  public error: string = '';
  public loading: boolean = false;

  public showSent: boolean = false;
  public showQR: boolean = false;
  public qrcodeLink: string = '';

  constructor(public dialogRef: MatDialogRef<DialogPasswordlessComponent>,
    @Inject(MAT_DIALOG_DATA) public data: { credOptions: any; type: U2FComponentDestination; },
    private service: GrpcAuthService, private translate: TranslateService, private toast: ToastService) {
    this.type = data.type;
  }

  public closeDialog(): void {
    this.dialogRef.close();
  }

  public closeDialogWithCode(): void {
    this.error = '';
    this.loading = true;
    if (this.name && this.data.credOptions.publicKey) {
      // this.data.credOptions.publicKey.rp.id = 'localhost';
      navigator.credentials.create(this.data.credOptions).then((resp) => {
        if (resp &&
          (resp as any).response.attestationObject &&
          (resp as any).response.clientDataJSON &&
          (resp as any).rawId) {

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
          if (this.type === U2FComponentDestination.MFA) {
            this.service.verifyMyMultiFactorU2F(base64, this.name).then(() => {
              this.translate.get('USER.MFA.U2F_SUCCESS').pipe(take(1)).subscribe(msg => {
                this.toast.showInfo(msg);
              });
              this.dialogRef.close(true);
              this.loading = false;
            }).catch(error => {
              this.loading = false;
              this.toast.showError(error);
            });
          } else if (this.type === U2FComponentDestination.PASSWORDLESS) {
            this.service.verifyMyPasswordless(base64, this.name).then(() => {
              this.translate.get('USER.PASSWORDLESS.U2F_SUCCESS').pipe(take(1)).subscribe(msg => {
                this.toast.showInfo(msg);
              });
              this.dialogRef.close(true);
              this.loading = false;
            }).catch(error => {
              this.loading = false;
              this.toast.showError(error);
            });
          }
        } else {
          this.loading = false;
          this.translate.get('USER.MFA.U2F_ERROR').pipe(take(1)).subscribe(msg => {
            this.toast.showInfo(msg);
          });
          this.dialogRef.close(true);
        }
      }).catch(error => {
        this.loading = false;
        this.error = error;
        this.toast.showInfo(error.message);
      });
    }

  }

  public sendMyPasswordlessLink(): void {
    this.service.sendMyPasswordlessLink().then(() => {
      this.toast.showInfo('USER.TOAST.PASSWORDLESSREGISTRATIONSENT');
      this.showSent = true;
    }).catch(error => {
      this.toast.showError(error);
    });
  }

  public addMyPasswordlessLink(): void {
    this.service.addMyPasswordlessLink().then((resp) => {
      console.log(resp);
      this.showQR = true;

      this.qrcodeLink = resp.link;

    }).catch(error => {
      this.toast.showError(error);
    });
  }
}
