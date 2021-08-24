import { Component, Inject, Injector, Type } from '@angular/core';
import { FormArray, FormControl, FormGroup, Validators } from '@angular/forms';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { BulkSetUserMetadataRequest } from 'src/app/proto/generated/zitadel/management_pb';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';


@Component({
  selector: 'app-metadata-dialog',
  templateUrl: './metadata-dialog.component.html',
  styleUrls: ['./metadata-dialog.component.scss'],
})
export class MetadataDialogComponent {
  public injData: any = {};
  private service!: GrpcAuthService | ManagementService;
  public loading: boolean = true;
  public ts!: Timestamp.AsObject | undefined;

  public formArray!: FormArray;
  public formGroup!: FormGroup;

  constructor(
    private injector: Injector,
    private toast: ToastService,
    public dialogRef: MatDialogRef<MetadataDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any) {

    this.injData = data;
    switch (this.data.serviceType) {
      case 'MGMT':
        this.service = this.injector.get(ManagementService as Type<ManagementService>);
        break;
      case 'AUTH':
        this.service = this.injector.get(GrpcAuthService as Type<GrpcAuthService>);
        break;
    }

    this.load();

    this.formGroup = new FormGroup({
      key: new FormControl('', [Validators.required]),
      value: new FormControl('', [Validators.required]),
    });

    this.formArray = new FormArray([this.formGroup]);
  }


  public addEntry(): void {
    const newGroup = new FormGroup({
      key: new FormControl('', [Validators.required]),
      value: new FormControl('', [Validators.required]),
    });

    this.formArray.push(newGroup);
  }

  public removeEntry(index: number): void {
    this.formArray.removeAt(index);
  }

  public load(): void {
    this.loadMetadata().then(() => {
      this.loading = false;
      if (this.formArray.length === 0) {
        this.addEntry();
      }
    }).catch(error => {
      this.loading = false;
      this.toast.showError(error);
      if (this.formArray.length === 0) {
        this.addEntry();
      }
    });
  }

  public loadMetadata(): Promise<any> {
    this.loading = true;
    if (this.data.serviceType === 'MGMT' && this.injData.userId) {
      return (this.service as ManagementService).listUserMetadata(this.injData.userId).then(resp => {
        this.formArray.patchValue(resp.resultList);
        this.ts = resp.details?.viewTimestamp;
      });
    } else {
      return (this.service as GrpcAuthService).listMyMetadata().then(resp => {
        this.formArray.patchValue(resp.resultList);
        this.ts = resp.details?.viewTimestamp;
      });
    }
  }

  public setMetadataAndClose(): void {
    this.loading = true;
    const metadataList = this.formArray.value;

    switch (this.injData.serviceType) {
      case 'MGMT':
        const bulk = metadataList.map((element: any) => {
          const e = new BulkSetUserMetadataRequest.Metadata();
          e.setKey(element.key);
          e.setValue(element.value);
          return e;
        });

        (this.service as ManagementService).bulkSetUserMetadata(metadataList, this.injData.userId)
          .then(() => {
            this.toast.showInfo('USER.METADATA.SETSUCCESS', true);
            this.loading = false;
            this.dialogRef.close();
          }).catch(error => {
            this.loading = false;
            this.toast.showError(error);
          });
        break;
      case 'AUTH':
        const mybulk = metadataList.map((element: any) => {
          const e = new BulkSetUserMetadataRequest.Metadata();
          e.setKey(element.key);
          e.setValue(element.value);
          return e;
        });

        (this.service as GrpcAuthService).bulkSetMyMetadata(mybulk)
          .then(() => {
            this.toast.showInfo('USER.METADATA.SETSUCCESS', true);
            this.loading = false;
            this.dialogRef.close();
          }).catch(error => {
            this.loading = false;
            this.toast.showError(error);
          });
        break;
    }
  }

  public removeMetadata(key: string): Promise<any> {
    switch (this.injData.serviceType) {
      case 'MGMT': return (this.service as ManagementService).removeUserMetadata(key, this.injData.userId)
        .then(() => {
          this.toast.showInfo('USER.METADATA.REMOVESUCCESS', true);
        }).catch(error => {
          this.toast.showError(error);
        });
      case 'AUTH': return (this.service as GrpcAuthService).removeMyMetadata(key)
        .then(() => {
          this.toast.showInfo('USER.METADATA.REMOVESUCCESS', true);
        }).catch(error => {
          this.toast.showError(error);
        });
      default:
        return Promise.reject();
    }
  }

  closeDialog(): void {
    this.dialogRef.close();
  }
}
