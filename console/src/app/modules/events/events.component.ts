import { LiveAnnouncer } from '@angular/cdk/a11y';
import { Component, OnDestroy, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatSort, Sort } from '@angular/material/sort';
import { MatTableDataSource } from '@angular/material/table';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { BehaviorSubject, Observable, Subject, takeUntil } from 'rxjs';
import { DisplayJsonDialogComponent } from 'src/app/modules/display-json-dialog/display-json-dialog.component';
import { PaginatorComponent } from 'src/app/modules/paginator/paginator.component';
import { ListEventsRequest, ListEventsResponse } from 'src/app/proto/generated/zitadel/admin_pb';
import { Event } from 'src/app/proto/generated/zitadel/event_pb';
import { AdminService } from 'src/app/services/admin.service';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ToastService } from 'src/app/services/toast.service';

enum EventFieldName {
  EDITOR = 'editor',
  AGGREGATE = 'aggregate',
  RESOURCEOWNER = 'resourceOwner',
  SEQUENCE = 'sequence',
  CREATIONDATE = 'creationDate',
  TYPE = 'type',
  PAYLOAD = 'payload',
}

type LoadRequest = {
  req: ListEventsRequest;
  override: boolean;
};

@Component({
  selector: 'cnsl-events',
  templateUrl: './events.component.html',
  styleUrls: ['./events.component.scss'],
  standalone: false,
})
export class EventsComponent implements OnDestroy {
  public INITPAGESIZE = 20;
  public sortAsc = false;
  private destroy$: Subject<void> = new Subject();

  public displayedColumns: string[] = [
    EventFieldName.TYPE,
    EventFieldName.AGGREGATE,
    EventFieldName.RESOURCEOWNER,
    EventFieldName.EDITOR,
    EventFieldName.SEQUENCE,
    EventFieldName.CREATIONDATE,
    EventFieldName.PAYLOAD,
  ];

  // the last request received from the filter component, without paging cursor
  private filterRequest: ListEventsRequest = new ListEventsRequest();

  public currentRequest$: BehaviorSubject<LoadRequest> = new BehaviorSubject<LoadRequest>({
    req: new ListEventsRequest().setLimit(this.INITPAGESIZE),
    override: true,
  });

  @ViewChild(MatSort) public sort!: MatSort;
  @ViewChild(PaginatorComponent) public paginator!: PaginatorComponent;
  public dataSource: MatTableDataSource<Event> = new MatTableDataSource<Event>([]);

  public _done: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public done: Observable<boolean> = this._done.asObservable();

  public _loading: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);

  private _data: BehaviorSubject<Event[]> = new BehaviorSubject<Event[]>([]);

  constructor(
    private adminService: AdminService,
    private breadcrumbService: BreadcrumbService,
    private _liveAnnouncer: LiveAnnouncer,
    private toast: ToastService,
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

    this.currentRequest$
      .pipe(
        // this would compare the requests if a duplicate and redundant request would be made
        // distinctUntilChanged(({ req: prev }, { req: next }) => {
        //   return JSON.stringify(prev.toObject()) === JSON.stringify(next.toObject());
        // }),
        takeUntil(this.destroy$),
      )
      .subscribe(({ req, override }) => {
        this._loading.next(true);
        this.adminService
          .listEvents(req)
          .then((res: ListEventsResponse) => {
            if (override) {
              this._data = new BehaviorSubject<Event[]>([]);
              this.dataSource = new MatTableDataSource<Event>([]);
            }

            const eventList = res.getEventsList();
            this._data.next(eventList);

            const concat = this.dataSource.data.concat(eventList);
            this.dataSource = new MatTableDataSource<Event>(concat);

            this._loading.next(false);

            if (eventList.length === 0) {
              this._done.next(true);
            } else {
              this._done.next(false);
            }
          })
          .catch((error) => {
            this.toast.showError(error);
            this._loading.next(false);
            this._data.next([]);
          });
      });
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  public loadEvents(filteredRequest: ListEventsRequest, override: boolean = false): void {
    this.currentRequest$.next({ req: filteredRequest, override });
  }

  public refresh(): void {
    const req = new ListEventsRequest();
    req.setLimit(this.paginator.pageSize);
  }

  public sortChange(sortState: Sort) {
    if (sortState.direction && sortState.active) {
      this._liveAnnouncer.announce(`Sorted ${sortState.direction}ending`);
      this.sortAsc = sortState.direction === 'asc';

      // start from the filter request again, so the paging cursor of the previous direction is dropped
      this.loadEvents(this.buildRequest(this.filterRequest, this.sortAsc), true);
    } else {
      this._liveAnnouncer.announce('Sorting cleared');
    }
  }

  public openDialog(event: Event): void {
    this.dialog.open(DisplayJsonDialogComponent, {
      data: {
        event: event,
      },
      width: '450px',
    });
  }

  public more(): void {
    const { req } = this.currentRequest$.value;
    const cursor = this.getCursor();
    if (cursor) {
      this.applyCursor(req, cursor);
    }
    this.loadEvents(req);
  }

  public filterChanged(filterRequest: ListEventsRequest) {
    this.filterRequest = filterRequest;
    const isAsc: boolean = filterRequest.getAsc();
    if (this.sortAsc !== isAsc) {
      this.sort.sort({ id: 'creationDate', start: isAsc ? 'asc' : 'desc', disableClear: true });
    }
    this.loadEvents(this.buildRequest(filterRequest, isAsc), true);
  }

  private buildRequest(filterRequest: ListEventsRequest, asc: boolean): ListEventsRequest {
    const req = new ListEventsRequest();
    req.setLimit(this.INITPAGESIZE);
    req.setAsc(asc);

    req.setAggregateTypesList(filterRequest.getAggregateTypesList());
    req.setAggregateId(filterRequest.getAggregateId());
    req.setEventTypesList(filterRequest.getEventTypesList());
    req.setEditorUserId(filterRequest.getEditorUserId());
    req.setResourceOwner(filterRequest.getResourceOwner());
    req.setSequence(filterRequest.getSequence());
    req.setRange(filterRequest.getRange());
    req.setFrom(filterRequest.getFrom());
    return req;
  }

  /**
   * The events API pages by creation date: with asc = false the returned events are older than `from`,
   * with asc = true they are younger. `from` and `range` share a oneof, so an active range filter is
   * tightened on the cursor side instead of being replaced.
   *
   * The sequence of an event must not be used as cursor, it is only counting per aggregate.
   */
  private applyCursor(req: ListEventsRequest, cursor: Timestamp): void {
    const range = req.getRange();
    if (range) {
      if (req.getAsc()) {
        range.setSince(cursor);
      } else {
        range.setUntil(cursor);
      }
      req.setRange(range);
      return;
    }
    req.setFrom(cursor);
  }

  private getCursor(): Timestamp | undefined {
    const current = this._data.value;
    if (!current.length) {
      return undefined;
    }
    return current[current.length - 1].getCreationDate();
  }
}
