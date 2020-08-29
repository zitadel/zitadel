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
import {
    Application,
    ApplicationSearchResponse,
    ProjectGrantView,
    ProjectMember,
    ProjectMemberSearchResponse,
    ProjectMemberView,
    ProjectRole,
    ProjectRoleSearchResponse,
    ProjectState,
    ProjectType,
    UserGrantSearchKey,
    UserView,
} from 'src/app/proto/generated/management_pb';
import { OrgService } from 'src/app/services/org.service';
import { ProjectService } from 'src/app/services/project.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
    selector: 'app-granted-project-detail',
    templateUrl: './granted-project-detail.component.html',
    styleUrls: ['./granted-project-detail.component.scss'],
})
export class GrantedProjectDetailComponent implements OnInit, OnDestroy {
    public projectId: string = '';
    public grantId: string = '';
    public project!: ProjectGrantView.AsObject;

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

    public userGrantContext: UserGrantContext = UserGrantContext.GRANTED_PROJECT;
    public userGrantSearchKey: UserGrantSearchKey = UserGrantSearchKey.USERGRANTSEARCHKEY_PROJECT_ID;

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
        private router: Router,
        private dialog: MatDialog,
    ) {
    }

    public ngOnInit(): void {
        this.subscription = this.route.params.subscribe(params => this.getData(params));
    }

    public ngOnDestroy(): void {
        this.subscription?.unsubscribe();
    }

    private async getData({ id, grantId }: Params): Promise<void> {
        this.projectId = id;
        this.grantId = grantId;

        console.log(id, grantId);

        this.orgService.GetIam().then(iam => {
            this.isZitadel = iam.toObject().iamProjectId === this.projectId;
        });

        if (this.projectId && this.grantId) {
            this.projectService.GetGrantedProjectByID(this.projectId, this.grantId).then(proj => {
                this.project = proj.toObject();
            }).catch(error => {
                this.toast.showError(error);
            });

            from(this.projectService.SearchProjectGrantMembers(this.projectId,
                this.grantId, 100, 0)).pipe(
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
    }

    public navigateBack(): void {
        this._location.back();
    }

    public openAddMember(): void {
        const dialogRef = this.dialog.open(MemberCreateDialogComponent, {
            data: {
                creationType: CreationType.PROJECT_GRANTED,
            },
            width: '400px',
        });

        dialogRef.afterClosed().subscribe(resp => {
            if (resp) {
                const users: UserView.AsObject[] = resp.users;
                const roles: string[] = resp.roles;

                if (users && users.length && roles && roles.length) {
                    users.forEach(user => {
                        return this.projectService.AddProjectGrantMember(
                            this.projectId,
                            this.grantId,
                            user.id,
                            roles,
                        ).then(() => {
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
        this.router.navigate(['granted-projects', this.project.projectId, 'grant', this.grantId, 'members']);
    }
}
