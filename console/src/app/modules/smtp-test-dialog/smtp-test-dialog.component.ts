import { Component, Inject, signal } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';

import { InfoSectionType } from '../info-section/info-section.component';
import { AdminService } from 'src/app/services/admin.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ToastService } from 'src/app/services/toast.service';
import { TranslateService } from '@ngx-translate/core';
import { TestSMTPConfigByIdRequest } from 'src/app/proto/generated/zitadel/admin_pb';
import { take } from 'rxjs';

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

    const req = new TestSMTPConfigByIdRequest();
    req.setId(this.data.id);
    req.setReceiverAddress(this.email);

    this.adminService
      .testSMTPConfigById(req)
      .then(() => {
        this.resultClass = 'test-success';
        this.isLoading.set(false);
        this.translate
          .get('SMTP.CREATE.STEPS.TEST.RESULT')
          .pipe(take(1))
          .subscribe((msg) => {
            this.testResult = msg;
          });
      })
      .catch((error) => {
        this.resultClass = 'test-error';
        this.isLoading.set(false);
        this.testResult = error;
      });
  }

  public closeDialog(): void {
    this.dialogRef.close(false);
  }
}
