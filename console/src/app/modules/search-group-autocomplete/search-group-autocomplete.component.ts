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
import { UntypedFormControl } from '@angular/forms';
import { MatAutocomplete, MatAutocompleteSelectedEvent } from '@angular/material/autocomplete';
import { MatChipInputEvent } from '@angular/material/chips';
import { from, of, Subject } from 'rxjs';
import { debounceTime, switchMap, takeUntil, tap } from 'rxjs/operators';
import { ListGroupsResponse } from 'src/app/proto/generated/zitadel/management_pb';
import { TextQueryMethod } from 'src/app/proto/generated/zitadel/object_pb';
import { GroupNameQuery, GroupQuery, Group } from 'src/app/proto/generated/zitadel/group_pb';
import { LoginNameQuery, User } from 'src/app/proto/generated/zitadel/user_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

const USER_LIMIT = 25;

@Component({
  selector: 'cnsl-search-group-autocomplete',
  templateUrl: './search-group-autocomplete.component.html',
  styleUrls: ['./search-group-autocomplete.component.scss'],
})
export class SearchGroupAutocompleteComponent implements OnInit, AfterContentChecked {
  public removable: boolean = true;
  public addOnBlur: boolean = true;
  public separatorKeysCodes: number[] = [ENTER, COMMA];

  public myControl: UntypedFormControl = new UntypedFormControl();
  public globalLoginNameControl: UntypedFormControl = new UntypedFormControl();

  @Input() public groups: Array<Group.AsObject> = [];
  @Input() public editState: boolean = true;
  public filteredGroups: Array<Group.AsObject> = [];
  public isLoading: boolean = false;
  public hint: string = '';
  @ViewChild('usernameInput') public usernameInput!: ElementRef<HTMLInputElement>;
  @ViewChild('auto') public matAutocomplete!: MatAutocomplete;
  @Output() public selectionChanged: EventEmitter<Group.AsObject[]> = new EventEmitter();

  private unsubscribed$: Subject<void> = new Subject();
  constructor(
    private groupService: ManagementService,
    private toast: ToastService,
    private cdref: ChangeDetectorRef,
  ) {}

  public ngOnInit(): void {
      const query = new GroupQuery();
      const lnQuery = new GroupNameQuery();
      lnQuery.setMethod(TextQueryMethod.TEXT_QUERY_METHOD_STARTS_WITH_IGNORE_CASE);
      query.setNameQuery(lnQuery);
      this.groupService.listGroups(USER_LIMIT, 0, [query]).then((groups) => {
        this.filteredGroups = groups.resultList;
      });

      this.getFilteredResults();
  }

  public ngAfterContentChecked(): void {
    this.cdref.detectChanges();
  }

  private getFilteredResults(): void {
    this.myControl.valueChanges
      .pipe(
        debounceTime(200),
        takeUntil(this.unsubscribed$),
        tap(() => (this.isLoading = true)),
        switchMap((value) => {
          const query = new GroupQuery();

          const lnQuery = new GroupNameQuery();
          lnQuery.setMethod(TextQueryMethod.TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE);
          lnQuery.setName(value);

          query.setNameQuery(lnQuery);

          return from(this.groupService.listGroups(USER_LIMIT, 0, [query]));
        }),
      )
      .subscribe((userresp: ListGroupsResponse.AsObject | unknown) => {
        this.isLoading = false;
        if (userresp) {
          const results = (userresp as ListGroupsResponse.AsObject).resultList;
          this.filteredGroups = results.filter((filteredUser) => {
            return !this.groups.map((u) => u.id).includes(filteredUser.id);
          });
        }
      });
  }

  public displayFn(group?: Group.AsObject): string {
    return group ? `${group.name}` : '';
  }

  public add(event: MatChipInputEvent): void {
    if (!this.matAutocomplete.isOpen) {
      const input = event.chipInput?.inputElement;
      const value = event.value;

      if ((value || '').trim()) {
        const index = this.filteredGroups.findIndex((group) => {
          if (group.name) {
            return group.name === value;
          } else {
            return false;
          }
        });
        if (index > -1) {
          if (this.groups && this.groups.length > 0) {
            this.groups.push(this.filteredGroups[index]);
            this.selectionChanged.emit(this.groups);
          } else {
            this.groups = [this.filteredGroups[index]];
            this.selectionChanged.emit(this.groups);
          }
        }
      }

      if (input) {
        input.value = '';
      }
    }
  }

  public remove(user: Group.AsObject): void {
    const index = this.groups.indexOf(user);

    if (index >= 0) {
      this.groups.splice(index, 1);
      this.selectionChanged.emit(this.groups);
    }
  }

  public selected(event: MatAutocompleteSelectedEvent): void {
    const index = this.filteredGroups.findIndex((user) => user === event.option.value);
    if (index !== -1) {
      if (this.groups && this.groups.length > 0) {
        this.groups.push(this.filteredGroups[index]);
      } else {
        this.groups = [this.filteredGroups[index]];
      }

      this.selectionChanged.emit(this.groups);

      this.usernameInput.nativeElement.value = '';
      this.myControl.setValue(null);
    }
  }
}
