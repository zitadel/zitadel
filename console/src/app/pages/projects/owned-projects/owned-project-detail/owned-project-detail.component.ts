import { Location } from '@angular/common';
import { Component, EventEmitter, OnDestroy, OnInit } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatTableDataSource } from '@angular/material/table';
import { ActivatedRoute, Params, Router } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { BehaviorSubject, from, Observable, of, Subscription } from 'rxjs';
import { catchError, finalize, map } from 'rxjs/operators';
import { CreationType, MemberCreateDialogComponent } from 'src/app/modules/add-member-dialog/member-create-dialog.component';
import { ChangeType } from 'src/app/modules/changes/changes.component';
import { ProjectType } from 'src/app/modules/project-members/project-members.component';
import { UserGrantContext } from 'src/app/modules/user-grants/user-grants-datasource';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import { App } from 'src/app/proto/generated/zitadel/app_pb';
import { ListAppsResponse, UpdateProjectRequest } from 'src/app/proto/generated/zitadel/management_pb';
import { Member } from 'src/app/proto/generated/zitadel/member_pb';
import { Project, ProjectState } from 'src/app/proto/generated/zitadel/project_pb';
import { User } from 'src/app/proto/generated/zitadel/user_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
    selector: 'app-owned-project-detail',
    templateUrl: './owned-project-detail.component.html',
    styleUrls: ['./owned-project-detail.component.scss'],

})
export class OwnedProjectDetailComponent implements OnInit, OnDestroy {
    public projectId: string = '';
    public project!: Project.AsObject;

    public pageSizeApps: number = 10;
    public appsDataSource: MatTableDataSource<App.AsObject> = new MatTableDataSource<App.AsObject>();
    public appsResult!: ListAppsResponse.AsObject;
    public appsColumns: string[] = ['name'];

    public ProjectState: any = ProjectState;
    public ProjectType: any = ProjectType;
    public ChangeType: any = ChangeType;

    public grid: boolean = true;
    private subscription?: Subscription;
    public editstate: boolean = false;

    public isZitadel: boolean = false;

    public UserGrantContext: any = UserGrantContext;

    // members
    public totalMemberResult: number = 0;
    public membersSubject: BehaviorSubject<Member.AsObject[]>
        = new BehaviorSubject<Member.AsObject[]>([]);
    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(true);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();
    public refreshChanges$: EventEmitter<void> = new EventEmitter();

    constructor(
        public translate: TranslateService,
        private route: ActivatedRoute,
        private toast: ToastService,
        private mgmtService: ManagementService,
        private _location: Location,
        private dialog: MatDialog,
        private router: Router,
    ) { }

    public ngOnInit(): void {
        this.subscription = this.route.params.subscribe(params => this.getData(params));
    }

    public ngOnDestroy(): void {
        this.subscription?.unsubscribe();
    }

    private async getData({ id }: Params): Promise<void> {
        this.projectId = id;

        this.mgmtService.getIAM().then(iam => {
            this.isZitadel = iam.iamProjectId === this.projectId;
        });

        this.mgmtService.getProjectByID(id).then(resp => {
            if (resp.project) {
                this.project = resp.project;
            }
        }).catch(error => {
            console.error(error);
            this.toast.showError(error);
        });

        this.loadMembers();
    }

    public loadMembers(): void {
        this.loadingSubject.next(true);
        from(this.mgmtService.listProjectMembers(this.projectId, 100, 0)).pipe(
            map(resp => {
                if (resp.details?.totalResult) {
                    this.totalMemberResult = resp.details?.totalResult;
                }
                return resp.resultList;
            }),
            catchError(() => of([])),
            finalize(() => this.loadingSubject.next(false)),
        ).subscribe(members => {
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
            dialogRef.afterClosed().subscribe(resp => {
                if (resp) {
                    this.mgmtService.reactivateProject(this.projectId).then(() => {
                        this.toast.showInfo('PROJECT.TOAST.REACTIVATED', true);
                        this.project.state = ProjectState.PROJECT_STATE_ACTIVE;
                        this.refreshChanges$.emit();
                    }).catch(error => {
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
            dialogRef.afterClosed().subscribe(resp => {
                if (resp) {
                    this.mgmtService.deactivateProject(this.projectId).then(() => {
                        this.toast.showInfo('PROJECT.TOAST.DEACTIVATED', true);
                        this.project.state = ProjectState.PROJECT_STATE_INACTIVE;
                        this.refreshChanges$.emit();
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
                this.mgmtService.removeProject(this.projectId).then(() => {
                    this.toast.showInfo('PROJECT.TOAST.DELETED', true);
                    const params: Params = {
                        'deferredReload': true,
                    };
                    this.router.navigate(['/projects'], { queryParams: params });
                }).catch(error => {
                    this.toast.showError(error);
                });
            }
        });
    }

    public saveProject(): void {
        const req = new UpdateProjectRequest();
        req.setId(this.project.id);
        req.setName(this.project.name);
        req.setProjectRoleAssertion(this.project.projectRoleAssertion);
        req.setProjectRoleCheck(this.project.projectRoleCheck);

        this.mgmtService.updateProject(req).then(() => {
            this.toast.showInfo('PROJECT.TOAST.UPDATED', true);
            this.refreshChanges$.emit();
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
                projectId: this.project.id,
            },
            width: '400px',
        });

        dialogRef.afterClosed().subscribe(resp => {
            if (resp) {
                const users: User.AsObject[] = resp.users;
                const roles: string[] = resp.roles;

                if (users && users.length && roles && roles.length) {
                    users.forEach(user => {
                        return this.mgmtService.addProjectMember(this.projectId, user.id, roles)
                            .then(() => {
                                this.toast.showInfo('PROJECT.TOAST.MEMBERADDED', true);
                                setTimeout(() => {
                                    this.loadMembers();
                                }, 1000);
                            }).catch(error => {
                                this.toast.showError(error);
                            });
                    });
                }
            }
        });
    }

    public showDetail(): void {
        this.router.navigate(['projects', this.project.id, 'members']);
    }
}
