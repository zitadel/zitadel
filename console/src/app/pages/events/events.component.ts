import { AfterViewInit, Component, Input, ViewChild } from '@angular/core';
import { MatLegacyPaginator as MatPaginator } from '@angular/material/legacy-paginator';
import { MatSort, Sort } from '@angular/material/sort';
import { MatLegacyTableDataSource as MatTableDataSource } from '@angular/material/legacy-table';
import { BehaviorSubject, from, Observable, of, Subject } from 'rxjs';
import { catchError, finalize, map, switchMap, takeUntil } from 'rxjs/operators';
import { ListEventsRequest, ListEventTypesRequest, View } from 'src/app/proto/generated/zitadel/admin_pb';
import { AdminService } from 'src/app/services/admin.service';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { Event } from 'src/app/proto/generated/zitadel/event_pb';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { PaginatorComponent } from 'src/app/modules/paginator/paginator.component';
import { LiveAnnouncer } from '@angular/cdk/a11y';
import { Router } from '@angular/router';
import { ToastService } from 'src/app/services/toast.service';
import { TranslateService } from '@ngx-translate/core';

type Request = ListEventsRequest;

enum EventFieldName {
  EDITOR = 'editor',
  AGGREGATE = 'aggregate',
  SEQUENCE = 'sequence',
  CREATIONDATE = 'creationDate',
  TYPE = 'type',
  PAYLOAD = 'payload',
}
const initRequest = new ListEventsRequest().setLimit(20);

@Component({
  selector: 'cnsl-events',
  templateUrl: './events.component.html',
  styleUrls: ['./events.component.scss'],
})
export class EventsComponent {
  //   public orgSearchKey: OrgListSearchKey | undefined = undefined;

  @ViewChild(PaginatorComponent) public paginator!: PaginatorComponent;
  @ViewChild('input') public filter!: Input;

  public dataSource: MatTableDataSource<Event.AsObject> = new MatTableDataSource<Event.AsObject>([]);
  public displayedColumns: string[] = [
    EventFieldName.EDITOR,
    EventFieldName.AGGREGATE,
    EventFieldName.SEQUENCE,
    EventFieldName.CREATIONDATE,
    EventFieldName.TYPE,
    EventFieldName.PAYLOAD,
  ];
  private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public loading$: Observable<boolean> = this.loadingSubject.asObservable();
  public initialLimit: number = 20;
  public timestamp: Timestamp.AsObject | undefined = undefined;
  public totalResult: number = 0;
  public filterOpen: boolean = false;

  @ViewChild(MatSort) public sort!: MatSort;

  private destroy$: Subject<void> = new Subject();
  private requestOrgs$: BehaviorSubject<Request> = new BehaviorSubject<Request>(initRequest);
  private requestOrgsObservable$ = this.requestOrgs$.pipe(takeUntil(this.destroy$));

  constructor(
    private adminService: AdminService,
    private breadcrumbService: BreadcrumbService,
    private _liveAnnouncer: LiveAnnouncer,
    private router: Router,
    private toast: ToastService,
    private translate: TranslateService,
  ) {
    const breadcrumbs = [
      new Breadcrumb({
        type: BreadcrumbType.INSTANCE,
        name: 'Instance',
        routerLink: ['/instance'],
      }),
    ];
    this.breadcrumbService.setBreadcrumb(breadcrumbs);

    this.requestOrgs$.next(initRequest);

    this.requestOrgsObservable$.pipe(switchMap((req) => this.loadEvents(req))).subscribe((orgs) => {
      this.dataSource = new MatTableDataSource<Event.AsObject>(orgs);
    });

    this.load();
  }

  public loadEvents(request: Request): Observable<Event.AsObject[]> {
    this.loadingSubject.next(true);

    let sortingField: EventFieldName | undefined = undefined;
    if (this.sort?.active && this.sort?.direction) sortingField = this.sort.active as EventFieldName;
    //   switch (this.sort.active) {
    //     case 'name':
    //       sortingField = OrgFieldName.ORG_FIELD_NAME_NAME;
    //       break;
    //   }

    const req = new ListEventsRequest();

    return from(this.adminService.listEvents(req)).pipe(
      map((resp) => {
        this.totalResult = resp?.eventsList.length ?? 0;
        return resp.eventsList;
      }),
      catchError((error) => {
        this.toast.showError(error);
        return of([]);
      }),
      finalize(() => this.loadingSubject.next(false)),
    );
  }

  public load() {
    const req = new ListEventTypesRequest();

    return this.adminService.listEventTypes(req).then((list) => {
      list.eventTypesList.forEach((el) => console.log(el));
    });
  }

  public refresh(): void {
    const req = new ListEventsRequest();
    req.setLimit(this.paginator.pageSize);
    // req.setSequence()
    this.requestOrgs$.next(req);
  }

  public sortChange(sortState: Sort) {
    if (sortState.direction && sortState.active) {
      this._liveAnnouncer.announce(`Sorted ${sortState.direction}ending`);
      this.refresh();
    } else {
      this._liveAnnouncer.announce('Sorting cleared');
    }
  }

  public applySearchQuery(req: ListEventsRequest): void {
    // this.requestOrgs$.next(req);
  }

  public changePage(): void {
    this.refresh();
  }

  public gotoRouterLink(rL: any) {
    this.router.navigate(rL);
  }
}
