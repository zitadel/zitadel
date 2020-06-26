import { SelectionModel } from '@angular/cdk/collections';
import { AfterViewInit, Component, ViewChild } from '@angular/core';
import { MatPaginator } from '@angular/material/paginator';
import { MatTable } from '@angular/material/table';
import { ActivatedRoute } from '@angular/router';
import { tap } from 'rxjs/operators';
import { ProjectMember, ProjectType, ProjectView } from 'src/app/proto/generated/management_pb';
import { ProjectService } from 'src/app/services/project.service';
import { ToastService } from 'src/app/services/toast.service';

import { ProjectMembersDataSource } from './project-members-datasource';

@Component({
    selector: 'app-project-members',
    templateUrl: './project-members.component.html',
    styleUrls: ['./project-members.component.scss'],
})
export class ProjectMembersComponent implements AfterViewInit {
    public project!: ProjectView.AsObject;
    public projectType: ProjectType = ProjectType.PROJECTTYPE_OWNED;
    public disabled: boolean = false;
    @ViewChild(MatPaginator) public paginator!: MatPaginator;
    @ViewChild(MatTable) public table!: MatTable<ProjectMember.AsObject>;
    public dataSource!: ProjectMembersDataSource;
    public selection: SelectionModel<ProjectMember.AsObject> = new SelectionModel<ProjectMember.AsObject>(true, []);

    /** Columns displayed in the table. Columns IDs can be added, removed, or reordered. */
    public displayedColumns: string[] = ['select', 'firstname', 'lastname', 'username', 'email', 'roles'];

    constructor(private projectService: ProjectService,
        private toast: ToastService,
        private route: ActivatedRoute) {
        this.route.params.subscribe(params => {
            this.projectService.GetProjectById(params.projectid).then(project => {
                this.project = project.toObject();
                this.dataSource = new ProjectMembersDataSource(this.projectService);
                this.dataSource.loadMembers(this.project, this.projectType, 0, 25, 'asc');
            });
        });
    }

    public ngAfterViewInit(): void {
        this.paginator.page
            .pipe(
                tap(() => this.loadMembersPage()),
            )
            .subscribe();

    }

    private loadMembersPage(): void {
        this.dataSource.loadMembers(
            this.project,
            this.projectType,
            this.paginator.pageIndex,
            this.paginator.pageSize,
        );
    }

    public removeProjectMemberSelection(): void {
        Promise.all(this.selection.selected.map(member => {
            return this.projectService.RemoveProjectMember(this.project.projectId, member.userId).then(() => {
                this.toast.showInfo('Removed successfully');
            }).catch(error => {
                this.toast.showError(error.message);
            });
        }));
    }

    public removeMember(member: ProjectMember.AsObject): void {
        this.projectService.RemoveProjectMember(this.project.projectId, member.userId).then(() => {
            this.toast.showInfo('Member removed successfully');
        }).catch(error => {
            this.toast.showError(error.message);
        });
    }

    public isAllSelected(): boolean {
        const numSelected = this.selection.selected.length;
        const numRows = this.dataSource.membersSubject.value.length;
        return numSelected === numRows;
    }

    public masterToggle(): void {
        this.isAllSelected() ?
            this.selection.clear() :
            this.dataSource.membersSubject.value.forEach(row => this.selection.select(row));
    }

    public openAddMember(): void {

        // TODO
        // const dialogRef = this.dialog.open(ProjectMemberCreateDialogComponent, {
        //     data: {
        //         creationType: this.project.type ===
        //             ProjectType.PROJECTTYPE_GRANTED ? CreationType.PROJECT_GRANTED :
        //             ProjectType.PROJECTTYPE_OWNED ? CreationType.PROJECT_OWNED : undefined,
        //         projectId: this.project.id,
        //     },
        //     width: '400px',
        // });

        // dialogRef.afterClosed().subscribe(resp => {
        //     if (resp) {
        //         const users: User.AsObject[] = resp.users;
        //         const roles: string[] = resp.roles;

        //         if (users && users.length && roles && roles.length) {
        //             Promise.all(users.map(user => {
        //                 return this.projectService.AddProjectMember(this.project.id, user.id, roles);
        //             })).then(() => {
        //                 this.toast.showError('members added');
        //             }).catch(error => {
        //                 this.toast.showError(error.message);
        //             });
        //         }
        //     }
        // });
    }
}
