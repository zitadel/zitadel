import { COMMA, ENTER } from '@angular/cdk/keycodes';
import {
  AfterContentChecked,
  ChangeDetectorRef,
  Component,
  ElementRef,
  EventEmitter,
  Input,
  OnInit,
  Output,
  ViewChild,
} from '@angular/core';
import { UntypedFormBuilder, UntypedFormControl, UntypedFormGroup, Validators } from '@angular/forms';
import { MatLegacyAutocomplete as MatAutocomplete } from '@angular/material/legacy-autocomplete';
import { ListAggregateTypesRequest } from 'src/app/proto/generated/zitadel/admin_pb';
import { AggregateType } from 'src/app/proto/generated/zitadel/event_pb';
import { User } from 'src/app/proto/generated/zitadel/user_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ToastService } from 'src/app/services/toast.service';

export enum UserTarget {
  SELF = 'self',
  EXTERNAL = 'external',
}

@Component({
  selector: 'cnsl-filter-events-aggregate',
  templateUrl: './filter-events-aggregate.component.html',
  styleUrls: ['./filter-events-aggregate.component.scss'],
})
export class FilterEventsAggregateComponent implements OnInit, AfterContentChecked {
  public myControl: UntypedFormControl = new UntypedFormControl();
  public aggregateTypes: Array<AggregateType.AsObject> = [];
  public isLoading: boolean = false;
  @Output() public selectionChanged: EventEmitter<AggregateType.AsObject[]> = new EventEmitter();

  public aggregateForm = this.fb.group({
    aggregateId: ['', []],
    aggregateType: ['', [Validators.required]],
  });

  constructor(
    private fb: UntypedFormBuilder,
    private adminService: AdminService,
    private toast: ToastService,
    private cdref: ChangeDetectorRef,
  ) {
    this.aggregateForm.valueChanges.subscribe(console.log);
  }

  public ngOnInit(): void {
    this.getAggregateTypes();
  }

  public ngAfterContentChecked(): void {
    this.cdref.detectChanges();
  }

  private getAggregateTypes(): void {
    const req = new ListAggregateTypesRequest();

    this.adminService.listAggregateTypes(req).then((list) => {
      this.aggregateTypes = list.aggregateTypesList ?? [];
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
}
