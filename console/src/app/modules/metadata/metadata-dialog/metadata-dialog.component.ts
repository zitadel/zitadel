import { Component, Inject } from '@angular/core';
import {
  MatLegacyDialogRef as MatDialogRef,
  MAT_LEGACY_DIALOG_DATA as MAT_DIALOG_DATA,
} from '@angular/material/legacy-dialog';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { Metadata } from 'src/app/proto/generated/zitadel/metadata_pb';
import { ToastService } from 'src/app/services/toast.service';

@Component({
  selector: 'cnsl-metadata-dialog',
  templateUrl: './metadata-dialog.component.html',
  styleUrls: ['./metadata-dialog.component.scss'],
})
export class MetadataDialogComponent {
  public metadata: Partial<Metadata.AsObject>[] = [];
  public loading: boolean = false;
  public ts!: Timestamp.AsObject | undefined;

  constructor(
    private toast: ToastService,
    public dialogRef: MatDialogRef<MetadataDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any,
  ) {
    this.metadata = data.metadata;
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
      this.data
        .setFcn(key, value)
        .then(() => {
          this.toast.showInfo('METADATA.SETSUCCESS', true);
        })
        .catch((error: any) => {
          this.toast.showError(error);
        });
    }
  }

  public removeMetadata(key: string): Promise<void> {
    return this.data
      .removeFcn(key)
      .then((resp: any) => {
        this.toast.showInfo('METADATA.REMOVESUCCESS', true);
      })
      .catch((error: any) => {
        this.toast.showError(error);
      });
  }

  closeDialog(): void {
    this.dialogRef.close();
  }
}
