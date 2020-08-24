import { AfterViewInit, Component, ViewChild } from '@angular/core';
import { MatPaginator } from '@angular/material/paginator';
import { MatTableDataSource } from '@angular/material/table';
import { BehaviorSubject, from, Observable, of } from 'rxjs';
import { catchError, finalize, map } from 'rxjs/operators';
import { View } from 'src/app/proto/generated/admin_pb';
import { AdminService } from 'src/app/services/admin.service';

@Component({
    selector: 'app-iam-views',
    templateUrl: './iam-views.component.html',
    styleUrls: ['./iam-views.component.scss'],
})
export class IamViewsComponent implements AfterViewInit {
    @ViewChild(MatPaginator) public paginator!: MatPaginator;
    public dataSource!: MatTableDataSource<View.AsObject>;

    public displayedColumns: string[] = ['viewName', 'database', 'sequence', 'timestamp', 'actions'];

    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();
    constructor(private adminService: AdminService) {
        this.loadViews();
    }

    ngAfterViewInit(): void {
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
        });
    }

    public cancelView(viewname: string, db: string): void {
        this.adminService.ClearView(viewname, db);
    }
}
