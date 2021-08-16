import { Component, Inject, Injector, Type } from '@angular/core';
import { FormArray, FormControl, FormGroup, Validators } from '@angular/forms';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { Metadata } from 'src/app/proto/generated/zitadel/metadata_pb';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';


@Component({
  selector: 'app-metadata-dialog',
  templateUrl: './metadata-dialog.component.html',
  styleUrls: ['./metadata-dialog.component.scss'],
})
export class MetadataDialogComponent {
  public metadata: Metadata.AsObject[] = [];

  public formGroup!: FormGroup;
  public formArray!: FormArray;

  public injData: any = {};
  private service!: GrpcAuthService | ManagementService;

  constructor(
    private injector: Injector,
    private toast: ToastService,
    public dialogRef: MatDialogRef<MetadataDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any) {
    this.formGroup = new FormGroup({
      key: new FormControl('', [Validators.required]),
      value: new FormControl('', [Validators.required]),
    });
    this.formArray = new FormArray([this.formGroup]);

    this.injData = data;
    switch (this.data.serviceType) {
      case 'MGMT':
        this.service = this.injector.get(ManagementService as Type<ManagementService>);
        break;
      case 'AUTH':
        this.service = this.injector.get(GrpcAuthService as Type<GrpcAuthService>);
        break;
    }

    this.loadMetadata(data.userId);
  }

  public loadMetadata(userId?: string): void {
    if (this.data.serviceType === 'MGMT' && userId) {
      (this.service as ManagementService).listUserMetadata(userId).then(resp => {
        this.metadata = resp.resultList;
        this.formArray.patchValue(this.metadata);
      });
    } else if (this.data.serviceType === 'AUTH') {
      (this.service as GrpcAuthService).listMyMetadata().then(resp => {
        this.metadata = resp.resultList;
        this.formArray.patchValue(this.metadata);
      });
    }
  }

  public addEntry(): void {
    const newGroup = new FormGroup({
      key: new FormControl('', [Validators.required]),
      value: new FormControl('', [Validators.required]),
    });

    this.formArray.push(newGroup);
  }

  public removeEntry(index: number): void {
    const key = this.formArray.controls[index].get('key')?.value;
    if (key) {

    } else {
      this.formArray.removeAt(index);
    }
  }

  public saveElement(index: number): void {
    const formControl = this.formArray.controls[index];

    if (formControl.valid) {
      this.setMetadata(formControl.get('key')?.value, formControl.get('value')?.value);
    }
  }

  public setMetadata(key: string, value: string): void {
    if (key && value) {
      switch (this.injData.serviceType) {
        case 'MGMT': (this.service as ManagementService).setUserMetadata(key, value, this.injData.userId)
          .then(() => {
            this.toast.showInfo('');
            this.formGroup.reset();
          }).catch(error => {
            this.toast.showError(error);
          });
          break;
        case 'AUTH': (this.service as GrpcAuthService).setMyMetadata(key, value)
          .then(() => {
            this.toast.showInfo('');
            this.formGroup.reset();
          }).catch(error => {
            this.toast.showError(error);
          });
          break;
      }
    }
  }

  public removeMetadata(key: string): void {
    if (key) {
      switch (this.injData.serviceType) {
        case 'MGMT': (this.service as ManagementService).removeUserMetadata(key, this.injData.userId)
          .then(() => {
            this.toast.showInfo('');
            this.formGroup.reset();
          }).catch(error => {
            this.toast.showError(error);
          });
          break;
        case 'AUTH': (this.service as GrpcAuthService).removeMyMetadata(key)
          .then(() => {
            this.toast.showInfo('');
            this.formGroup.reset();
          }).catch(error => {
            this.toast.showError(error);
          });
          break;
      }
    }
  }

  closeDialog(): void {
    this.dialogRef.close();
  }
}
