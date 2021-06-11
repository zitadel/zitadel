import { Component, Inject, OnInit } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { AssetService } from 'src/app/services/asset.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
  selector: 'cnsl-profile-picture',
  templateUrl: './profile-picture.component.html',
  styleUrls: ['./profile-picture.component.scss'],
})
export class ProfilePictureComponent implements OnInit {
  public isHovering: boolean = false;

  public selectedFile: any = null;
  public imageChangedEvent: any = '';
  // public imageChangedFormat: string = '';
  public croppedImage: any = '';
  public showCropperError: boolean = false;

  constructor(
    private authService: GrpcAuthService,
    public dialogRef: MatDialogRef<ProfilePictureComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any,
    private toast: ToastService,
    private assetService: AssetService) { }

  public ngOnInit(): void {
  }

  public toggleHover(isHovering: boolean): void {
    this.isHovering = isHovering;
  }

  public onDrop(event: any): Promise<any> | void {
    const filelist: FileList = event.target.files;
    this.imageChangedEvent = event;
    const file = filelist.item(0);
    const reader = new FileReader();


    if (file) {
      this.selectedFile = file;

      reader.readAsDataURL(file); //FileStream response from .NET core backend
      reader.onload = _event => {
        this.croppedImage = reader.result;
      };
      // this.imageChangedFormat = file.type;
    }
  }

  public deletePic(): void {
    console.log('delete');
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
      this.data.profilePic = this.croppedImage;
    }).catch(error => this.toast.showError(error));
  }

  // public fileChangeEvent(event: any): void {
  //   this.imageChangedEvent = event;
  // }

  // public imageCropped(event: ImageCroppedEvent): void {
  //   this.showCropperError = false;
  //   this.croppedImage = event.base64;
  // }

  public upload(): Promise<any> | void {
    // const formData = new FormData();
    // const splitted = this.croppedImage.split(';base64,');
    // if (splitted[1]) {
    //   const blob = this.base64toBlob(splitted[1]);
    //   formData.append('file', blob);
    //   return this.handleUploadPromise(this.assetService.upload('users/me/avatar', formData));
    // }

    const formData = new FormData();
    formData.append('file', this.selectedFile);
    return this.handleUploadPromise(this.assetService.upload('users/me/avatar', formData));
  }

  // public loadImageFailed(): void {
  //   this.showCropperError = true;
  // }

  public closeDialog(): void {
    this.dialogRef.close(false);
  }

  // public base64toBlob(b64: string): Blob {
  //   const byteCharacters = atob(b64);
  //   const byteNumbers = new Array(byteCharacters.length);
  //   for (let i = 0; i < byteCharacters.length; i++) {
  //     byteNumbers[i] = byteCharacters.charCodeAt(i);
  //   }
  //   const byteArray = new Uint8Array(byteNumbers);
  //   const blob = new Blob([byteArray], { type: this.imageChangedFormat });
  //   return blob;
  // }
}
