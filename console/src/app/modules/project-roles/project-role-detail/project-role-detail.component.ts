import { Component, Inject } from '@angular/core';
import { AbstractControl, FormControl, FormGroup, Validators } from '@angular/forms';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
    selector: 'app-project-role-detail',
    templateUrl: './project-role-detail.component.html',
    styleUrls: ['./project-role-detail.component.scss'],
})
export class ProjectRoleDetailComponent {
    public projectId: string = '';

    public formGroup!: FormGroup;
    constructor(private mgmtService: ManagementService, private toast: ToastService,
        public dialogRef: MatDialogRef<ProjectRoleDetailComponent>,
        @Inject(MAT_DIALOG_DATA) public data: any) {

        this.projectId = data.projectId;
        this.formGroup = new FormGroup({
            key: new FormControl({ value: '', disabled: true }, [Validators.required]),
            displayName: new FormControl(''),
            group: new FormControl(''),
        });

        this.formGroup.patchValue(data.role);
    }

    submitForm(): void {
        if (this.formGroup.valid && this.key?.value && this.group?.value && this.displayName?.value) {
            this.mgmtService.ChangeProjectRole(this.projectId, this.key.value, this.displayName.value, this.group.value)
                .then(() => {
                    this.toast.showInfo('PROJECT.TOAST.ROLECHANGED', true);
                    this.dialogRef.close(true);
                }).catch(error => {
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
