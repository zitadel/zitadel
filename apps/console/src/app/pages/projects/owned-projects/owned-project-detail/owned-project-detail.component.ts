import { Location } from '@angular/common';
import { Component, EventEmitter, OnInit } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { ActivatedRoute, Params, Router } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { BehaviorSubject, from, Observable, of } from 'rxjs';
import { catchError, finalize, map, take } from 'rxjs/operators';
import { CreationType, MemberCreateDialogComponent } from 'src/app/modules/add-member-dialog/member-create-dialog.component';
import { ChangeType } from 'src/app/modules/changes/changes.component';
import { ProjectPrivateLabelingDialogComponent } from 'src/app/modules/project-private-labeling-dialog/project-private-labeling-dialog.component';
import { SidenavSetting } from 'src/app/modules/sidenav/sidenav.component';
import { UserGrantContext } from 'src/app/modules/user-grants/user-grants-datasource';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import { UpdateProjectRequest } from 'src/app/proto/generated/zitadel/management_pb';
import { Member } from 'src/app/proto/generated/zitadel/member_pb';
import { PrivateLabelingSetting, Project, ProjectState } from 'src/app/proto/generated/zitadel/project_pb';
import { User } from 'src/app/proto/generated/zitadel/user_pb';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { NameDialogComponent } from 'src/app/modules/name-dialog/name-dialog.component';

const ROUTEPARAM = 'projectid';

const GENERAL: SidenavSetting = { id: 'general', i18nKey: 'USER.SETTINGS.GENERAL' };
const ROLES: SidenavSetting = { id: 'roles', i18nKey: 'MENU.ROLES' };
const PROJECTGRANTS: SidenavSetting = { id: 'projectgrants', i18nKey: 'MENU.PROJECTGRANTS' };
const GRANTS: SidenavSetting = { id: 'grants', i18nKey: 'MENU.GRANTS' };

@Component({
  selector: 'cnsl-owned-project-detail',
  templateUrl: './owned-project-detail.component.html',
  styleUrls: ['./owned-project-detail.component.scss'],
})
export class OwnedProjectDetailComponent implements OnInit {
  public projectId: string = '';
  public project?: Project.AsObject;

  public ProjectState: any = ProjectState;
  public ChangeType: any = ChangeType;

  public grid: boolean = true;

  public isZitadel: boolean = false;

  public UserGrantContext: any = UserGrantContext;

  // members
  public totalMemberResult: number = 0;
  public membersSubject: BehaviorSubject<Member.AsObject[]> = new BehaviorSubject<Member.AsObject[]>([]);
  private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(true);
  public loading$: Observable<boolean> = this.loadingSubject.asObservable();
  public refreshChanges$: EventEmitter<void> = new EventEmitter();

  public settingsList: SidenavSetting[] = [GENERAL, ROLES, PROJECTGRANTS, GRANTS];
  public currentSetting = this.settingsList[0];

  constructor(
    public translate: TranslateService,
    private route: ActivatedRoute,
    private toast: ToastService,
    private mgmtService: ManagementService,
    private dialog: MatDialog,
    private router: Router,
    private breadcrumbService: BreadcrumbService,
  ) {
    route.queryParamMap.pipe(take(1)).subscribe((params) => {
      const id = params.get('id');
      if (!id) {
        return;
      }
      const setting = this.settingsList.find((setting) => setting.id === id);
      if (setting) {
        this.currentSetting = setting;
      }
    });
  }

  ngOnInit(): void {
    const projectId = this.route.snapshot.paramMap.get(ROUTEPARAM);
    if (projectId) {
      this.projectId = projectId;
      this.getData(projectId).then();
    }
  }

  public openNameDialog(): void {
    const dialogRef = this.dialog.open(NameDialogComponent, {
      data: {
        name: this.project?.name,
        titleKey: 'PROJECT.NAMEDIALOG.TITLE',
        descKey: 'PROJECT.NAMEDIALOG.DESCRIPTION',
        labelKey: 'PROJECT.NAMEDIALOG.NAME',
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((name) => {
      if (name) {
        this.project!.name = name;
        this.updateName();
      }
    });
  }

  public openPrivateLabelingDialog(): void {
    const dialogRef = this.dialog.open(ProjectPrivateLabelingDialogComponent, {
      data: {
        setting: this.project?.privateLabelingSetting,
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((resp: PrivateLabelingSetting) => {
      if (resp !== undefined) {
        this.project!.privateLabelingSetting = resp;
      }
    });
  }

  private async getData(projectId: string): Promise<void> {
    this.mgmtService
      .getProjectByID(projectId)
      .then((resp) => {
        if (resp.project) {
          this.project = resp.project;

          this.mgmtService.getIAM().then((iam) => {
            this.isZitadel = iam.iamProjectId === this.projectId;

            const breadcrumbs = [
              new Breadcrumb({
                type: BreadcrumbType.ORG,
                routerLink: ['/org'],
              }),
              new Breadcrumb({
                type: BreadcrumbType.PROJECT,
                name: this.project?.name,
                param: { key: ROUTEPARAM, value: projectId },
                routerLink: ['/projects', projectId],
                isZitadel: this.isZitadel,
              }),
            ];
            this.breadcrumbService.setBreadcrumb(breadcrumbs);
          });
        }
      })
      .catch((error) => {
        console.error(error);
        this.toast.showError(error);
      });

    this.loadMembers();
  }

  public loadMembers(): void {
    this.loadingSubject.next(true);
    from(this.mgmtService.listProjectMembers(this.projectId, 100, 0))
      .pipe(
        map((resp) => {
          if (resp.details?.totalResult) {
            this.totalMemberResult = resp.details?.totalResult;
          } else {
            this.totalMemberResult = 0;
          }
          return resp.resultList;
        }),
        catchError(() => of([])),
        finalize(() => this.loadingSubject.next(false)),
      )
      .subscribe((members) => {
        this.membersSubject.next(members);
      });
  }

  public changeState(newState: ProjectState): void {
    if (newState === ProjectState.PROJECT_STATE_ACTIVE) {
      const dialogRef = this.dialog.open(WarnDialogComponent, {
        data: {
          confirmKey: 'ACTIONS.REACTIVATE',
          cancelKey: 'ACTIONS.CANCEL',
          titleKey: 'PROJECT.PAGES.DIALOG.REACTIVATE.TITLE',
          descriptionKey: 'PROJECT.PAGES.DIALOG.REACTIVATE.DESCRIPTION',
        },
        width: '400px',
      });
      dialogRef.afterClosed().subscribe((resp) => {
        if (resp) {
          this.mgmtService
            .reactivateProject(this.projectId)
            .then(() => {
              this.toast.showInfo('PROJECT.TOAST.REACTIVATED', true);
              this.project!.state = ProjectState.PROJECT_STATE_ACTIVE;
              this.refreshChanges$.emit();
            })
            .catch((error) => {
              this.toast.showError(error);
            });
        }
      });
    } else if (newState === ProjectState.PROJECT_STATE_INACTIVE) {
      const dialogRef = this.dialog.open(WarnDialogComponent, {
        data: {
          confirmKey: 'ACTIONS.DEACTIVATE',
          cancelKey: 'ACTIONS.CANCEL',
          titleKey: 'PROJECT.PAGES.DIALOG.DEACTIVATE.TITLE',
          descriptionKey: 'PROJECT.PAGES.DIALOG.DEACTIVATE.DESCRIPTION',
        },
        width: '400px',
      });
      dialogRef.afterClosed().subscribe((resp) => {
        if (resp) {
          this.mgmtService
            .deactivateProject(this.projectId)
            .then(() => {
              this.toast.showInfo('PROJECT.TOAST.DEACTIVATED', true);
              this.project!.state = ProjectState.PROJECT_STATE_INACTIVE;
              this.refreshChanges$.emit();
            })
            .catch((error) => {
              this.toast.showError(error);
            });
        }
      });
    }
  }

  public deleteProject(): void {
    const dialogRef = this.dialog.open(WarnDialogComponent, {
      data: {
        confirmKey: 'ACTIONS.DELETE',
        cancelKey: 'ACTIONS.CANCEL',
        titleKey: 'PROJECT.PAGES.DIALOG.DELETE.TITLE',
        descriptionKey: 'PROJECT.PAGES.DIALOG.DELETE.DESCRIPTION',
      },
      width: '400px',
    });
    dialogRef.afterClosed().subscribe((resp) => {
      if (resp) {
        this.mgmtService
          .removeProject(this.projectId)
          .then(() => {
            this.toast.showInfo('PROJECT.TOAST.DELETED', true);
            const params: Params = {
              deferredReload: true,
            };
            this.router.navigate(['/projects'], { queryParams: params }).then();
          })
          .catch((error) => {
            this.toast.showError(error);
          });
      }
    });
  }

  public saveProject(): void {
    if (this.project) {
      const req = new UpdateProjectRequest();
      req.setId(this.project.id);
      req.setName(this.project.name);
      req.setProjectRoleAssertion(this.project.projectRoleAssertion);
      req.setProjectRoleCheck(this.project.projectRoleCheck);
      req.setHasProjectCheck(this.project.hasProjectCheck);
      req.setPrivateLabelingSetting(this.project.privateLabelingSetting);

      this.mgmtService
        .updateProject(req)
        .then(() => {
          this.toast.showInfo('PROJECT.TOAST.UPDATED', true);
          this.refreshChanges$.emit();
        })
        .catch((error) => {
          this.toast.showError(error);
        });
    }
  }

  public updateName(): void {
    this.saveProject();
  }

  public openAddMember(): void {
    const dialogRef = this.dialog.open(MemberCreateDialogComponent, {
      data: {
        creationType: CreationType.PROJECT_OWNED,
        projectId: this.project?.id,
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((resp) => {
      if (resp) {
        const users: User.AsObject[] = resp.users;
        const roles: string[] = resp.roles;

        if (users && users.length && roles && roles.length) {
          users.forEach((user) => {
            return this.mgmtService
              .addProjectMember(this.projectId, user.id, roles)
              .then(() => {
                this.toast.showInfo('PROJECT.TOAST.MEMBERADDED', true);
                setTimeout(() => {
                  this.loadMembers();
                }, 1000);
              })
              .catch((error) => {
                this.toast.showError(error);
              });
          });
        }
      }
    });
  }

  public showDetail(): void {
    if (this.project) {
      this.router.navigate(['projects', this.project.id, 'members']).then();
    }
  }
}
