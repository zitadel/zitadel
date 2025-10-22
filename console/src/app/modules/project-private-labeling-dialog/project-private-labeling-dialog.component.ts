import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { PrivateLabelingSetting } from 'src/app/proto/generated/zitadel/project_pb';

@Component({
  selector: 'cnsl-project-private-labeling-dialog',
  templateUrl: './project-private-labeling-dialog.component.html',
  styleUrls: ['./project-private-labeling-dialog.component.scss'],
  standalone: false,
})
export class ProjectPrivateLabelingDialogComponent {
  public setting: PrivateLabelingSetting = PrivateLabelingSetting.PRIVATE_LABELING_SETTING_UNSPECIFIED;
  public settings: PrivateLabelingSetting[] = [
    PrivateLabelingSetting.PRIVATE_LABELING_SETTING_UNSPECIFIED,
    PrivateLabelingSetting.PRIVATE_LABELING_SETTING_ENFORCE_PROJECT_RESOURCE_OWNER_POLICY,
    PrivateLabelingSetting.PRIVATE_LABELING_SETTING_ALLOW_LOGIN_USER_RESOURCE_OWNER_POLICY,
  ];
  constructor(
    public dialogRef: MatDialogRef<ProjectPrivateLabelingDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any,
  ) {
    this.setting = data.setting;
  }

  closeDialog(setting?: PrivateLabelingSetting): void {
    this.dialogRef.close(setting);
  }
}
