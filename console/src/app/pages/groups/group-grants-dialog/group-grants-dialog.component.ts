import { CommonModule } from '@angular/common';
import { Component, Inject, signal } from '@angular/core';
import { FormBuilder, FormControl, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MAT_DIALOG_DATA, MatDialogModule, MatDialogRef } from '@angular/material/dialog';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { Group, GroupGrant } from '@zitadel/proto/zitadel/group/v2/group_pb';
import { requiredValidator } from 'src/app/modules/form-field/validators/validators';
import { InputModule } from 'src/app/modules/input/input.module';
import { GroupService } from 'src/app/services/group.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
  selector: 'cnsl-group-grants-dialog',
  templateUrl: './group-grants-dialog.component.html',
  styleUrls: ['./group-grants-dialog.component.scss'],
  imports: [
    CommonModule,
    MatButtonModule,
    MatDialogModule,
    MatTooltipModule,
    ReactiveFormsModule,
    TranslateModule,
    InputModule,
  ],
})
export class GroupGrantsDialogComponent {
  protected readonly grants = signal<GroupGrant[]>([]);
  protected readonly loading = signal(true);
  protected readonly grantForm: ReturnType<typeof this.buildGrantForm>;

  constructor(
    public dialogRef: MatDialogRef<GroupGrantsDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public readonly data: { group: Group },
    private readonly fb: FormBuilder,
    private readonly groupService: GroupService,
    private readonly toast: ToastService,
  ) {
    this.grantForm = this.buildGrantForm();
    this.loadGrants().then();
  }

  private buildGrantForm() {
    return this.fb.group({
      projectId: new FormControl<string>('', { nonNullable: true, validators: [requiredValidator] }),
      roleKeys: new FormControl<string>('', { nonNullable: true, validators: [requiredValidator] }),
    });
  }

  private async loadGrants(): Promise<void> {
    this.loading.set(true);
    try {
      const resp = await this.groupService.listGroupGrants({
        filters: [{ filter: { case: 'groupIds', value: { ids: [this.data.group.id] } } }],
      });
      this.grants.set(resp.groupGrants);
    } catch (error) {
      this.toast.showError(error);
    } finally {
      this.loading.set(false);
    }
  }

  protected async addGrant(): Promise<void> {
    if (this.grantForm.invalid) {
      return;
    }
    const { projectId, roleKeys } = this.grantForm.getRawValue();
    try {
      await this.groupService.createGroupGrant({
        groupId: this.data.group.id,
        projectId: projectId.trim(),
        roleKeys: roleKeys
          .split(',')
          .map((key) => key.trim())
          .filter(Boolean),
      });
      this.toast.showInfo('GROUPS.TOAST.GRANTADDED', true);
      this.grantForm.reset();
      await new Promise((res) => setTimeout(res, 1000));
      await this.loadGrants();
    } catch (error) {
      this.toast.showError(error);
    }
  }

  protected async removeGrant(grant: GroupGrant): Promise<void> {
    try {
      await this.groupService.deleteGroupGrant({ id: grant.id });
      this.toast.showInfo('GROUPS.TOAST.GRANTREMOVED', true);
      await new Promise((res) => setTimeout(res, 1000));
      await this.loadGrants();
    } catch (error) {
      this.toast.showError(error);
    }
  }
}
