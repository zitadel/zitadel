import { Component, Inject, Injector, Type } from '@angular/core';
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
  public metadata: Partial<Metadata.AsObject>[] = [];
  public injData: any = {};
  private service!: GrpcAuthService | ManagementService;
  public loading: boolean = true;
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

    this.loadMetadata(data.userId).then(() => {
      this.loading = false;
      if (this.metadata.length === 0) {
        this.addEntry();
      }
    }).catch(error => {
      this.loading = false;
      this.toast.showError(error);
      if (this.metadata.length === 0) {
        this.addEntry();
      }
    });
  }

  public loadMetadata(userId?: string): Promise<any> {
    if (this.data.serviceType === 'MGMT' && userId) {
      return (this.service as ManagementService).listUserMetadata(userId).then(resp => {
        this.metadata = resp.resultList;
      });
    } else {
      return (this.service as GrpcAuthService).listMyMetadata().then(resp => {
        this.metadata = resp.resultList;
      });
    }
  }

  public addEntry(): void {
    const newGroup = {
      key: '',
      value: '',
    };

    this.metadata.push(newGroup);
  }

  public removeEntry(index: number): void {
    const key = this.metadata[index].key;
    if (key) {
      this.removeMetadata(key).then(() => {
        this.metadata.splice(index, 1);
        if (this.metadata.length === 0) {
          this.addEntry();
        }
      });
    } else {
      this.metadata.splice(index, 1);
    }
  }

  public saveElement(index: number): void {
    const metadataElement = this.metadata[index];

    if (metadataElement.key && metadataElement.value) {
      this.setMetadata(metadataElement.key, metadataElement.value as string);
    }
  }

  public setMetadata(key: string, value: string): void {
    console.log(key, value, this.injData.userId);
    if (key && value) {
      switch (this.injData.serviceType) {
        case 'MGMT': (this.service as ManagementService).setUserMetadata(key, value, this.injData.userId)
          .then(() => {
            this.toast.showInfo('USER.METADATA.SETSUCCESS', true);
          }).catch(error => {
            this.toast.showError(error);
          });
          break;
        case 'AUTH': (this.service as GrpcAuthService).setMyMetadata(key, value)
          .then(() => {
            this.toast.showInfo('USER.METADATA.SETSUCCESS', true);
          }).catch(error => {
            this.toast.showError(error);
          });
          break;
      }
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
