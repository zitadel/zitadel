import { Component, EventEmitter } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { PageEvent } from '@angular/material/paginator';
import { ActivatedRoute } from '@angular/router';
import { ActionKeysType } from 'src/app/modules/action-keys/action-keys.component';
import { CreationType, MemberCreateDialogComponent } from 'src/app/modules/add-member-dialog/member-create-dialog.component';
import { UserGrantRoleDialogComponent } from 'src/app/modules/user-grant-role-dialog/user-grant-role-dialog.component';
import { Member } from 'src/app/proto/generated/zitadel/member_pb';
import { GrantedProject, ProjectGrantState, Role } from 'src/app/proto/generated/zitadel/project_pb';
import { User } from 'src/app/proto/generated/zitadel/user_pb';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { ProjectGrantMembersDataSource } from './project-grant-members-datasource';

@Component({
  selector: 'cnsl-project-grant-detail',
  templateUrl: './project-grant-detail.component.html',
  styleUrls: ['./project-grant-detail.component.scss'],
  standalone: false,
})
export class ProjectGrantDetailComponent {
  public INITIALPAGESIZE: number = 25;

  public grant: GrantedProject.AsObject | undefined;
  public projectid: string = '';
  public grantid: string = '';

  public disabled: boolean = false;

  public isZitadel: boolean = false;
  ProjectGrantState: any = ProjectGrantState;

  public projectRoleOptions: Role.AsObject[] = [];
  public memberRoleOptions: Array<string> = [];

  public changePageFactory!: Function;
  public changePage: EventEmitter<void> = new EventEmitter();
  public selection: Array<Member.AsObject> = [];
  public dataSource!: ProjectGrantMembersDataSource;

  public ActionKeysType: any = ActionKeysType;

  constructor(
    private mgmtService: ManagementService,
    private route: ActivatedRoute,
    private toast: ToastService,
    private dialog: MatDialog,
    private breadcrumbService: BreadcrumbService,
  ) {
    this.route.params.subscribe((params) => {
      this.projectid = params['projectid'];
      this.grantid = params['grantid'];

      this.dataSource = new ProjectGrantMembersDataSource(this.mgmtService);
      this.dataSource.loadMembers(params['projectid'], params['grantid'], 0, this.INITIALPAGESIZE);

      this.getRoleOptions(params['projectid']);
      this.getMemberRoleOptions();

      this.changePageFactory = (event?: PageEvent) => {
        return this.dataSource.loadMembers(
          params['projectid'],
          params['grantid'],
          event?.pageIndex ?? 0,
          event?.pageSize ?? this.INITIALPAGESIZE,
        );
      };

      this.mgmtService
        .getProjectGrantByID(this.grantid, this.projectid)
        .then((resp) => {
          if (resp.projectGrant) {
            this.grant = resp.projectGrant;

            const breadcrumbs = [
              new Breadcrumb({
                type: BreadcrumbType.ORG,
                routerLink: ['/org'],
              }),
              new Breadcrumb({
                type: BreadcrumbType.PROJECT,
                name: '',
                param: { key: 'projectid', value: resp.projectGrant.projectId },
                routerLink: ['/projects', resp.projectGrant.projectId],
              }),
            ];
            this.breadcrumbService.setBreadcrumb(breadcrumbs);
          }
        })
        .catch((error) => {
          this.toast.showError(error);
        });
    });
  }

  public changeState(grant: GrantedProject.AsObject, newState: ProjectGrantState): void {
    if (newState === ProjectGrantState.PROJECT_GRANT_STATE_ACTIVE) {
      this.mgmtService
        .reactivateProjectGrant(this.grantid, this.projectid)
        .then(() => {
          this.toast.showInfo('PROJECT.TOAST.REACTIVATED', true);
          grant.state = newState;
        })
        .catch((error) => {
          this.toast.showError(error);
        });
    } else if (newState === ProjectGrantState.PROJECT_GRANT_STATE_INACTIVE) {
      this.mgmtService
        .deactivateProjectGrant(this.grantid, this.projectid)
        .then(() => {
          this.toast.showInfo('PROJECT.TOAST.DEACTIVATED', true);
          grant.state = newState;
          setTimeout(() => {
            this.changePage.emit();
          }, 1000);
        })
        .catch((error) => {
          this.toast.showError(error);
          this.changePage.emit();
        });
    }
  }

  public getRoleOptions(projectId: string): void {
    this.mgmtService.listProjectRoles(projectId, 100, 0).then((resp) => {
      this.projectRoleOptions = resp.resultList;
    });
  }

  public getMemberRoleOptions(): void {
    this.mgmtService
      .listProjectGrantMemberRoles()
      .then((resp) => {
        this.memberRoleOptions = resp.resultList;
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public removeProjectMemberSelection(grant: GrantedProject.AsObject): void {
    Promise.all(
      this.selection.map((member) => {
        return this.mgmtService
          .removeProjectGrantMember(grant.projectId, grant.grantId, member.userId)
          .then(() => {
            this.toast.showInfo('PROJECT.GRANT.TOAST.PROJECTGRANTMEMBERREMOVED', true);
            setTimeout(() => {
              this.changePage.emit();
            }, 1000);
          })
          .catch((error) => {
            this.changePage.emit();
            this.toast.showError(error);
          });
      }),
    );
  }

  public removeProjectMember(grant: GrantedProject.AsObject, member: Member.AsObject): void {
    this.mgmtService
      .removeProjectGrantMember(grant.projectId, grant.grantId, member.userId)
      .then(() => {
        this.toast.showInfo('PROJECT.GRANT.TOAST.PROJECTGRANTMEMBERREMOVED', true);
        setTimeout(() => {
          this.changePage.emit();
        }, 1000);
      })
      .catch((error) => {
        this.changePage.emit();
        this.toast.showError(error);
      });
  }

  public async openAddMember(grant: GrantedProject.AsObject): Promise<any> {
    const dialogRef = this.dialog.open(MemberCreateDialogComponent, {
      data: {
        creationType: CreationType.PROJECT_GRANTED,
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((resp) => {
      if (resp) {
        const users: User.AsObject[] = resp.users;
        const roles: string[] = resp.roles;

        if (users && users.length && roles && roles.length) {
          const userIds = users.map((user) => user.id);
          Promise.all(
            userIds.map((userid: string) => {
              return this.mgmtService.addProjectGrantMember(grant.projectId, grant.grantId, userid, resp.roles);
            }),
          )
            .then(() => {
              this.toast.showInfo('PROJECT.GRANT.TOAST.PROJECTGRANTMEMBERADDED', true);
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

  updateMemberRoles(grant: GrantedProject.AsObject, member: Member.AsObject, selectionChange: string[]): void {
    this.mgmtService
      .updateProjectGrantMember(grant.projectId, grant.grantId, member.userId, selectionChange)
      .then(() => {
        setTimeout(() => {
          this.changePage.emit();
        }, 1000);
        this.toast.showInfo('PROJECT.GRANT.TOAST.PROJECTGRANTMEMBERCHANGED', true);
      })
      .catch((error) => {
        this.changePage.emit();
        this.toast.showError(error);
      });
  }

  removeRole(grant: GrantedProject.AsObject, role: string): void {
    const index = grant.grantedRoleKeysList.findIndex((r) => r === role);
    if (index > -1) {
      grant.grantedRoleKeysList.splice(index, 1);

      this.mgmtService
        .updateProjectGrant(grant.grantId, grant.projectId, grant.grantedRoleKeysList)
        .then(() => {
          setTimeout(() => {
            this.changePage.emit();
          }, 1000);
          this.toast.showInfo('PROJECT.GRANT.TOAST.PROJECTGRANTUPDATED', true);
        })
        .catch((error) => {
          this.changePage.emit();
          this.toast.showError(error);
        });
    }
  }

  public editRoles(grant: GrantedProject.AsObject): void {
    const dialogRef = this.dialog.open(UserGrantRoleDialogComponent, {
      data: {
        projectId: grant.projectId,
        selectedRoleKeysList: grant.grantedRoleKeysList,
        i18nTitle: 'PROJECT.GRANT.EDITTITLE',
      },
      width: '600px',
    });

    dialogRef.afterClosed().subscribe((resp) => {
      if (resp && resp.roles) {
        this.mgmtService
          .updateProjectGrant(grant.grantId, grant.projectId, resp.roles)
          .then(() => {
            this.toast.showInfo('PROJECT.GRANT.TOAST.PROJECTGRANTUPDATED', true);
            grant.grantedRoleKeysList = resp.roles;
            setTimeout(() => {
              this.changePage.emit();
            }, 1000);
          })
          .catch((error) => {
            this.changePage.emit();
            this.toast.showError(error);
          });
      }
    });
  }
}
