import { Component, Inject, OnInit } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { ToastService } from 'src/app/services/toast.service';
import { UploadService } from 'src/app/services/upload.service';

@Component({
  selector: 'cnsl-profile-picture',
  templateUrl: './profile-picture.component.html',
  styleUrls: ['./profile-picture.component.scss']
})
export class ProfilePictureComponent implements OnInit {
  public isHovering: boolean = false;

  constructor(
    public dialogRef: MatDialogRef<ProfilePictureComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any,
    private toast: ToastService,
    private uploadService: UploadService) { }

  ngOnInit(): void {
  }

  toggleHover(isHovering: boolean) {
    this.isHovering = isHovering;
  }

  public onDrop(filelist: FileList): Promise<any> | void {
    const file = filelist.item(0);
    if (file) {

      const formData = new FormData();
      formData.append('file', file);
      // switch (this.serviceType) {
      //   case PolicyComponentServiceType.MGMT:
      return this.handleUploadPromise(this.uploadService.upload('', formData));
      // case PolicyComponentServiceType.ADMIN:
      //   return this.handleUploadPromise(this.uploadService.upload(UploadEndpoint.IAMDARKLOGO, formData));
    }

  }


  private handleUploadPromise(task: Promise<any>): Promise<any> {
    return task.then(() => {
      this.toast.showInfo('POLICY.TOAST.UPLOADSUCCESS', true);
    }).catch(error => this.toast.showError(error));
  }

  public closeDialog(): void {
    this.dialogRef.close(false);
  }
}
