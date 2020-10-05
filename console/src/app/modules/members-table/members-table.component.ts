import { SelectionModel } from '@angular/cdk/collections';
import { Component, EventEmitter, Input, OnDestroy, OnInit, Output, ViewChild } from '@angular/core';
import { MatPaginator, PageEvent } from '@angular/material/paginator';
import { MatSelectChange } from '@angular/material/select';
import { MatTable } from '@angular/material/table';
import { BehaviorSubject, Observable } from 'rxjs';
import { takeUntil } from 'rxjs/operators';
import { IamMembersDataSource } from 'src/app/pages/iam/iam-members/iam-members-datasource';
import { OrgMembersDataSource } from 'src/app/pages/orgs/org-members/org-members-datasource';
import { IamMemberView } from 'src/app/proto/generated/admin_pb';
import { OrgMemberView, ProjectMemberView } from 'src/app/proto/generated/management_pb';

import { ProjectMembersDataSource } from '../project-members/project-members-datasource';

type View = OrgMemberView.AsObject | ProjectMemberView.AsObject | IamMemberView.AsObject;
type MemberDatasource = OrgMembersDataSource | ProjectMembersDataSource | IamMembersDataSource;

@Component({
    selector: 'app-members-table',
    templateUrl: './members-table.component.html',
    styleUrls: ['./members-table.component.scss'],
})
export class MembersTableComponent implements OnInit, OnDestroy {
    public INITIALPAGESIZE: number = 25;
    @Input() public disableWrite: boolean = false;
    @ViewChild(MatPaginator) public paginator!: MatPaginator;
    @ViewChild(MatTable) public table!: MatTable<View>;
    @Input() public dataSource!: MemberDatasource;
    public selection: SelectionModel<any> = new SelectionModel<any>(true, []);
    @Input() public memberRoleOptions: string[] = [];
    @Input() public factoryLoadFunc!: Function;
    @Input() public refreshTrigger!: Observable<void>;
    @Output() public updateRoles: EventEmitter<{ member: View, change: MatSelectChange; }> = new EventEmitter();
    private destroyed: BehaviorSubject<any> = new BehaviorSubject(false);

    /** Columns displayed in the table. Columns IDs can be added, removed, or reordered. */
    public displayedColumns: string[] = ['select', 'userId', 'firstname', 'lastname', 'username', 'email', 'roles'];

    constructor() { }

    public ngOnInit(): void {
        this.refreshTrigger.pipe(takeUntil(this.destroyed)).subscribe(() => {
            this.changePage(this.paginator);
        });
    }

    public ngOnDestroy(): void {
        this.destroyed.next(true);
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

    public changePage(event?: PageEvent | MatPaginator): any {
        return this.factoryLoadFunc(event ?? this.paginator);
    }
}
