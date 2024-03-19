import {Component, Inject, Input} from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import {Observable} from "rxjs";
import {Next} from "../provider-next/provider-next.component";

@Component({
  templateUrl: './provider-next-dialog.component.html',
  styleUrls: ['./provider-next-dialog.component.scss'],
})
export class ProviderNextDialogComponent {
  constructor(
    public dialogRef: MatDialogRef<ProviderNextDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public next$: Observable<Next>,
  ) {}

  public closeDialog(): void {
    this.dialogRef.close(false);
  }
}
