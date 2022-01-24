import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { MatCheckboxChange } from '@angular/material/checkbox';
import { SearchQuery as MemberSearchQuery } from 'src/app/proto/generated/zitadel/member_pb';
import { TextQueryMethod } from 'src/app/proto/generated/zitadel/object_pb';
import { ProjectQuery } from 'src/app/proto/generated/zitadel/project_pb';
import {
  DisplayNameQuery,
  EmailQuery,
  SearchQuery as UserSearchQuery,
  UserGrantQuery,
} from 'src/app/proto/generated/zitadel/user_pb';

enum FilterType {
  USER,
  USERGRANT,
  PROJECT,
}

enum SubQuery {
  DISPLAYNAME,
  EMAIL,
}

type FilterSearchQuery = UserSearchQuery | MemberSearchQuery | UserGrantQuery | ProjectQuery;

@Component({
  selector: 'cnsl-filter',
  templateUrl: './filter.component.html',
  styleUrls: ['./filter.component.scss'],
})
export class FilterComponent implements OnInit {
  @Input() type: FilterType = FilterType.USER;
  public FilterType: any = FilterType;
  public SubQuery: any = SubQuery;
  public searchQuery: FilterSearchQuery = new UserSearchQuery();
  @Output() public closedCard: EventEmitter<void> = new EventEmitter();
  @Output() public filterChanged: EventEmitter<FilterSearchQuery | undefined> = new EventEmitter();
  public showFilter: boolean = false;
  constructor() {}

  ngOnInit(): void {
    this.reset();
  }

  public changeCheckbox(subquery: SubQuery, event: MatCheckboxChange) {
    if (event.checked) {
      switch (subquery) {
        case SubQuery.DISPLAYNAME:
          const dnq = new DisplayNameQuery();
          dnq.setMethod(TextQueryMethod.TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE);
          dnq.setDisplayName('');
          (this.searchQuery as UserSearchQuery).setDisplayNameQuery(dnq);
          break;
        case SubQuery.EMAIL:
          const eq = new EmailQuery();
          eq.setMethod(TextQueryMethod.TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE);
          eq.setEmailAddress('');
          (this.searchQuery as UserSearchQuery).setEmailQuery(eq);
          break;
      }
    } else {
      switch (subquery) {
        case SubQuery.DISPLAYNAME:
          (this.searchQuery as UserSearchQuery).setDisplayNameQuery(undefined);
          break;
        case SubQuery.EMAIL:
          (this.searchQuery as UserSearchQuery).setEmailQuery(undefined);
          break;
      }
    }
  }

  public setValue(subquery: SubQuery, query: any, event: any) {
    switch (subquery) {
      case SubQuery.DISPLAYNAME:
        (query as DisplayNameQuery).setDisplayName(event?.target?.value);
        this.filterChanged.emit(this.filterCount ? this.searchQuery : undefined);
        break;
      case SubQuery.EMAIL:
        (query as EmailQuery).setEmailAddress(event?.target?.value);
        this.filterChanged.emit(this.filterCount ? this.searchQuery : undefined);
        break;
    }
  }

  public reset() {
    switch (this.type) {
      case FilterType.USER:
        this.searchQuery = new UserSearchQuery();
        this.searchQuery.setTypeQuery();
        this.searchQuery.setDisplayNameQuery();
        this.searchQuery.setUserNameQuery();
        this.searchQuery.setEmailQuery();
        this.searchQuery.setStateQuery();

        break;
      case FilterType.USERGRANT:
        this.searchQuery = new UserGrantQuery();
        break;
      case FilterType.PROJECT:
        this.searchQuery = new ProjectQuery();
        break;
    }
  }

  public emitFilter(): void {
    this.filterChanged.emit(this.filterCount ? this.searchQuery : undefined);
    this.showFilter = false;
  }

  public get filterCount(): number {
    return this.searchQuery
      ? Object.entries(this.searchQuery.toObject()).filter(([key, value]) => value !== undefined).length
      : 0;
  }
}
