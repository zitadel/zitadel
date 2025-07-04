import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { ToastService } from 'src/app/services/toast.service';
import { Metadata as MetadataV2 } from '@zitadel/proto/zitadel/metadata_pb';
import { Metadata } from 'src/app/proto/generated/zitadel/metadata_pb';

export type MetadataDialogData = {
  metadata: (Metadata.AsObject | MetadataV2)[];
  setFcn: (key: string, value: string) => Promise<any>;
  removeFcn: (key: string) => Promise<any>;
};

@Component({
  selector: 'cnsl-metadata-dialog',
  templateUrl: './metadata-dialog.component.html',
  styleUrls: ['./metadata-dialog.component.scss'],
})
export class MetadataDialogComponent {
  public metadata: { key: string; value: string }[] = [];
  public ts!: Timestamp.AsObject | undefined;

  constructor(
    private toast: ToastService,
    public dialogRef: MatDialogRef<MetadataDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: MetadataDialogData,
  ) {
    const decoder = new TextDecoder();
    this.metadata = data.metadata.map(({ key, value }) => ({
      key,
      value: typeof value === 'string' ? value : decoder.decode(value),
    }));
  }

  public addEntry(): void {
    this.metadata.push({
      key: '',
      value: '',
    });
  }

  public async removeEntry(index: number) {
    const key = this.metadata[index].key;
    if (!key) {
      this.metadata.splice(index, 1);
      return;
    }

    try {
      await this.data.removeFcn(key);
    } catch (error) {
      this.toast.showError(error);
      return;
    }

    this.toast.showInfo('METADATA.REMOVESUCCESS', true);
    this.metadata.splice(index, 1);
    if (this.metadata.length === 0) {
      this.addEntry();
    }
  }

  public async saveElement(index: number) {
    const { key, value } = this.metadata[index];

    if (!key || !value) {
      return;
    }

    try {
      await this.data.setFcn(key, value);
      this.toast.showInfo('METADATA.SETSUCCESS', true);
    } catch (error) {
      this.toast.showError(error);
    }
  }

  closeDialog(): void {
    this.dialogRef.close();
  }
}
