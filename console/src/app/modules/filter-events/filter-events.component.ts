import { ConnectedPosition, ConnectionPositionPair } from '@angular/cdk/overlay';
import {
  AfterContentChecked,
  ChangeDetectorRef,
  Component,
  EventEmitter,
  Input,
  OnDestroy,
  OnInit,
  Output,
} from '@angular/core';
import { MatLegacyCheckboxChange as MatCheckboxChange } from '@angular/material/legacy-checkbox';
import { MatLegacySelectChange } from '@angular/material/legacy-select';
import { ActivatedRoute, Router } from '@angular/router';
import { Subject, take } from 'rxjs';
import { ListAggregateTypesRequest, ListEventsRequest } from 'src/app/proto/generated/zitadel/admin_pb';
import { AggregateType, EventType } from 'src/app/proto/generated/zitadel/event_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ToastService } from 'src/app/services/toast.service';
import { ActionKeysType } from '../action-keys/action-keys.component';

export enum UserTarget {
  SELF = 'self',
  EXTERNAL = 'external',
}

@Component({
  selector: 'cnsl-filter-events',
  templateUrl: './filter-events.component.html',
  styleUrls: ['./filter-events.component.scss'],
})
export class FilterEventsComponent implements OnInit, OnDestroy, AfterContentChecked {
  public showFilter: boolean = false;
  public ActionKeysType: any = ActionKeysType;

  public positions: ConnectedPosition[] = [
    new ConnectionPositionPair({ originX: 'start', originY: 'bottom' }, { overlayX: 'start', overlayY: 'top' }, 0, 10),
    new ConnectionPositionPair({ originX: 'end', originY: 'bottom' }, { overlayX: 'end', overlayY: 'top' }, 0, 10),
  ];

  public userFilterSet: boolean = false;
  public aggregateFilterSet: boolean = false;
  public eventTypeFilterSet: boolean = false;

  @Input() public request: ListEventsRequest = new ListEventsRequest();

  public aggregateTypes: AggregateType.AsObject[] = [];
  public eventTypes: Array<EventType.AsObject> = [];

  public isLoading: boolean = false;

  @Output() public requestChanged: EventEmitter<ListEventsRequest> = new EventEmitter();
  private destroy$: Subject<void> = new Subject();

  constructor(
    private adminService: AdminService,
    private toast: ToastService,
    private cdref: ChangeDetectorRef,
    private route: ActivatedRoute,
    private router: Router,
  ) {}

  public ngAfterContentChecked(): void {
    this.cdref.detectChanges();
  }

  public ngOnInit(): void {
    this.getAggregateTypes();
    this.getEventTypes();

    this.route.queryParams.pipe(take(1)).subscribe((params) => {
      const { filter } = params;
      if (filter) {
        const stringifiedFilters = filter as string;
        const filters: ListEventsRequest.AsObject = JSON.parse(stringifiedFilters);
        if (filters.aggregateId) {
          this.request.setAggregateId(filters.aggregateId);
          this.aggregateFilterSet = true;
        }
        if (filters.aggregateTypesList && filters.aggregateTypesList.length) {
          this.request.setAggregateTypesList(filters.aggregateTypesList);
          this.aggregateFilterSet = true;
        }
        if (filters.editorUserId) {
          this.request.setEditorUserId(filters.editorUserId);
          console.log('set user filter');
          this.userFilterSet = true;
        }
        if (filters.eventTypesList && filters.eventTypesList.length) {
          this.request.setEventTypesList(filters.eventTypesList);
          this.eventTypeFilterSet = true;
        }
        this.emitChange();
        this.cdref.detectChanges();
      }
    });
  }

  public reset(): void {
    this.request = new ListEventsRequest();
    this.request.setLimit(20);
    this.userFilterSet = false;
    this.aggregateFilterSet = false;
    this.eventTypeFilterSet = false;
    this.emitChange();
  }

  public finish(): void {
    this.showFilter = false;
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  //   public ngAfterContentChecked(): void {
  //     this.cdref.detectChanges();
  //   }

  private getAggregateTypes(): void {
    const req = new ListAggregateTypesRequest();

    this.adminService.listAggregateTypes(req).then((list) => {
      this.aggregateTypes = list.aggregateTypesList ?? [];
    });
  }

  private getEventTypes(): void {
    const req = new ListAggregateTypesRequest();

    this.adminService.listEventTypes(req).then((list) => {
      this.eventTypes = list.eventTypesList ?? [];
    });
  }

  public resetAggregateValues(event: MatCheckboxChange): void {
    if (!event.checked) {
      this.request.setAggregateId('');
      this.request.setAggregateTypesList([]);
      this.emitChange();
    }
  }

  public resetUserValues(event: MatCheckboxChange): void {
    if (!event.checked) {
      this.request.setEditorUserId('');
      this.emitChange();
    }
  }

  public resetTypeValues(event: MatCheckboxChange): void {
    if (!event.checked) {
      this.request.setEventTypesList([]);
      this.emitChange();
    }
  }

  public aggregateTypeObject(type: string): AggregateType.AsObject | null {
    const index = this.aggregateTypes.findIndex((agg) => agg.type === type);
    if (index > -1) {
      return this.aggregateTypes[index];
    } else {
      return null;
    }
  }

  public eventTypeObject(type: string): EventType.AsObject | null {
    const index = this.eventTypes.findIndex((eventType) => eventType.type === type);
    if (index > -1) {
      return this.eventTypes[index];
    } else {
      return null;
    }
  }

  public updateAggregateTypesList(event: MatLegacySelectChange): void {
    console.log(event.value);
    this.request.setAggregateTypesList(event.value);
    this.emitChange();
  }

  public emitChange(): void {
    console.log(this.request.toObject());
    this.requestChanged.emit(this.request);

    const req = this.request.toObject();

    let filterObject: Partial<ListEventsRequest.AsObject> = {};
    if (req.aggregateId) {
      filterObject.aggregateId = req.aggregateId;
    }
    if (req && req.aggregateTypesList.length) {
      filterObject.aggregateTypesList = req.aggregateTypesList;
    }
    if (req.editorUserId) {
      filterObject.editorUserId = req.editorUserId;
    }
    if (req.eventTypesList && req.eventTypesList.length) {
      filterObject.eventTypesList = req.eventTypesList;
    }

    if (Object.keys(filterObject).length) {
      this.router.navigate([], {
        relativeTo: this.route,
        queryParams: {
          ['filter']: JSON.stringify(filterObject),
        },
        replaceUrl: true,
        queryParamsHandling: 'merge',
        skipLocationChange: false,
      });
    } else {
      this.router.navigate([], {
        relativeTo: this.route,
        replaceUrl: true,
        skipLocationChange: false,
      });
    }
  }
}
