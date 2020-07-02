import { SelectionModel } from '@angular/cdk/collections';
import { Component, EventEmitter, Output, ViewChild } from '@angular/core';
import { MatTable } from '@angular/material/table';
import { View } from 'src/app/proto/generated/admin_pb';
import { AdminService } from 'src/app/services/admin.service';

import { IamViewsDataSource } from './iam-views.datasource';

@Component({
    selector: 'app-iam-views',
    templateUrl: './iam-views.component.html',
    styleUrls: ['./iam-views.component.scss'],
})
export class IamViewsComponent {
    public views: View.AsObject[] = [];
    @ViewChild(MatTable) public table!: MatTable<View.AsObject>;
    public dataSource!: IamViewsDataSource;
    public selection: SelectionModel<View.AsObject> = new SelectionModel<View.AsObject>(true, []);
    @Output() public changedSelection: EventEmitter<Array<View.AsObject>> = new EventEmitter();

    public displayedColumns: string[] = ['select', 'viewname', 'database', 'sequence', 'actions'];

    constructor(private adminService: AdminService) {
        this.dataSource = new IamViewsDataSource(this.adminService);
        this.dataSource.loadViews();

        this.selection.changed.subscribe(() => {
            this.changedSelection.emit(this.selection.selected);
        });
    }

    public isAllSelected(): boolean {
        const numSelected = this.selection.selected.length;
        const numRows = this.dataSource.viewsSubject.value.length;
        return numSelected === numRows;
    }

    public masterToggle(): void {
        this.isAllSelected() ?
            this.selection.clear() :
            this.dataSource.viewsSubject.value.forEach((row: View.AsObject) => this.selection.select(row));
    }

    public cancelSelectedViews(): void {

    }
}
