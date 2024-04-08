import { AfterViewInit, Component, ViewChild } from '@angular/core';
import { MatPaginator } from '@angular/material/paginator';
import { MatTableDataSource } from '@angular/material/table';
import { BehaviorSubject, from, Observable, of } from 'rxjs';
import { catchError, finalize, map } from 'rxjs/operators';
import { FailedEvent } from 'src/app/proto/generated/zitadel/admin_pb';
import { AdminService } from 'src/app/services/admin.service';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
  selector: 'cnsl-iam-failed-events',
  templateUrl: './failed-events.component.html',
  styleUrls: ['./failed-events.component.scss'],
})
export class FailedEventsComponent implements AfterViewInit {
  @ViewChild(MatPaginator) public eventPaginator!: MatPaginator;
  public eventDataSource: MatTableDataSource<FailedEvent.AsObject> = new MatTableDataSource<FailedEvent.AsObject>([]);

  public eventDisplayedColumns: string[] = [
    'viewName',
    'database',
    'failedSequence',
    'failureCount',
    'lastFailed',
    'errorMessage',
    'actions',
  ];

  private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public loading$: Observable<boolean> = this.loadingSubject.asObservable();
  constructor(
    private adminService: AdminService,
    private breadcrumbService: BreadcrumbService,
    private toast: ToastService,
  ) {
    this.loadEvents();

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
    this.loadEvents();
  }

  public loadEvents(): void {
    this.loadingSubject.next(true);
    from(this.adminService.listFailedEvents())
      .pipe(
        map((resp) => {
          return resp?.resultList;
        }),
        catchError(() => of([])),
        finalize(() => this.loadingSubject.next(false)),
      )
      .subscribe((events) => {
        this.eventDataSource = new MatTableDataSource<FailedEvent.AsObject>(events);
        this.eventDataSource.paginator = this.eventPaginator;
      });
  }

  public cancelEvent(viewname: string, db: string, seq: number): void {
    this.adminService.removeFailedEvent(viewname, db, seq).then(() => {
      this.toast.showInfo('IAM.FAILEDEVENTS.DELETESUCCESS', true);
    });
  }
}
