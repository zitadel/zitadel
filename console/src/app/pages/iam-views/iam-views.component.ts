import { AfterViewInit, Component, ViewChild } from '@angular/core';
import { MatLegacyPaginator as MatPaginator } from '@angular/material/legacy-paginator';
import { MatLegacyTableDataSource as MatTableDataSource } from '@angular/material/legacy-table';
import { MatSort } from '@angular/material/sort';
import { BehaviorSubject, from, Observable, of } from 'rxjs';
import { catchError, finalize, map } from 'rxjs/operators';
import { View } from 'src/app/proto/generated/zitadel/admin_pb';
import { AdminService } from 'src/app/services/admin.service';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';

@Component({
  selector: 'cnsl-iam-views',
  templateUrl: './iam-views.component.html',
  styleUrls: ['./iam-views.component.scss'],
})
export class IamViewsComponent implements AfterViewInit {
  @ViewChild(MatSort) sort!: MatSort;

  @ViewChild(MatPaginator) public paginator!: MatPaginator;
  public dataSource: MatTableDataSource<View.AsObject> = new MatTableDataSource<View.AsObject>([]);

  public displayedColumns: string[] = ['viewName', 'database', 'sequence', 'eventTimestamp', 'lastSuccessfulSpoolerRun'];

  private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public loading$: Observable<boolean> = this.loadingSubject.asObservable();
  constructor(private adminService: AdminService, private breadcrumbService: BreadcrumbService) {
    this.loadViews();

    const breadcrumbs = [
      new Breadcrumb({
        type: BreadcrumbType.INSTANCE,
        name: 'Instance',
        routerLink: ['/instance'],
      }),
    ];
    this.breadcrumbService.setBreadcrumb(breadcrumbs);
  }

  ngAfterViewInit(): void {
    this.loadViews();
  }

  public loadViews(): void {
    this.loadingSubject.next(true);
    from(this.adminService.listViews())
      .pipe(
        map((resp) => {
          return resp.resultList;
        }),
        catchError(() => of([])),
        finalize(() => this.loadingSubject.next(false)),
      )
      .subscribe((views) => {
        this.dataSource = new MatTableDataSource<View.AsObject>(views);
        this.dataSource.paginator = this.paginator;
        this.dataSource.sort = this.sort;
      });
  }
}
