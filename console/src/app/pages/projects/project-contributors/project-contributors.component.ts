import { Component, Input, OnInit } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { Router } from '@angular/router';
import { BehaviorSubject, from, of } from 'rxjs';
import { catchError, finalize, map } from 'rxjs/operators';
import { User } from 'src/app/proto/generated/auth_pb';
import {
    Project,
    ProjectMember,
    ProjectMemberSearchResponse,
    ProjectState,
    ProjectType,
} from 'src/app/proto/generated/management_pb';
import { ProjectService } from 'src/app/services/project.service';
import { ToastService } from 'src/app/services/toast.service';

import {
    CreationType,
    ProjectMemberCreateDialogComponent,
} from '../../../modules/add-member-dialog/project-member-create-dialog.component';

@Component({
    selector: 'app-project-contributors',
    templateUrl: './project-contributors.component.html',
    styleUrls: ['./project-contributors.component.scss'],
})
export class ProjectContributorsComponent implements OnInit {
    @Input() public project!: Project.AsObject;
    @Input() public disabled: boolean = false;

    public totalResult: number = 0;
    public membersSubject: BehaviorSubject<ProjectMember.AsObject[]> = new BehaviorSubject<ProjectMember.AsObject[]>([]);
    public ProjectState: any = ProjectState;
    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);

    constructor(private projectService: ProjectService,
        private dialog: MatDialog,
        private toast: ToastService,
        private router: Router) { }

    public ngOnInit(): void {
        const promise: Promise<ProjectMemberSearchResponse> | undefined =
            this.project.type === ProjectType.PROJECTTYPE_SELF ?
                this.projectService.SearchProjectMembers(this.project.id, 100, 0) :
                this.project.type === ProjectType.PROJECTTYPE_GRANTED ?
                    this.projectService.SearchProjectGrantMembers(this.project.id, this.project.grantId, 100, 0) : undefined;
        if (promise) {
            from(promise).pipe(
                map(resp => {
                    this.totalResult = resp.toObject().totalResult;
                    return resp.toObject().resultList;
                }),
                catchError(() => of([])),
                finalize(() => this.loadingSubject.next(false)),
            ).subscribe(members => {
                console.log(members);
                this.membersSubject.next(members);
            });
        }
    }

    public openAddMember(): void {
        const dialogRef = this.dialog.open(ProjectMemberCreateDialogComponent, {
            data: {
                creationType: this.project.type ===
                    ProjectType.PROJECTTYPE_GRANTED ? CreationType.PROJECT_GRANTED :
                    ProjectType.PROJECTTYPE_SELF ? CreationType.PROJECT_OWNED : undefined,
                projectId: this.project.id,
            },
            width: '400px',
        });

        dialogRef.afterClosed().subscribe(resp => {
            if (resp) {
                const users: User.AsObject[] = resp.users;
                const roles: string[] = resp.roles;

                if (users && users.length && roles && roles.length) {
                    Promise.all(users.map(user => {
                        return this.projectService.AddProjectMember(this.project.id, user.id, roles);
                    })).then(() => {
                        this.toast.showError('members added');
                    }).catch(error => {
                        this.toast.showError(error.message);
                    });
                }
            }
        });
    }

    public showDetail(): void {
        if (this.project?.state === ProjectState.PROJECTSTATE_ACTIVE) {
            this.router.navigate(['projects', this.project.id, 'members']);
        }
    }
}
