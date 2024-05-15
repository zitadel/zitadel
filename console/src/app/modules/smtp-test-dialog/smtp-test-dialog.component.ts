import { Component, Inject, signal } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';

import { InfoSectionType } from '../info-section/info-section.component';
import { AdminService } from 'src/app/services/admin.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ToastService } from 'src/app/services/toast.service';
import { TranslateService } from '@ngx-translate/core';

@Component({
  selector: 'cnsl-smtp-test-dialog',
  templateUrl: './smtp-test-dialog.component.html',
  styleUrls: ['./smtp-test-dialog.component.scss'],
})
export class SmtpTestDialogComponent {
  public resultClass = 'test-success';
  public isLoading = signal(false);
  public email: string = '';
  public testResult: string = '';
  InfoSectionType: any = InfoSectionType;
  constructor(
    public dialogRef: MatDialogRef<SmtpTestDialogComponent>,
    private adminService: AdminService,
    private authService: GrpcAuthService,
    private toast: ToastService,
    private translate: TranslateService,
    @Inject(MAT_DIALOG_DATA) public data: any,
  ) {
    this.authService
      .getMyUser()
      .then((resp) => {
        if (resp.user) {
          this.email = resp.user.human?.email?.email || '';
        }
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public testEmailConfiguration(): void {
    this.isLoading.set(true);
    this.adminService
      .testSMTPConfigById(this.data.id, this.email)
      .then(() => {
        this.resultClass = 'test-success';
        this.isLoading.set(false);
        // this.translate.get('ORG.TOAST.ORG_WAS_DELETED').subscribe((data) => {
        //   this.toast.showInfo(data);
        // });

        this.testResult = 'Your email was succesfully sent';
      })
      .catch((error) => {
        this.resultClass = 'test-error';
        this.isLoading.set(false);
        let errorMsg: string = error.message;
        if (errorMsg.includes('could not make smtp dial')) {
          errorMsg =
            "We couldn't contact with the SMTP server, check the server port, if you're behind a proxy or firewall... ";
          // TODO translate error message
        }

        if (errorMsg.includes('could not add smtp auth')) {
          errorMsg =
            "There was an issue with authentication, check that your user and password are correct. If they're correct maybe your provider requires an auth method not supported by Zitadel";
          // TODO translate error message
        }

        this.testResult = errorMsg;
      });
  }

  public closeDialog(): void {
    this.dialogRef.close(false);
  }
}
