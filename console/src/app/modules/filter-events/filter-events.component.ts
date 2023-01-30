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
import { AbstractControl, FormBuilder, FormControl, FormGroup } from '@angular/forms';
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

  @Input() public request: ListEventsRequest = new ListEventsRequest();

  public aggregateTypes: AggregateType.AsObject[] = [];
  public eventTypes: Array<EventType.AsObject> = [];

  public isLoading: boolean = false;

  @Output() public requestChanged: EventEmitter<ListEventsRequest> = new EventEmitter();
  private destroy$: Subject<void> = new Subject();
  public form = new FormGroup({
    userFilterSet: new FormControl(false),
    editorUserId: new FormControl(''),
    aggregateFilterSet: new FormControl(false),
    aggregateId: new FormControl(''),
    aggregateTypesList: new FormControl<string[]>([]),
    eventTypesFilterSet: new FormControl(false),
    eventTypesList: new FormControl<string[]>([]),
  });

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
          this.aggregateFilterSet?.setValue(true);
        }
        if (filters.aggregateTypesList && filters.aggregateTypesList.length) {
          this.request.setAggregateTypesList(filters.aggregateTypesList);
          this.aggregateFilterSet?.setValue(true);
        }
        if (filters.editorUserId) {
          this.request.setEditorUserId(filters.editorUserId);
          this.userFilterSet?.setValue(true);
        }
        if (filters.eventTypesList && filters.eventTypesList.length) {
          this.request.setEventTypesList(filters.eventTypesList);
          this.eventTypeFilterSet?.setValue(true);
        }
        this.emitChange();
        this.cdref.detectChanges();
      }
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

  //   public resetAggregateValues(event: MatCheckboxChange): void {
  //     if (!event.checked) {
  //       this.request.setAggregateId('');
  //       this.request.setAggregateTypesList([]);
  //       this.emitChange();
  //     }
  //   }

  //   public resetUserValues(event: MatCheckboxChange): void {
  //     if (!event.checked) {
  //       this.request.setEditorUserId('');
  //       this.emitChange();
  //     }
  //   }

  //   public resetTypeValues(event: MatCheckboxChange): void {
  //     if (!event.checked) {
  //       this.request.setEventTypesList([]);
  //       this.emitChange();
  //     }
  //   }

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

  public compareTypes(t1: string, t2: string) {
    console.log(t1, t2);
    if (t1 && t2) {
      return t1 === t2;
    }
    return false;
  }

  public emitChange(): void {
    console.log(this.form.value);
    const formValues = this.form.value;

    const constructRequest = new ListEventsRequest();
    let filterObject: Partial<ListEventsRequest.AsObject> = {};

    if (formValues.userFilterSet && formValues.editorUserId) {
      constructRequest.setEditorUserId(formValues.editorUserId);
      filterObject.editorUserId = formValues.editorUserId;
    }
    if (formValues.aggregateFilterSet && formValues.aggregateTypesList) {
      constructRequest.setAggregateTypesList(formValues.aggregateTypesList);
      filterObject.aggregateTypesList = formValues.aggregateTypesList;
    }
    if (formValues.aggregateFilterSet && formValues.aggregateId) {
      constructRequest.setAggregateId(formValues.aggregateId);
      filterObject.aggregateId = formValues.aggregateId;
    }
    if (formValues.eventTypesFilterSet && formValues.eventTypesList) {
      constructRequest.setEventTypesList(formValues.eventTypesList);
      filterObject.eventTypesList = formValues.eventTypesList;
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

  public get eventTypeFilterSet(): AbstractControl | null {
    return this.form.get('eventTypeFilterSet');
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
}
