import { CommonModule } from '@angular/common';
import { Component, Inject, signal } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MAT_DIALOG_DATA, MatDialogModule, MatDialogRef } from '@angular/material/dialog';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { Group, GroupUser } from '@zitadel/proto/zitadel/group/v2/group_pb';
import { User } from 'src/app/proto/generated/zitadel/user_pb';
import { AvatarModule } from 'src/app/modules/avatar/avatar.module';
import { SearchUserAutocompleteModule } from 'src/app/modules/search-user-autocomplete/search-user-autocomplete.module';
import { UserTarget } from 'src/app/modules/search-user-autocomplete/search-user-autocomplete.component';
import { GroupService } from 'src/app/services/group.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
  selector: 'cnsl-group-members-dialog',
  templateUrl: './group-members-dialog.component.html',
  styleUrls: ['./group-members-dialog.component.scss'],
  imports: [
    AvatarModule,
    CommonModule,
    MatButtonModule,
    MatDialogModule,
    MatTooltipModule,
    TranslateModule,
    SearchUserAutocompleteModule,
  ],
})
export class GroupMembersDialogComponent {
  protected readonly members = signal<GroupUser[]>([]);
  protected readonly loading = signal(true);
  protected readonly UserTarget = UserTarget;
  protected usersToAdd: string[] = [];
  private changed = false;

  constructor(
    public dialogRef: MatDialogRef<GroupMembersDialogComponent, boolean>,
    @Inject(MAT_DIALOG_DATA) public readonly data: { group: Group },
    private readonly groupService: GroupService,
    private readonly toast: ToastService,
  ) {
    this.loadMembers().then();
  }

  protected close(): void {
    this.dialogRef.close(this.changed);
  }

  private async loadMembers(): Promise<void> {
    this.loading.set(true);
    try {
      const resp = await this.groupService.listGroupUsers({
        filters: [{ filter: { case: 'groupIds', value: { ids: [this.data.group.id] } } }],
      });
      this.members.set(resp.groupUsers);
    } catch (error) {
      this.toast.showError(error);
    } finally {
      this.loading.set(false);
    }
  }

  protected selectionChanged(users: User.AsObject[]): void {
    this.usersToAdd = users.map((user) => user.id);
  }

  protected async addUsers(): Promise<void> {
    if (!this.usersToAdd.length) {
      return;
    }
    try {
      await this.groupService.addUsersToGroup({
        id: this.data.group.id,
        userIds: this.usersToAdd,
      });
      this.toast.showInfo('GROUPS.TOAST.MEMBERSADDED', true);
      this.usersToAdd = [];
      this.changed = true;
      await new Promise((res) => setTimeout(res, 1000));
      await this.loadMembers();
    } catch (error) {
      this.toast.showError(error);
    }
  }

  protected async removeUser(member: GroupUser): Promise<void> {
    const userId = member.user?.id;
    if (!userId) {
      return;
    }
    try {
      await this.groupService.removeUsersFromGroup({
        id: this.data.group.id,
        userIds: [userId],
      });
      this.toast.showInfo('GROUPS.TOAST.MEMBERREMOVED', true);
      this.changed = true;
      await new Promise((res) => setTimeout(res, 1000));
      await this.loadMembers();
    } catch (error) {
      this.toast.showError(error);
    }
  }
}
