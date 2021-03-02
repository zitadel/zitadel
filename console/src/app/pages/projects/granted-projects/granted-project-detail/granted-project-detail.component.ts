import { Location } from '@angular/common';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { ActivatedRoute, Params, Router } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { BehaviorSubject, from, Observable, of, Subscription } from 'rxjs';
import { catchError, finalize, map } from 'rxjs/operators';
import { CreationType, MemberCreateDialogComponent } from 'src/app/modules/add-member-dialog/member-create-dialog.component';
import { ChangeType } from 'src/app/modules/changes/changes.component';
import { ProjectType } from 'src/app/modules/project-members/project-members.component';
import { UserGrantContext } from 'src/app/modules/user-grants/user-grants-datasource';
import { Member } from 'src/app/proto/generated/zitadel/member_pb';
import { GrantedProject, ProjectState } from 'src/app/proto/generated/zitadel/project_pb';
import { User } from 'src/app/proto/generated/zitadel/user_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
    selector: 'app-granted-project-detail',
    templateUrl: './granted-project-detail.component.html',
    styleUrls: ['./granted-project-detail.component.scss'],
})
export class GrantedProjectDetailComponent implements OnInit, OnDestroy {
    public projectId: string = '';
    public grantId: string = '';
    public project!: GrantedProject.AsObject;

    public ProjectState: any = ProjectState;
    public ProjectType: any = ProjectType;
    public ChangeType: any = ChangeType;

    private subscription?: Subscription;

    public isZitadel: boolean = false;

    UserGrantContext: any = UserGrantContext;

    // members
    public totalMemberResult: number = 0;
    public membersSubject: BehaviorSubject<Member.AsObject[]>
        = new BehaviorSubject<Member.AsObject[]>([]);
    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(true);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();

    constructor(
        public translate: TranslateService,
        private route: ActivatedRoute,
        private toast: ToastService,
        private mgmtService: ManagementService,
        private _location: Location,
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

        this.mgmtService.getIAM().then(iam => {
            this.isZitadel = iam.iamProjectId === this.projectId;
        });

        if (this.projectId && this.grantId) {
            this.mgmtService.getGrantedProjectByID(this.projectId, this.grantId).then(proj => {
                if (proj.grantedProject) {
                    this.project = proj.grantedProject;
                }
            }).catch(error => {
                this.toast.showError(error);
            });

            this.loadMembers();
        }
    }

    public loadMembers(): void {
        this.loadingSubject.next(true);
        from(this.mgmtService.listProjectGrantMembers(this.projectId,
            this.grantId, 100, 0)).pipe(
                map(resp => {
                    if (resp.metaData?.totalResult) {
                        this.totalMemberResult = resp.metaData.totalResult;
                    }
                    return resp.resultList;
                }),
                catchError(() => of([])),
                finalize(() => this.loadingSubject.next(false)),
            ).subscribe(members => {
                this.membersSubject.next(members);
            });
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
                const users: User.AsObject[] = resp.users;
                const roles: string[] = resp.roles;

                if (users && users.length && roles && roles.length) {
                    users.forEach(user => {
                        return this.mgmtService.addProjectGrantMember(
                            this.projectId,
                            this.grantId,
                            user.id,
                            roles,
                        ).then(() => {
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
        this.router.navigate(['granted-projects', this.project.projectId, 'grant', this.grantId, 'members']);
    }
}
