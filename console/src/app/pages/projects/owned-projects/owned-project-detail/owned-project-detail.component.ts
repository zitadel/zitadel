import { SelectionModel } from '@angular/cdk/collections';
import { Location } from '@angular/common';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatTableDataSource } from '@angular/material/table';
import { ActivatedRoute, Params, Router } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { BehaviorSubject, from, Observable, of, Subscription } from 'rxjs';
import { catchError, finalize, map } from 'rxjs/operators';
import { CreationType, MemberCreateDialogComponent } from 'src/app/modules/add-member-dialog/member-create-dialog.component';
import { ChangeType } from 'src/app/modules/changes/changes.component';
import { UserGrantContext } from 'src/app/modules/user-grants/user-grants-datasource';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import {
    Application,
    ApplicationSearchResponse,
    ProjectMember,
    ProjectMemberSearchResponse,
    ProjectMemberView,
    ProjectRole,
    ProjectRoleSearchResponse,
    ProjectState,
    ProjectType,
    ProjectView,
    UserView,
    UserGrantSearchKey,
} from 'src/app/proto/generated/management_pb';
import { OrgService } from 'src/app/services/org.service';
import { ProjectService } from 'src/app/services/project.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
    selector: 'app-owned-project-detail',
    templateUrl: './owned-project-detail.component.html',
    styleUrls: ['./owned-project-detail.component.scss'],

})
export class OwnedProjectDetailComponent implements OnInit, OnDestroy {
    public projectId: string = '';
    public project!: ProjectView.AsObject;

    public pageSizeRoles: number = 10;
    public roleDataSource: MatTableDataSource<ProjectRole.AsObject> = new MatTableDataSource<ProjectRole.AsObject>();
    public roleResult!: ProjectRoleSearchResponse.AsObject;
    public roleColumns: string[] = ['name', 'displayname', 'group', 'actions'];

    public pageSizeMembers: number = 10;
    public projectDataSource: MatTableDataSource<ProjectMember.AsObject> = new MatTableDataSource<ProjectMember.AsObject>();
    public memberResult!: ProjectMemberSearchResponse.AsObject;
    public memberColumns: string[] = ['firstname', 'lastname', 'username', 'email', 'roles'];
    public selection: SelectionModel<ProjectMember.AsObject> = new SelectionModel<ProjectMember.AsObject>(true, []);

    public pageSizeApps: number = 10;
    public appsDataSource: MatTableDataSource<Application.AsObject> = new MatTableDataSource<Application.AsObject>();
    public appsResult!: ApplicationSearchResponse.AsObject;
    public appsColumns: string[] = ['name'];

    public ProjectState: any = ProjectState;
    public ProjectType: any = ProjectType;
    public ChangeType: any = ChangeType;

    public grid: boolean = true;
    private subscription?: Subscription;
    public editstate: boolean = false;

    public isZitadel: boolean = false;

    public userGrantSearchKey: UserGrantSearchKey = UserGrantSearchKey.USERGRANTSEARCHKEY_PROJECT_ID;
    public userGrantContext: UserGrantContext = UserGrantContext.OWNED_PROJECT;

    // members
    public totalMemberResult: number = 0;
    public membersSubject: BehaviorSubject<ProjectMemberView.AsObject[]>
        = new BehaviorSubject<ProjectMemberView.AsObject[]>([]);
    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(true);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();

    constructor(
        public translate: TranslateService,
        private route: ActivatedRoute,
        private toast: ToastService,
        private projectService: ProjectService,
        private _location: Location,
        private orgService: OrgService,
        private dialog: MatDialog,
        private router: Router,
    ) {
    }

    public ngOnInit(): void {
        this.subscription = this.route.params.subscribe(params => this.getData(params));
    }

    public ngOnDestroy(): void {
        this.subscription?.unsubscribe();
    }

    private async getData({ id }: Params): Promise<void> {
        this.projectId = id;

        this.orgService.GetIam().then(iam => {
            this.isZitadel = iam.toObject().iamProjectId === this.projectId;
        });

        this.projectService.GetProjectById(id).then(proj => {
            this.project = proj.toObject();
        }).catch(error => {
            console.error(error);
            this.toast.showError(error);
        });

        from(this.projectService.SearchProjectMembers(this.projectId, 100, 0)).pipe(
            map(resp => {
                this.totalMemberResult = resp.toObject().totalResult;
                return resp.toObject().resultList;
            }),
            catchError(() => of([])),
            finalize(() => this.loadingSubject.next(false)),
        ).subscribe(members => {
            this.membersSubject.next(members);
        });
    }

    public changeState(newState: ProjectState): void {
        if (newState === ProjectState.PROJECTSTATE_ACTIVE) {
            const dialogRef = this.dialog.open(WarnDialogComponent, {
                data: {
                    confirmKey: 'ACTIONS.REACTIVATE',
                    cancelKey: 'ACTIONS.CANCEL',
                    titleKey: 'PROJECT.PAGES.DIALOG.REACTIVATE.TITLE',
                    descriptionKey: 'PROJECT.PAGES.DIALOG.REACTIVATE.DESCRIPTION',
                },
                width: '400px',
            });
            dialogRef.afterClosed().subscribe(resp => {
                if (resp) {
                    this.projectService.ReactivateProject(this.projectId).then(() => {
                        this.toast.showInfo('PROJECT.TOAST.REACTIVATED', true);
                        this.project.state = ProjectState.PROJECTSTATE_ACTIVE;
                    }).catch(error => {
                        this.toast.showError(error);
                    });
                }
            });

        } else if (newState === ProjectState.PROJECTSTATE_INACTIVE) {
            const dialogRef = this.dialog.open(WarnDialogComponent, {
                data: {
                    confirmKey: 'ACTIONS.DEACTIVATE',
                    cancelKey: 'ACTIONS.CANCEL',
                    titleKey: 'PROJECT.PAGES.DIALOG.DEACTIVATE.TITLE',
                    descriptionKey: 'PROJECT.PAGES.DIALOG.DEACTIVATE.DESCRIPTION',
                },
                width: '400px',
            });
            dialogRef.afterClosed().subscribe(resp => {
                if (resp) {
                    this.projectService.DeactivateProject(this.projectId).then(() => {
                        this.toast.showInfo('PROJECT.TOAST.DEACTIVATED', true);
                        this.project.state = ProjectState.PROJECTSTATE_INACTIVE;
                    }).catch(error => {
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
        dialogRef.afterClosed().subscribe(resp => {
            if (resp) {
                this.projectService.RemoveProject(this.projectId).then(() => {
                    this.toast.showInfo('PROJECT.TOAST.DELETED', true);
                    this.router.navigate(['/projects']);
                }).catch(error => {
                    this.toast.showError(error);
                });
            }
        });
    }

    public saveProject(): void {
        this.projectService.UpdateProject(this.project.projectId, this.project.name).then(() => {
            this.toast.showInfo('PROJECT.TOAST.UPDATED', true);
        }).catch(error => {
            this.toast.showError(error);
        });
    }

    public navigateBack(): void {
        this._location.back();
    }

    public updateName(): void {
        this.saveProject();
        this.editstate = false;
    }

    public openAddMember(): void {
        const dialogRef = this.dialog.open(MemberCreateDialogComponent, {
            data: {
                creationType: CreationType.PROJECT_OWNED,
                projectId: this.project.projectId,
            },
            width: '400px',
        });

        dialogRef.afterClosed().subscribe(resp => {
            if (resp) {
                const users: UserView.AsObject[] = resp.users;
                const roles: string[] = resp.roles;

                if (users && users.length && roles && roles.length) {
                    users.forEach(user => {
                        return this.projectService.AddProjectMember(this.projectId, user.id, roles)
                            .then(() => {
                                this.toast.showInfo('PROJECT.TOAST.MEMBERADDED', true);
                            }).catch(error => {
                                this.toast.showError(error);
                            });
                    });
                }
            }
        });
    }

    public showDetail(): void {
        if (this.project?.state === ProjectState.PROJECTSTATE_ACTIVE) {
            this.router.navigate(['projects', this.project.projectId, 'members']);
        }
    }
}
