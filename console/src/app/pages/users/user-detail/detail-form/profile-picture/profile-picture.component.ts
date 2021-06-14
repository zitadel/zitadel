import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { DomSanitizer } from '@angular/platform-browser';
import { AssetService } from 'src/app/services/asset.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
  selector: 'cnsl-profile-picture',
  templateUrl: './profile-picture.component.html',
  styleUrls: ['./profile-picture.component.scss'],
})
export class ProfilePictureComponent {
  constructor(
    private authService: GrpcAuthService,
    public dialogRef: MatDialogRef<ProfilePictureComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any,
    private toast: ToastService,
    private assetService: AssetService,
    private sanitizer: DomSanitizer,
  ) {
  }

  public onDrop(event: any): Promise<any> | void {
    const filelist: FileList = event.target.files;
    const file = filelist.item(0);

    if (file) {
      const formData = new FormData();
      formData.append('file', file);
      return this.handleUploadPromise(this.assetService.upload('users/me/avatar', formData));
    }
  }

  public deletePic(): void {
    this.authService.removeMyAvatar().then(() => {
      this.toast.showInfo('USER.PROFILE.AVATAR.DELETESUCCESS', true);
      this.data.profilePic = null;
    }).catch(error => {
      this.toast.showError(error);
    });
  }

  private handleUploadPromise(task: Promise<any>): Promise<any> {
    return task.then(() => {
      this.toast.showInfo('POLICY.TOAST.UPLOADSUCCESS', true);
      this.assetService.load('users/me/avatar').then(data => {
        const objectURL = URL.createObjectURL(data);
        const pic = this.sanitizer.bypassSecurityTrustUrl(objectURL);
        this.data.profilePic = pic;
      }).catch(error => {
        this.toast.showError(error);
      });
    }).catch(error => this.toast.showError(error));
  }

  public closeDialog(): void {
    this.dialogRef.close(false);
  }
}
