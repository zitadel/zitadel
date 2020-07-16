import { Component, Input, OnInit } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { Router } from '@angular/router';
import { BehaviorSubject, from, of } from 'rxjs';
import { catchError, finalize, map } from 'rxjs/operators';
import {
    ProjectGrantView,
    ProjectMemberSearchResponse,
    ProjectMemberView,
    ProjectState,
    ProjectType,
    ProjectView,
    User,
} from 'src/app/proto/generated/management_pb';
import { ProjectService } from 'src/app/services/project.service';
import { ToastService } from 'src/app/services/toast.service';

import { CreationType, MemberCreateDialogComponent } from '../../modules/add-member-dialog/member-create-dialog.component';

@Component({
    selector: 'app-project-contributors',
    templateUrl: './project-contributors.component.html',
    styleUrls: ['./project-contributors.component.scss'],
})
export class ProjectContributorsComponent implements OnInit {
    @Input() public project!: ProjectView.AsObject | ProjectGrantView.AsObject;
    @Input() public grantId: string = '';
    @Input() public projectType!: ProjectType;

    @Input() public disabled: boolean = false;

    public totalResult: number = 0;
    public membersSubject: BehaviorSubject<ProjectMemberView.AsObject[]>
        = new BehaviorSubject<ProjectMemberView.AsObject[]>([]);
    public ProjectState: any = ProjectState;
    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);

    public ProjectType: any = ProjectType;

    constructor(private projectService: ProjectService,
        private dialog: MatDialog,
        private toast: ToastService,
        private router: Router) { }

    public ngOnInit(): void {
        const promise: Promise<ProjectMemberSearchResponse> | undefined =
            this.projectType === ProjectType.PROJECTTYPE_OWNED ?
                this.projectService.SearchProjectMembers(this.project.projectId, 100, 0) :
                this.projectType === ProjectType.PROJECTTYPE_GRANTED ?
                    this.projectService.SearchProjectGrantMembers(this.project.projectId,
                        this.project.projectId, 100, 0) : undefined;
        if (promise) {
            from(promise).pipe(
                map(resp => {
                    this.totalResult = resp.toObject().totalResult;
                    return resp.toObject().resultList;
                }),
                catchError(() => of([])),
                finalize(() => this.loadingSubject.next(false)),
            ).subscribe(members => {
                this.membersSubject.next(members);
            });
        }
    }

    public openAddMember(): void {
        const dialogRef = this.dialog.open(MemberCreateDialogComponent, {
            data: {
                // TODO replace
                creationType: this.projectType === ProjectType.PROJECTTYPE_OWNED ? CreationType.PROJECT_OWNED :
                    ProjectType.PROJECTTYPE_GRANTED ? CreationType.PROJECT_GRANTED : ProjectType.PROJECTTYPE_OWNED,
                projectId: this.project.projectId,
            },
            width: '400px',
        });

        dialogRef.afterClosed().subscribe(resp => {
            if (resp) {
                const users: User.AsObject[] = resp.users;
                const roles: string[] = resp.roles;

                if (users && users.length && roles && roles.length) {
                    users.forEach(user => {
                        switch (this.projectType) {
                            case ProjectType.PROJECTTYPE_OWNED:
                                return this.projectService.AddProjectMember(this.project.projectId, user.id, roles)
                                    .then(() => {
                                        this.toast.showInfo('PROJECT.TOAST.MEMBERADDED', true);
                                    }).catch(error => {
                                        this.toast.showError(error);
                                    });
                            case ProjectType.PROJECTTYPE_GRANTED:
                                return this.projectService.AddProjectGrantMember(
                                    this.project.projectId,
                                    this.grantId,
                                    user.id,
                                    roles,
                                ).then(() => {
                                    this.toast.showInfo('PROJECT.TOAST.MEMBERADDED', true);
                                }).catch(error => {
                                    this.toast.showError(error);
                                });
                        }
                    });
                }
            }
        });
    }

    public showDetail(): void {
        if (this.project?.state === ProjectState.PROJECTSTATE_ACTIVE) {
            if (this.projectType === ProjectType.PROJECTTYPE_GRANTED) {
                this.router.navigate(['granted-projects', this.project.projectId, 'grant', this.grantId, 'members']);
            } else if (this.projectType === ProjectType.PROJECTTYPE_OWNED) {
                this.router.navigate(['projects', this.project.projectId, 'members']);
            }
        }
    }
}
