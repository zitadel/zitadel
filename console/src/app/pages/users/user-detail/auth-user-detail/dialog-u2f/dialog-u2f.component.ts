import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { TranslateService } from '@ngx-translate/core';
import { take } from 'rxjs/operators';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ToastService } from 'src/app/services/toast.service';

export function _arrayBufferToBase64(buffer: any): string {
    let binary = '';
    const bytes = new Uint8Array(buffer);
    const len = bytes.byteLength;
    for (let i = 0; i < len; i++) {
        binary += String.fromCharCode(bytes[i]);
    }
    return btoa(binary).replace(/\+/g, '-')
        .replace(/\//g, '_')
        .replace(/=/g, '');
}

@Component({
    selector: 'app-dialog-u2f',
    templateUrl: './dialog-u2f.component.html',
    styleUrls: ['./dialog-u2f.component.scss'],
})
export class DialogU2FComponent {
    public name: string = '';
    public error: string = '';
    public loading: boolean = false;

    constructor(public dialogRef: MatDialogRef<DialogU2FComponent>,
        @Inject(MAT_DIALOG_DATA) public data: { credOptions: any; },
        private service: GrpcAuthService, private translate: TranslateService, private toast: ToastService) { }

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
                    this.service.VerifyMyMfaU2F(base64, this.name).then(() => {
                        this.translate.get('USER.MFA.U2F_SUCCESS').pipe(take(1)).subscribe(msg => {
                            this.toast.showInfo(msg);
                        });
                        this.dialogRef.close(true);
                        this.loading = false;
                    }).catch(error => {
                        this.loading = false;
                        this.toast.showError(error);
                    });
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
}
