import { COMMA, ENTER } from '@angular/cdk/keycodes';
import {
  AfterContentChecked,
  ChangeDetectorRef,
  Component,
  ElementRef,
  EventEmitter,
  Input,
  OnDestroy,
  OnInit,
  Output,
  ViewChild,
} from '@angular/core';
import { AbstractControl, UntypedFormBuilder, UntypedFormControl, UntypedFormGroup, Validators } from '@angular/forms';
import { MatLegacyAutocomplete as MatAutocomplete } from '@angular/material/legacy-autocomplete';
import { Observable, Subject, takeUntil } from 'rxjs';
import { ListAggregateTypesRequest, ListEventsRequest } from 'src/app/proto/generated/zitadel/admin_pb';
import { AggregateType, EventType } from 'src/app/proto/generated/zitadel/event_pb';
import { User } from 'src/app/proto/generated/zitadel/user_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ToastService } from 'src/app/services/toast.service';

export enum UserTarget {
  SELF = 'self',
  EXTERNAL = 'external',
}

@Component({
  selector: 'cnsl-filter-events',
  templateUrl: './filter-events.component.html',
  styleUrls: ['./filter-events.component.scss'],
})
export class FilterEventsComponent implements OnInit, OnDestroy {
  @Input() public initialRequest: ListEventsRequest = new ListEventsRequest();

  public userFilterSet: boolean = false;
  public aggregateFilterSet: boolean = false;
  public aggregateTypes: Array<AggregateType.AsObject> = [];

  public typeFilterSet: boolean = false;
  public eventTypes: Array<EventType.AsObject> = [];

  public isLoading: boolean = false;
  @Output() public requestChanged: EventEmitter<ListEventsRequest> = new EventEmitter();
  private destroy$: Subject<void> = new Subject();
  public form = this.fb.group({
    aggregateId: ['', []],
    aggregateTypesList: [[], []],
    eventTypesList: [[], []],
    userId: ['', []],
  });

  constructor(
    private fb: UntypedFormBuilder,
    private adminService: AdminService,
    private toast: ToastService, // private cdref: ChangeDetectorRef,
  ) {
    this.form.valueChanges.pipe(takeUntil(this.destroy$)).subscribe((value) => {
      const req = new ListEventsRequest();
      if (value.aggregateId) {
        req.setAggregateId(value.aggregateId);
      }
      if (value.aggregateTypesList) {
        req.setAggregateTypesList((value.aggregateTypesList as AggregateType.AsObject[]).map((type) => type.type));
      }
      if (value.eventTypesList) {
        console.log(value.eventTypesList);
        req.setEventTypesList((value.eventTypesList as EventType.AsObject[]).map((type) => type.type));
      }
      if (value.userId) {
        req.setEditorUserId(value.userId);
      }

      console.log(req.toObject());
      this.requestChanged.emit(req);
    });
  }

  public ngOnInit(): void {
    console.log('init');
    this.getAggregateTypes();
    this.getEventTypes();

    if (this.initialRequest) {
      const req = this.initialRequest.toObject();

      if (req.editorUserId) {
        this.userFilterSet = true;
        this.userIdFormValue?.setValue(req.editorUserId);
      }

      if (req.eventTypesList.length) {
        this.typeFilterSet = true;
        this.typesFormValue?.setValue(req.eventTypesList);
      }

      if (req.aggregateTypesList.length) {
        this.aggregateFilterSet = true;
        this.aggregateTypesFormValue?.setValue(req.aggregateTypesList);
      }
    }
  }

  public reset(): void {
    this.form.reset();
    this.typeFilterSet = false;
    this.userFilterSet = false;
    this.aggregateFilterSet = false;
  }

  public finish(): void {}

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

  //   public add(event: MatChipInputEvent): void {
  //     if (!this.matAutocomplete.isOpen) {
  //       const input = event.chipInput?.inputElement;
  //       const value = event.value;

  //       if ((value || '').trim()) {
  //         const index = this.aggregateTypes.findIndex((user) => {
  //           if (user.preferredLoginName) {
  //             return user.preferredLoginName === value;
  //           } else {
  //             return false;
  //           }
  //         });
  //         if (index > -1) {
  //           if (this.users && this.users.length > 0) {
  //             this.users.push(this.filteredUsers[index]);
  //             this.selectionChanged.emit(this.users);
  //           } else {
  //             this.users = [this.filteredUsers[index]];
  //             this.selectionChanged.emit(this.users);
  //           }
  //         }
  //       }

  //       if (input) {
  //         input.value = '';
  //       }
  //     }
  //   }

  //   public remove(user: User.AsObject): void {
  //     const index = this.users.indexOf(user);

  //     if (index >= 0) {
  //       this.users.splice(index, 1);
  //       this.selectionChanged.emit(this.users);
  //     }
  //   }

  //   public selected(event: MatAutocompleteSelectedEvent): void {
  //     const index = this.filteredUsers.findIndex((user) => user === event.option.value);
  //     if (index !== -1) {
  //       if (this.singleOutput) {
  //         this.selectionChanged.emit([this.filteredUsers[index]]);
  //       } else {
  //         if (this.users && this.users.length > 0) {
  //           this.users.push(this.filteredUsers[index]);
  //         } else {
  //           this.users = [this.filteredUsers[index]];
  //         }

  //         this.selectionChanged.emit(this.users);

  //         this.usernameInput.nativeElement.value = '';
  //         this.myControl.setValue(null);
  //       }
  //     }
  //   }

  public get userIdFormValue(): AbstractControl | null {
    return this.form.get('userId');
  }
  public get aggregateTypesFormValue(): AbstractControl | null {
    return this.form.get('aggregateTypesList');
  }

  public get typesFormValue(): AbstractControl | null {
    return this.form.get('eventTypesList');
  }
}
