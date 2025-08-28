import { Component, EventEmitter } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { PageEvent } from '@angular/material/paginator';
import { ActionKeysType } from 'src/app/modules/action-keys/action-keys.component';
import { CreationType, MemberCreateDialogComponent } from 'src/app/modules/add-member-dialog/member-create-dialog.component';
import { Member } from 'src/app/proto/generated/zitadel/member_pb';
import { Org } from 'src/app/proto/generated/zitadel/org_pb';
import { User } from 'src/app/proto/generated/zitadel/user_pb';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { OrgMembersDataSource } from './org-members-datasource';

@Component({
  selector: 'cnsl-org-members',
  templateUrl: './org-members.component.html',
  styleUrls: ['./org-members.component.scss'],
})
export class OrgMembersComponent {
  public INITIALPAGESIZE: number = 25;
  public org?: Org.AsObject;
  public disableWrite: boolean = false;
  public dataSource!: OrgMembersDataSource;

  public memberRoleOptions: string[] = [];
  public changePageFactory!: Function;
  public changePage: EventEmitter<void> = new EventEmitter();
  public selection: Array<Member.AsObject> = [];
  public ActionKeysType: any = ActionKeysType;

  constructor(
    private mgmtService: ManagementService,
    private dialog: MatDialog,
    private toast: ToastService,
    breadcrumbService: BreadcrumbService,
  ) {
    const bread: Breadcrumb = {
      type: BreadcrumbType.ORG,
      routerLink: ['/org'],
    };
    breadcrumbService.setBreadcrumb([bread]);

    this.mgmtService.getMyOrg().then((resp) => {
      if (resp.org) {
        this.org = resp.org;
        this.dataSource = new OrgMembersDataSource(this.mgmtService);
        this.dataSource.loadMembers(0, this.INITIALPAGESIZE);
      }
    });

    this.getRoleOptions();

    this.changePageFactory = (event?: PageEvent) => {
      return this.dataSource.loadMembers(event?.pageIndex ?? 0, event?.pageSize ?? this.INITIALPAGESIZE);
    };
  }

  public getRoleOptions(): void {
    this.mgmtService
      .listOrgMemberRoles()
      .then((resp) => {
        this.memberRoleOptions = resp.resultList;
      })
      .catch((error) => {
        this.toast.showError(error);
        this.changePage.emit();
      });
  }

  updateRoles(member: Member.AsObject, selectionChange: string[]): void {
    this.mgmtService
      .updateOrgMember(member.userId, selectionChange)
      .then(() => {
        this.toast.showInfo('ORG.TOAST.MEMBERCHANGED', true);
        setTimeout(() => {
          this.changePage.emit();
        }, 1000);
      })
      .catch((error) => {
        this.changePage.emit();
        this.toast.showError(error);
      });
  }

  public removeOrgMemberSelection(): void {
    Promise.all(
      this.selection.map((member) => {
        return this.mgmtService
          .removeOrgMember(member.userId)
          .then(() => {
            this.toast.showInfo('ORG.TOAST.MEMBERREMOVED', true);
            setTimeout(() => {
              this.changePage.emit();
            }, 1000);
          })
          .catch((error) => {
            this.toast.showError(error);
          });
      }),
    )
      .then(() => {
        setTimeout(() => {
          this.changePage.emit();
        }, 1000);
      })
      .catch((error) => {
        this.toast.showError(error);
        this.changePage.emit();
      });
  }

  public removeOrgMember(member: Member.AsObject): void {
    this.mgmtService
      .removeOrgMember(member.userId)
      .then(() => {
        this.toast.showInfo('ORG.TOAST.MEMBERREMOVED', true);

        setTimeout(() => {
          this.changePage.emit();
        }, 1000);
      })
      .catch((error) => {
        this.toast.showError(error);
        this.changePage.emit();
      });
  }

  public openAddMember(): void {
    const dialogRef = this.dialog.open(MemberCreateDialogComponent, {
      data: {
        creationType: CreationType.ORG,
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((resp) => {
      if (resp) {
        const users: User.AsObject[] = resp.users;
        const roles: string[] = resp.roles;

        if (users && users.length && roles && roles.length) {
          Promise.all(
            users.map((user) => {
              return this.mgmtService.addOrgMember(user.id, roles);
            }),
          )
            .then(() => {
              this.toast.showInfo('ORG.TOAST.MEMBERADDED', true);
              setTimeout(() => {
                this.changePage.emit();
              }, 1000);
            })
            .catch((error) => {
              this.changePage.emit();
              this.toast.showError(error);
            });
        }
      }
    });
  }
}
