import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';

import { InfoSectionType } from '../info-section/info-section.component';
import { AdminService } from 'src/app/services/admin.service';

@Component({
  selector: 'cnsl-smtp-test-dialog',
  templateUrl: './smtp-test-dialog.component.html',
  styleUrls: ['./smtp-test-dialog.component.scss'],
})
export class SmtpTestDialogComponent {
  public email: string = '';
  public testResult: string = '';
  InfoSectionType: any = InfoSectionType;
  constructor(
    public dialogRef: MatDialogRef<SmtpTestDialogComponent>,
    private adminService: AdminService,
    @Inject(MAT_DIALOG_DATA) public data: any,
  ) {}

  public testEmailConfiguration(): void {
    this.adminService
      .testSMTPConfigById(this.data.id, this.email)
      .then(() => {
        this.testResult = 'Your email was succesfully sent';
      })
      .catch((error) => {
        this.testResult = error.message;
      });
  }

  public closeDialog(): void {
    this.dialogRef.close(false);
  }
}
