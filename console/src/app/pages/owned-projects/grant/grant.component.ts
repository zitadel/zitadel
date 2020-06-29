import { SelectionModel } from '@angular/cdk/collections';
import { AfterViewInit, Component, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatPaginator } from '@angular/material/paginator';
import { MatTable } from '@angular/material/table';
import { ActivatedRoute } from '@angular/router';
import { tap } from 'rxjs/operators';
import { CreationType, MemberCreateDialogComponent } from 'src/app/modules/add-member-dialog/member-create-dialog.component';
import { ProjectMember, ProjectType, User } from 'src/app/proto/generated/management_pb';
import { ProjectService } from 'src/app/services/project.service';
import { ToastService } from 'src/app/services/toast.service';

import { GrantDataSource } from './grant-datasource';

@Component({
    selector: 'app-grant',
    templateUrl: './grant.component.html',
    styleUrls: ['./grant.component.scss'],
})
export class GrantComponent implements AfterViewInit {
    public projectid: string = '';
    public grantid: string = '';

    public projectType: ProjectType = ProjectType.PROJECTTYPE_OWNED;
    public disabled: boolean = false;
    @ViewChild(MatPaginator) public paginator!: MatPaginator;
    @ViewChild(MatTable) public table!: MatTable<ProjectMember.AsObject>;
    public dataSource!: GrantDataSource;
    public selection: SelectionModel<ProjectMember.AsObject> = new SelectionModel<ProjectMember.AsObject>(true, []);

    /** Columns displayed in the table. Columns IDs can be added, removed, or reordered. */
    public displayedColumns: string[] = ['select', 'userId', 'firstname', 'lastname', 'username', 'email', 'roles'];

    constructor(
        private projectService: ProjectService,
        private dialog: MatDialog,
        private toast: ToastService,
        private route: ActivatedRoute) {
        this.route.params.subscribe(params => {
            this.projectid = params.projectid;
            this.grantid = params.grantid;

            this.dataSource = new GrantDataSource(this.projectService);
            this.dataSource.loadMembers(params.projectid, params.grantid, 0, 25, 'asc');
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
            this.projectid,
            this.grantid,
            this.paginator.pageIndex,
            this.paginator.pageSize,
        );
    }

    public removeProjectMemberSelection(): void {
        Promise.all(this.selection.selected.map(member => {
            return this.projectService.RemoveProjectGrantMember(this.projectid, this.grantid, member.userId).then(() => {
                this.toast.showInfo('Removed successfully');
            }).catch(error => {
                this.toast.showError(error.message);
            });
        }));
    }

    public removeMember(member: ProjectMember.AsObject): void {
        this.projectService.RemoveProjectGrantMember(this.projectid, this.grantid, member.userId).then(() => {
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
        const dialogRef = this.dialog.open(MemberCreateDialogComponent, {
            data: {
                creationType: CreationType.PROJECT_OWNED,
                projectId: this.projectid,
            },
            width: '400px',
        });

        dialogRef.afterClosed().subscribe(resp => {
            if (resp) {
                const users: User.AsObject[] = resp.users;
                const roles: string[] = resp.roles;

                if (users && users.length && roles && roles.length) {
                    Promise.all(users.map(user => {
                        return this.projectService.AddProjectGrantMember(this.projectid, this.grantid, user.id, roles);
                    })).then(() => {
                        this.toast.showError('members added');
                    }).catch(error => {
                        this.toast.showError(error.message);
                    });
                }
            }
        });
    }
}
