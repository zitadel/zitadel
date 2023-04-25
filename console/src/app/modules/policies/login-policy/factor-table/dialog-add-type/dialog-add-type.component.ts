import { Component, Inject } from '@angular/core';
import {
  MatLegacyDialogRef as MatDialogRef,
  MAT_LEGACY_DIALOG_DATA as MAT_DIALOG_DATA,
} from '@angular/material/legacy-dialog';
import { MultiFactorType, SecondFactorType } from 'src/app/proto/generated/zitadel/policy_pb';

enum LoginMethodComponentType {
  MultiFactor = 1,
  SecondFactor = 2,
}

@Component({
  selector: 'cnsl-dialog-add-type',
  templateUrl: './dialog-add-type.component.html',
  styleUrls: ['./dialog-add-type.component.scss'],
})
export class DialogAddTypeComponent {
  public LoginMethodComponentType: any = LoginMethodComponentType;
  public availableMfaTypes: Array<MultiFactorType | SecondFactorType> = [];
  public newMfaType!: MultiFactorType | SecondFactorType;

  constructor(public dialogRef: MatDialogRef<DialogAddTypeComponent>, @Inject(MAT_DIALOG_DATA) public data: any) {
    this.availableMfaTypes = data.types;
    this.newMfaType = data.types && data.types[0] ? data.types[0] : undefined;
  }

  public closeDialog(): void {
    this.dialogRef.close();
  }

  public closeDialogWithCode(): void {
    this.dialogRef.close(this.newMfaType);
  }
}
