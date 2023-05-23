import { ConnectedPosition, ConnectionPositionPair } from '@angular/cdk/overlay';
import { Component, EventEmitter, OnInit, Output } from '@angular/core';
import { AbstractControl, FormControl, FormGroup } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { take } from 'rxjs';
import { ListAggregateTypesRequest, ListEventsRequest } from 'src/app/proto/generated/zitadel/admin_pb';
import { AggregateType, EventType } from 'src/app/proto/generated/zitadel/event_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ToastService } from 'src/app/services/toast.service';
import { ActionKeysType } from '../action-keys/action-keys.component';

export enum UserTarget {
  SELF = 'self',
  EXTERNAL = 'external',
}

function dateToTs(date: Date): Timestamp {
  const ts = new Timestamp();
  const milliseconds = date.getTime();
  const seconds = Math.abs(milliseconds / 1000);
  const nanos = (milliseconds - seconds * 1000) * 1000 * 1000;
  ts.setSeconds(seconds);
  ts.setNanos(nanos);
  return ts;
}

@Component({
  selector: 'cnsl-filter-events',
  templateUrl: './filter-events.component.html',
  styleUrls: ['./filter-events.component.scss'],
})
export class FilterEventsComponent implements OnInit {
  public showFilter: boolean = false;
  public ActionKeysType: any = ActionKeysType;

  public positions: ConnectedPosition[] = [
    new ConnectionPositionPair({ originX: 'start', originY: 'bottom' }, { overlayX: 'start', overlayY: 'top' }, 0, 10),
    new ConnectionPositionPair({ originX: 'end', originY: 'bottom' }, { overlayX: 'end', overlayY: 'top' }, 0, 10),
  ];

  private request: ListEventsRequest = new ListEventsRequest();

  public aggregateTypes: AggregateType.AsObject[] = [];
  public eventTypes: Array<EventType.AsObject> = [];

  public isLoading: boolean = false;

  @Output() public requestChanged: EventEmitter<ListEventsRequest> = new EventEmitter();
  public form: FormGroup = new FormGroup({
    resourceOwnerFilterSet: new FormControl(false),
    resourceOwner: new FormControl(''),
    sequenceFilterSet: new FormControl(false),
    sequence: new FormControl(''),
    isAsc: new FormControl<boolean>(false),
    creationDateFilterSet: new FormControl(false),
    creationDate: new FormControl<Date>(new Date()),
    userFilterSet: new FormControl(false),
    editorUserId: new FormControl(''),
    aggregateFilterSet: new FormControl(false),
    aggregateId: new FormControl(''),
    aggregateTypesList: new FormControl<AggregateType.AsObject[]>([]),
    eventTypesFilterSet: new FormControl(false),
    eventTypesList: new FormControl<EventType.AsObject[]>([]),
  });

  constructor(
    private adminService: AdminService,
    private toast: ToastService,
    private route: ActivatedRoute,
    private router: Router,
  ) {}

  public ngOnInit(): void {
    this.loadAvailableTypes().then(() => {
      this.route.queryParams.pipe(take(1)).subscribe((params) => {
        this.loadAvailableTypes().then(() => {
          const { filter } = params;
          if (filter) {
            const stringifiedFilters = filter as string;
            const filters = JSON.parse(stringifiedFilters);

            if (filters.aggregateId) {
              this.request.setAggregateId(filters.aggregateId);
              this.aggregateId?.setValue(filters.aggregateId);
              this.aggregateFilterSet?.setValue(true);
            }
            if (filters.creationDate) {
              const milliseconds = filters.creationDate;
              const date = new Date(milliseconds);
              const ts = dateToTs(date);
              this.request.setCreationDate(ts);
              this.creationDate?.setValue(date);
              this.creationDateFilterSet?.setValue(true);
            }
            if (filters.aggregateTypesList && filters.aggregateTypesList.length) {
              const values = this.aggregateTypes.filter((agg) => filters.aggregateTypesList.includes(agg.type));
              this.request.setAggregateTypesList(filters.aggregateTypesList);
              this.aggregateTypesList?.setValue(values);
              this.aggregateFilterSet?.setValue(true);
            }
            if (filters.editorUserId) {
              this.request.setEditorUserId(filters.editorUserId);
              this.editorUserId?.setValue(filters.editorUserId);
              this.userFilterSet?.setValue(true);
            }
            if (filters.resourceOwner) {
              this.request.setResourceOwner(filters.resourceOwner);
              this.resourceOwner?.setValue(filters.resourceOwner);
              this.resourceOwnerFilterSet?.setValue(true);
            }
            if (filters.sequence) {
              this.request.setSequence(filters.sequence);
              this.sequence?.setValue(filters.sequence);
              this.sequenceFilterSet?.setValue(true);
            }
            if (filters.isAsc) {
              this.request.setAsc(filters.isAsc);
              this.isAsc?.setValue(filters.isAsc);
            }
            if (filters.eventTypesList && filters.eventTypesList.length) {
              const values = this.eventTypes.filter((ev) => filters.eventTypesList.includes(ev.type));
              this.request.setEventTypesList(filters.eventTypesList);
              this.eventTypesList?.setValue(values);
              this.eventTypesFilterSet?.setValue(true);
            }
            this.emitChange();
          }
        });
      });
    });
  }

  private loadAvailableTypes(): Promise<void> {
    this.isLoading = true;
    const aT = this.getAggregateTypes();
    const eT = this.getEventTypes();
    return Promise.all([aT, eT])
      .then(() => {
        this.isLoading = false;
      })
      .catch(() => {
        this.isLoading = false;
      });
  }

  public reset(): void {
    this.form.reset();
    this.emitChange();
  }

  public finish(): void {
    this.showFilter = false;
    this.emitChange();
  }

  private getAggregateTypes(): Promise<void> {
    const req = new ListAggregateTypesRequest();

    return this.adminService
      .listAggregateTypes(req)
      .then((list) => {
        this.aggregateTypes = list.aggregateTypesList ?? [];
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  private getEventTypes(): Promise<void> {
    const req = new ListAggregateTypesRequest();

    return this.adminService
      .listEventTypes(req)
      .then((list) => {
        this.eventTypes = list.eventTypesList ?? [];
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public emitChange(): void {
    const formValues = this.form.value;

    const constructRequest = new ListEventsRequest();
    let filterObject: any = {};

    if (formValues.userFilterSet && formValues.editorUserId) {
      constructRequest.setEditorUserId(formValues.editorUserId);
      filterObject.editorUserId = formValues.editorUserId;
    }
    if (formValues.aggregateFilterSet && formValues.aggregateTypesList && formValues.aggregateTypesList.length) {
      constructRequest.setAggregateTypesList(
        formValues.aggregateTypesList.map((aggType: AggregateType.AsObject) => aggType.type),
      );
      filterObject.aggregateTypesList = formValues.aggregateTypesList.map((aggType: AggregateType.AsObject) => aggType.type);
    }
    if (formValues.aggregateFilterSet && formValues.aggregateId) {
      constructRequest.setAggregateId(formValues.aggregateId);
      filterObject.aggregateId = formValues.aggregateId;
    }
    if (formValues.eventTypesFilterSet && formValues.eventTypesList && formValues.eventTypesList.length) {
      constructRequest.setEventTypesList(formValues.eventTypesList.map((eventType: EventType.AsObject) => eventType.type));
      filterObject.eventTypesList = formValues.eventTypesList.map((eventType: EventType.AsObject) => eventType.type);
    }
    if (formValues.resourceOwnerFilterSet && formValues.resourceOwner) {
      constructRequest.setResourceOwner(formValues.resourceOwner);
      filterObject.resourceOwner = formValues.resourceOwner;
    }
    if (formValues.sequenceFilterSet && formValues.sequence) {
      constructRequest.setSequence(formValues.sequence);
      filterObject.sequence = formValues.sequence;
    }
    if (formValues.isAsc) {
      constructRequest.setAsc(formValues.isAsc);
      filterObject.isAsc = formValues.isAsc;
    }
    if (formValues.creationDateFilterSet && formValues.creationDate) {
      const date = new Date(formValues.creationDate);
      const ts = dateToTs(date);
      constructRequest.setCreationDate(ts);
      filterObject.creationDate = date.getTime();
    }

    this.requestChanged.emit(constructRequest);

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

  public get userFilterSet(): AbstractControl | null {
    return this.form.get('userFilterSet');
  }

  public get aggregateFilterSet(): AbstractControl | null {
    return this.form.get('aggregateFilterSet');
  }

  public get eventTypesFilterSet(): AbstractControl | null {
    return this.form.get('eventTypesFilterSet');
  }

  public get sequence(): AbstractControl | null {
    return this.form.get('sequence');
  }

  public get isAsc(): AbstractControl | null {
    return this.form.get('isAsc');
  }

  public get sequenceFilterSet(): AbstractControl | null {
    return this.form.get('sequenceFilterSet');
  }

  public get creationDate(): AbstractControl | null {
    return this.form.get('creationDate');
  }

  public get creationDateFilterSet(): AbstractControl | null {
    return this.form.get('creationDateFilterSet');
  }

  public get resourceOwnerFilterSet(): AbstractControl | null {
    return this.form.get('resourceOwnerFilterSet');
  }

  public get resourceOwner(): AbstractControl | null {
    return this.form.get('resourceOwner');
  }

  public get editorUserId(): AbstractControl | null {
    return this.form.get('editorUserId');
  }

  public get aggregateId(): AbstractControl | null {
    return this.form.get('aggregateId');
  }

  public get aggregateTypesList(): AbstractControl | null {
    return this.form.get('aggregateTypesList');
  }

  public get eventTypesList(): AbstractControl | null {
    return this.form.get('eventTypesList');
  }

  public get queryCount(): number {
    let count = 0;
    if (this.userFilterSet?.value && this.editorUserId?.value) {
      ++count;
    }
    if (this.creationDateFilterSet?.value && this.creationDate?.value) {
      ++count;
    }
    if (this.aggregateFilterSet?.value && this.aggregateId?.value) {
      ++count;
    }
    if (this.sequenceFilterSet?.value && this.sequence?.value) {
      ++count;
    }
    if (this.resourceOwnerFilterSet?.value && this.resourceOwner?.value) {
      ++count;
    }
    if (this.aggregateFilterSet?.value && this.aggregateTypesList?.value && this.aggregateTypesList.value.length) {
      ++count;
    }
    if (this.eventTypesFilterSet?.value && this.eventTypesList?.value && this.eventTypesList.value.length) {
      ++count;
    }
    return count;
  }
}
