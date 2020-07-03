import { animate, state, style, transition, trigger } from '@angular/animations';
import { SelectionModel } from '@angular/cdk/collections';
import { AfterViewInit, Component, Input, OnInit, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatPaginator } from '@angular/material/paginator';
import { MatTable } from '@angular/material/table';
import { tap } from 'rxjs/operators';
import { ProjectGrant, ProjectMemberView } from 'src/app/proto/generated/management_pb';
import { ProjectService } from 'src/app/services/project.service';
import { ToastService } from 'src/app/services/toast.service';

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
}
