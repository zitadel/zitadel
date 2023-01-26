import { Component, Input, OnInit, ViewChild } from '@angular/core';
import { MatSort, Sort } from '@angular/material/sort';
import { MatLegacyTableDataSource as MatTableDataSource } from '@angular/material/legacy-table';
import { BehaviorSubject, Observable, Subject } from 'rxjs';
import { scan } from 'rxjs/operators';
import { ListEventsRequest, ListEventsResponse } from 'src/app/proto/generated/zitadel/admin_pb';
import { AdminService } from 'src/app/services/admin.service';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { Event } from 'src/app/proto/generated/zitadel/event_pb';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { PageEvent, PaginatorComponent } from 'src/app/modules/paginator/paginator.component';
import { LiveAnnouncer } from '@angular/cdk/a11y';
import { Router } from '@angular/router';
import { ToastService } from 'src/app/services/toast.service';
import { TranslateService } from '@ngx-translate/core';
import { ConnectedPosition, ConnectionPositionPair } from '@angular/cdk/overlay';
import { ActionKeysType } from 'src/app/modules/action-keys/action-keys.component';
import { DisplayJsonDialogComponent } from 'src/app/modules/display-json-dialog/display-json-dialog.component';
import { MatLegacyDialog as MatDialog } from '@angular/material/legacy-dialog';

enum EventFieldName {
  EDITOR = 'editor',
  AGGREGATE = 'aggregate',
  SEQUENCE = 'sequence',
  CREATIONDATE = 'creationDate',
  TYPE = 'type',
  PAYLOAD = 'payload',
}

@Component({
  selector: 'cnsl-events',
  templateUrl: './events.component.html',
  styleUrls: ['./events.component.scss'],
})
export class EventsComponent implements OnInit {
  public INITPAGESIZE = 20;

  public showFilter: boolean = false;
  public ActionKeysType: any = ActionKeysType;

  public positions: ConnectedPosition[] = [
    new ConnectionPositionPair({ originX: 'start', originY: 'bottom' }, { overlayX: 'start', overlayY: 'top' }, 0, 10),
    new ConnectionPositionPair({ originX: 'end', originY: 'bottom' }, { overlayX: 'end', overlayY: 'top' }, 0, 10),
  ];

  @ViewChild(PaginatorComponent) public paginator!: PaginatorComponent;
  @ViewChild('input') public filter!: Input;

  public dataSource: MatTableDataSource<Event.AsObject> = new MatTableDataSource<Event.AsObject>([]);
  public displayedColumns: string[] = [
    EventFieldName.TYPE,
    EventFieldName.AGGREGATE,
    EventFieldName.EDITOR,
    EventFieldName.SEQUENCE,
    EventFieldName.CREATIONDATE,
    EventFieldName.PAYLOAD,
  ];
  public timestamp: Timestamp.AsObject | undefined = undefined;
  public totalResult: number = 0;
  public filterOpen: boolean = false;

  public currentRequest: ListEventsRequest = new ListEventsRequest();

  @ViewChild(MatSort) public sort!: MatSort;

  private destroy$: Subject<void> = new Subject();
  //   private requestOrgs$: BehaviorSubject<ListEventsRequest> = new BehaviorSubject<ListEventsRequest>(initRequest);
  //   private requestOrgsObservable$ = this.requestOrgs$.pipe(takeUntil(this.destroy$));

  public _done: BehaviorSubject<any> = new BehaviorSubject(false);
  public loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  private _data: BehaviorSubject<Event.AsObject[]> = new BehaviorSubject<Event.AsObject[]>([]);
  public data!: Observable<Event.AsObject[]>;

  constructor(
    private adminService: AdminService,
    private breadcrumbService: BreadcrumbService,
    private _liveAnnouncer: LiveAnnouncer,
    private router: Router,
    private toast: ToastService,
    private translate: TranslateService,
    private dialog: MatDialog,
  ) {
    const breadcrumbs = [
      new Breadcrumb({
        type: BreadcrumbType.INSTANCE,
        name: 'Instance',
        routerLink: ['/instance'],
      }),
    ];
    this.breadcrumbService.setBreadcrumb(breadcrumbs);

    this.data = this._data.asObservable().pipe(
      scan((acc, val) => {
        console.log('scan');
        return false ? val.concat(acc) : acc.concat(val);
      }),
    );
    // .pipe(
    //   scan((acc, val) => {
    //     console.log(val);
    //     return false ? val.concat(acc) : acc.concat(val);
    //   }),
    //   tap((data) => {
    //     console.log(data);
    //     this.dataSource = new MatTableDataSource<Event.AsObject>(data);
    //   }),
  }

  ngOnInit(): void {
    const req = new ListEventsRequest();
    req.setLimit(this.INITPAGESIZE);
    this.loadEvents(req);
  }

  public loadEvents(filteredRequest: ListEventsRequest): Promise<any> | void {
    this.loadingSubject.next(true);

    // let sortingField: EventFieldName | undefined = undefined;
    // if (this.sort?.active === EventFieldName.SEQUENCE && this.sort?.direction) {
    //   sortingField = EventFieldName.SEQUENCE;
    // }

    this.currentRequest = filteredRequest;
    console.log('load', this.currentRequest.toObject());

    if (this._done.value) {
      this.loadingSubject.next(false);
      console.log('done');
      return;
    }

    return this.adminService
      .listEvents(this.currentRequest)
      .then((res: ListEventsResponse.AsObject) => {
        this._data.next(res.eventsList);

        this.loadingSubject.next(false);

        if (res.eventsList.length === 0) {
          this._done.next(true);
        }
      })
      .catch((error) => {
        console.error(error);
        this.loadingSubject.next(false);
        this._data.next([]);
      });
  }

  public refresh(): void {
    const req = new ListEventsRequest();
    req.setLimit(this.paginator.pageSize);
    // req.setSequence()
    // this.requestOrgs$.next(req);
  }

  public sortChange(sortState: Sort) {
    if (sortState.direction && sortState.active) {
      this._liveAnnouncer.announce(`Sorted ${sortState.direction}ending`);
      this.currentRequest.setAsc(sortState.direction === 'asc' ? true : false);
      this.loadEvents(this.currentRequest);
    } else {
      this._liveAnnouncer.announce('Sorting cleared');
    }
  }

  public openDialog(event: Event.AsObject): void {
    this.dialog.open(DisplayJsonDialogComponent, {
      data: {
        event: event,
      },
    });
  }

  public applySearchQuery(req: ListEventsRequest): void {
    // this.requestOrgs$.next(req);
  }

  public changePage(event: PageEvent): void {
    this.currentRequest.setLimit(event.pageSize);
    // todo use cursor to load more and add batch to local datasource
    // const lastSequenceOnNext = this.dataSource.data[this.dataSource.data.length - 1].sequence; // desc
    // const lastSequenceOnPrevious = this.dataSource.data[this.dataSource.data.length - 1].sequence; // desc

    // this.currentRequest.setSequence();
    this.loadEvents(this.currentRequest);
  }

  public more(): void {
    const sequence = this.getCursor();
    this.currentRequest.setSequence(sequence);
    this.loadEvents(this.currentRequest);
  }

  // Determines the snapshot to paginate query
  private getCursor(): number {
    const current = this._data.value;

    if (current.length) {
      const sequence = current[current.length - 1].sequence;
      return sequence;
    }
    return 0;
  }

  public gotoRouterLink(rL: any) {
    this.router.navigate(rL);
  }
}
