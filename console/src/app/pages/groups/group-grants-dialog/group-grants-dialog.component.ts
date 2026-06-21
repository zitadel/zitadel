import { CommonModule } from '@angular/common';
import { Component, Inject, signal } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MAT_DIALOG_DATA, MatDialogModule, MatDialogRef } from '@angular/material/dialog';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { Group, GroupGrant } from '@zitadel/proto/zitadel/group/v2/group_pb';
import { ProjectType } from 'src/app/modules/project-members/project-members-datasource';
import { ProjectRolesTableModule } from 'src/app/modules/project-roles-table/project-roles-table.module';
import { SearchProjectAutocompleteModule } from 'src/app/modules/search-project-autocomplete/search-project-autocomplete.module';
import { GrantedProject, Project } from 'src/app/proto/generated/zitadel/project_pb';
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
    TranslateModule,
    SearchProjectAutocompleteModule,
    ProjectRolesTableModule,
  ],
})
export class GroupGrantsDialogComponent {
  protected readonly grants = signal<GroupGrant[]>([]);
  protected readonly loading = signal(true);
  protected readonly selectedProjectId = signal<string>('');
  protected readonly selectedGrantId = signal<string>('');
  protected selectedRoleKeys: string[] = [];

  constructor(
    public dialogRef: MatDialogRef<GroupGrantsDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public readonly data: { group: Group },
    private readonly groupService: GroupService,
    private readonly toast: ToastService,
  ) {
    this.loadGrants().then();
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

  protected selectProject(project: Project.AsObject | GrantedProject.AsObject, type: ProjectType): void {
    if (type === ProjectType.PROJECTTYPE_OWNED) {
      this.selectedProjectId.set((project as Project.AsObject).id);
      this.selectedGrantId.set('');
    } else {
      const granted = project as GrantedProject.AsObject;
      this.selectedProjectId.set(granted.projectId);
      this.selectedGrantId.set(granted.grantId);
    }
    this.selectedRoleKeys = [];
  }

  protected selectRoles(roleKeys: string[]): void {
    this.selectedRoleKeys = roleKeys;
  }

  protected get canSave(): boolean {
    return !!this.selectedProjectId() && this.selectedRoleKeys.length > 0;
  }

  protected async addGrant(): Promise<void> {
    if (!this.canSave) {
      return;
    }
    try {
      await this.groupService.createGroupGrant({
        groupId: this.data.group.id,
        projectId: this.selectedProjectId(),
        projectGrantId: this.selectedGrantId() || undefined,
        roleKeys: this.selectedRoleKeys,
      });
      this.toast.showInfo('GROUPS.TOAST.GRANTADDED', true);
      this.selectedProjectId.set('');
      this.selectedGrantId.set('');
      this.selectedRoleKeys = [];
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
