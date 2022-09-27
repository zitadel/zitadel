import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { Buffer } from 'buffer';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { Metadata } from 'src/app/proto/generated/zitadel/metadata_pb';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
  selector: 'cnsl-metadata-dialog',
  templateUrl: './metadata-dialog.component.html',
  styleUrls: ['./metadata-dialog.component.scss'],
})
export class MetadataDialogComponent {
  public metadata: Partial<Metadata.AsObject>[] = [];
  public injData: any = {};
  public loading: boolean = true;
  public ts!: Timestamp.AsObject | undefined;

  constructor(
    private managementService: ManagementService,
    private authService: GrpcAuthService,
    private toast: ToastService,
    public dialogRef: MatDialogRef<MetadataDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any,
  ) {
    this.injData = data;
    this.load();
  }

  public load(): void {
    this.loadMetadata()
      .then(() => {
        this.loading = false;
        if (this.metadata.length === 0) {
          this.addEntry();
        }
      })
      .catch((error) => {
        this.loading = false;
        this.toast.showError(error);
        if (this.metadata.length === 0) {
          this.addEntry();
        }
      });
  }

  public loadMetadata(): Promise<void> {
    this.loading = true;
    if (this.injData.userId) {
      return this.managementService.listUserMetadata(this.injData.userId).then((resp) => {
        this.metadata = resp.resultList.map((md) => {
          return {
            key: md.key,
            value: Buffer.from(md.value as string, 'base64'),
          };
        });
        this.ts = resp.details?.viewTimestamp;
      });
    } else {
      return this.authService.listMyMetadata().then((resp) => {
        this.metadata = resp.resultList.map((md) => {
          return {
            key: md.key,
            value: Buffer.from(md.value as string, 'base64'),
          };
        });
        this.ts = resp.details?.viewTimestamp;
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
    if (key && value) {
      this.managementService
        .setUserMetadata(key, btoa(value), this.injData.userId)
        .then(() => {
          this.toast.showInfo('USER.METADATA.SETSUCCESS', true);
        })
        .catch((error) => {
          this.toast.showError(error);
        });
    }
  }

  public removeMetadata(key: string): Promise<void> {
    return this.managementService
      .removeUserMetadata(key, this.injData.userId)
      .then((resp) => {
        this.toast.showInfo('USER.METADATA.REMOVESUCCESS', true);
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  closeDialog(): void {
    this.dialogRef.close();
  }
}
