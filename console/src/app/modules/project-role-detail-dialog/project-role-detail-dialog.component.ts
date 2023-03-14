import { Component, Inject } from '@angular/core';
import { AbstractControl, UntypedFormControl, UntypedFormGroup } from '@angular/forms';
import {
  MatLegacyDialogRef as MatDialogRef,
  MAT_LEGACY_DIALOG_DATA as MAT_DIALOG_DATA,
} from '@angular/material/legacy-dialog';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';
import { requiredValidator } from '../form-field/validators/validators';

@Component({
  selector: 'cnsl-project-role-detail-dialog',
  templateUrl: './project-role-detail-dialog.component.html',
  styleUrls: ['./project-role-detail-dialog.component.scss'],
})
export class ProjectRoleDetailDialogComponent {
  public projectId: string = '';

  public formGroup!: UntypedFormGroup;
  constructor(
    private mgmtService: ManagementService,
    private toast: ToastService,
    public dialogRef: MatDialogRef<ProjectRoleDetailDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any,
  ) {
    this.projectId = data.projectId;
    this.formGroup = new UntypedFormGroup({
      key: new UntypedFormControl({ value: '', disabled: true }, [requiredValidator]),
      displayName: new UntypedFormControl(''),
      group: new UntypedFormControl(''),
    });

    this.formGroup.patchValue(data.role);
  }

  submitForm(): void {
    if (this.formGroup.valid && this.key?.value && this.group?.value && this.displayName?.value) {
      this.mgmtService
        .updateProjectRole(this.projectId, this.key.value, this.displayName.value, this.group.value)
        .then(() => {
          this.toast.showInfo('PROJECT.TOAST.ROLECHANGED', true);
          this.dialogRef.close(true);
        })
        .catch((error) => {
          this.toast.showError(error);
        });
    }
  }

  public closeDialog(): void {
    this.dialogRef.close(false);
  }

  public get key(): AbstractControl | null {
    return this.formGroup.get('key');
  }
  public get displayName(): AbstractControl | null {
    return this.formGroup.get('displayName');
  }
  public get group(): AbstractControl | null {
    return this.formGroup.get('group');
  }
}
