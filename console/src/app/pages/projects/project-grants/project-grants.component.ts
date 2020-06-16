import { animate, state, style, transition, trigger } from '@angular/animations';
import { SelectionModel } from '@angular/cdk/collections';
import { AfterViewInit, Component, Input, OnInit, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatPaginator } from '@angular/material/paginator';
import { MatTable } from '@angular/material/table';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { tap } from 'rxjs/operators';
import { ProjectGrant, ProjectMemberView } from 'src/app/proto/generated/management_pb';
import { ProjectService } from 'src/app/services/project.service';
import { ToastService } from 'src/app/services/toast.service';

import {
    ProjectGrantMembersCreateDialogComponent,
    ProjectGrantMembersCreateDialogExportType,
} from '../project-grant-members-create-dialog/project-grant-members-create-dialog.component';
import { ProjectGrantsDataSource } from './project-grants-datasource';

@Component({
    selector: 'app-project-grants',
    templateUrl: './project-grants.component.html',
    styleUrls: ['./project-grants.component.scss'],
    animations: [
        trigger('detailExpand', [
            state('collapsed', style({ height: '0px', minHeight: '0' })),
            state('expanded', style({ height: '*' })),
            transition('expanded <=> collapsed', animate('225ms cubic-bezier(0.4, 0.0, 0.2, 1)')),
        ]),
    ],
})
export class ProjectGrantsComponent implements OnInit, AfterViewInit {
    @Input() public projectId: string = '';
    @Input() public disabled: boolean = false;
    @ViewChild(MatPaginator) public paginator!: MatPaginator;
    @ViewChild(MatTable) public table!: MatTable<ProjectGrant.AsObject>;
    public dataSource!: ProjectGrantsDataSource;
    public selection: SelectionModel<ProjectGrant.AsObject> = new SelectionModel<ProjectGrant.AsObject>(true, []);
    public expandedElement: ProjectGrant.AsObject | null = null;
    public selectedGrantMembers: ProjectMemberView.AsObject[] = [];

    /** Columns displayed in the table. Columns IDs can be added, removed, or reordered. */
    public displayedColumns: string[] = ['select', 'grantedOrgName', 'grantedOrgDomain', 'creationDate', 'changeDate', 'roleNamesList'];

    constructor(private projectService: ProjectService, private toast: ToastService, private dialog: MatDialog) { }

    public ngOnInit(): void {
        this.dataSource = new ProjectGrantsDataSource(this.projectService);
        this.dataSource.loadGrants(this.projectId, 0, 25, 'asc');
    }

    public ngAfterViewInit(): void {
        this.paginator.page
            .pipe(
                tap(() => this.loadGrantsPage()),
            )
            .subscribe();

    }

    private loadGrantsPage(): void {
        this.dataSource.loadGrants(
            this.projectId,
            this.paginator.pageIndex,
            this.paginator.pageSize,
        );
    }

    public isAllSelected(): boolean {
        const numSelected = this.selection.selected.length;
        const numRows = this.dataSource.grantsSubject.value.length;
        return numSelected === numRows;
    }

    public masterToggle(): void {
        this.isAllSelected() ?
            this.selection.clear() :
            this.dataSource.grantsSubject.value.forEach(row => this.selection.select(row));
    }

    public setExpandableRow(grant: ProjectGrant.AsObject): void {
        this.expandedElement = this.expandedElement === grant ? null : grant;
        this.projectService.SearchProjectGrantMembers(this.projectId, grant.id, 10, 0).then(ret => {
            this.selectedGrantMembers = ret.toObject().resultList;
            console.log(this.selectedGrantMembers);
        });
    }

    // TODO
    public addProjectGrantMember(grant: ProjectGrant.AsObject): void {
        const dialogRef = this.dialog.open(ProjectGrantMembersCreateDialogComponent, {
            data: {
                orgId: grant.grantedOrgId,
                grantId: grant.id,
                projectId: grant.projectId,
                roleKeysList: grant.roleKeysList,
            },
            width: '400px',
        });

        console.log({
            orgId: grant.grantedOrgId,
            grantId: grant.id,
            projectId: grant.projectId,
            roleKeysList: grant.roleKeysList,
        });

        dialogRef.afterClosed().subscribe((dataToAdd: ProjectGrantMembersCreateDialogExportType) => {
            if (dataToAdd) {
                dataToAdd.userIds.forEach(userid => {
                    this.projectService.AddProjectGrantMember(
                        this.projectId,
                        grant.id,
                        userid,
                        dataToAdd.rolesKeyList,
                    ).then(() => {
                        this.toast.showInfo('Project Grant Member successfully added!');
                    }).catch(error => {
                        this.toast.showError(error.message);
                    });
                });

            }
        });
    }

    // TODO
    public removeProjectGrantMember(grantId: string, userId: string): void {
        this.projectService.RemoveProjectGrantMember(this.projectId, grantId, userId).then(() => {
            this.toast.showInfo('Project Grant Member successfully removed');
        }).catch(error => {
            this.toast.showInfo(error.message);
        });
    }

    public dateFromTimestamp(date: Timestamp.AsObject): any {
        const ts: Date = new Date(date.seconds * 1000 + date.nanos / 1000);
        return ts;
    }
}
