import { Component, ViewChild } from '@angular/core';
import { MatPaginator } from '@angular/material/paginator';
import { MatSort } from '@angular/material/sort';
import { MatTable, MatTableDataSource } from '@angular/material/table';
import { BehaviorSubject, from, Observable, of } from 'rxjs';
import { catchError, finalize, map } from 'rxjs/operators';
import { View } from 'src/app/proto/generated/admin_pb';
import { AdminService } from 'src/app/services/admin.service';

@Component({
    selector: 'app-iam-views',
    templateUrl: './iam-views.component.html',
    styleUrls: ['./iam-views.component.scss'],
})
export class IamViewsComponent {
    public views: View.AsObject[] = [];


    @ViewChild(MatPaginator) public paginator!: MatPaginator;
    @ViewChild(MatTable) public table!: MatTable<View.AsObject>;
    @ViewChild(MatSort) public sort!: MatSort;
    public dataSource!: MatTableDataSource<View.AsObject>;

    public displayedColumns: string[] = ['viewName', 'database', 'sequence', 'actions'];
    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();
    constructor(private adminService: AdminService) {
        this.loadViews();
    }

    public loadViews(): void {
        this.loadingSubject.next(true);
        from(this.adminService.GetViews()).pipe(
            map(resp => {
                return resp.toObject().viewsList;
            }),
            catchError(() => of([])),
            finalize(() => this.loadingSubject.next(false)),
        ).subscribe(views => {
            this.dataSource = new MatTableDataSource(views);
            this.dataSource.paginator = this.paginator;
            this.dataSource.sort = this.sort;
        });
    }

    public cancelView(viewname: string, db: string): void {
        this.adminService.ClearView(viewname, db);
    }
}
