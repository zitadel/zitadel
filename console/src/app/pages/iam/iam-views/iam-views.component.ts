import { AfterViewInit, Component, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatPaginator } from '@angular/material/paginator';
import { MatSort } from '@angular/material/sort';
import { MatTableDataSource } from '@angular/material/table';
import { BehaviorSubject, from, Observable, of } from 'rxjs';
import { catchError, finalize, map } from 'rxjs/operators';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import { View } from 'src/app/proto/generated/admin_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
    selector: 'app-iam-views',
    templateUrl: './iam-views.component.html',
    styleUrls: ['./iam-views.component.scss'],
})
export class IamViewsComponent implements AfterViewInit {
    @ViewChild(MatSort) sort!: MatSort;

    @ViewChild(MatPaginator) public paginator!: MatPaginator;
    public dataSource!: MatTableDataSource<View.AsObject>;

    public displayedColumns: string[] = ['viewName', 'database', 'sequence', 'eventTimestamp', 'lastSuccessfulSpoolerRun', 'actions'];

    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();
    constructor(private adminService: AdminService, private dialog: MatDialog, private toast: ToastService) {
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
            this.dataSource.sort = this.sort;
        });
    }

    public cancelView(viewname: string, db: string): void {
        const dialogRef = this.dialog.open(WarnDialogComponent, {
            data: {
                confirmKey: 'ACTIONS.CLEAR',
                cancelKey: 'ACTIONS.CANCEL',
                titleKey: 'IAM.VIEWS.DIALOG.VIEW_CLEAR_TITLE',
                descriptionKey: 'IAM.VIEWS.DIALOG.VIEW_CLEAR_DESCRIPTION',
            },
            width: '400px',
        });

        dialogRef.afterClosed().subscribe(resp => {
            if (resp) {
                this.adminService.ClearView(viewname, db).then(() => {
                    this.toast.showInfo('IAM.VIEWS.CLEARED', true);
                    this.loadViews();
                }).catch(error => {
                    this.toast.showError(error);
                });
            }
        });
    }
}
