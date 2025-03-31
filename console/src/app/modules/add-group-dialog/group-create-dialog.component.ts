import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { Group } from 'src/app/proto/generated/zitadel/group_pb';

@Component({
  selector: 'cnsl-group-create-dialog',
  templateUrl: './group-create-dialog.component.html',
  styleUrls: ['./group-create-dialog.component.scss'],
})
export class GroupCreateDialogComponent {
  public group?: Group.AsObject;
  public groupIds: string[] = [];
  public roles: string[] = [];

  constructor(
    public dialogRef: MatDialogRef<GroupCreateDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any,
  ) {
  }

  public closeDialog(): void {
    this.dialogRef.close(false);
  }

  public closeDialogWithSuccess(): void {
    this.dialogRef.close({
      groups: this.groupIds
    });
  }

  public selectGroups(group: Group.AsObject[]): void {
    if (group && group.length) {
      this.groupIds = (group as Group.AsObject[]).map((u) => u.id);
    }
  }
}
