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
  UserNameQuery,
} from 'src/app/proto/generated/zitadel/user_pb';

enum FilterType {
  USER,
  USERGRANT,
  PROJECT,
}

enum SubQuery {
  DISPLAYNAME,
  EMAIL,
  USERNAME,
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
  public searchQueries: FilterSearchQuery[] = [];
  @Output() public closedCard: EventEmitter<void> = new EventEmitter();
  @Output() public filterChanged: EventEmitter<FilterSearchQuery[] | undefined> = new EventEmitter();
  @Output() public filterOpen: EventEmitter<boolean> = new EventEmitter<boolean>(false);
  public showFilter: boolean = false;
  public methods: TextQueryMethod[] = [
    TextQueryMethod.TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE,
    TextQueryMethod.TEXT_QUERY_METHOD_ENDS_WITH_IGNORE_CASE,
    TextQueryMethod.TEXT_QUERY_METHOD_EQUALS_IGNORE_CASE,
  ];
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

          const dn_sq = new UserSearchQuery();
          dn_sq.setDisplayNameQuery(dnq);

          this.searchQueries.push(dn_sq);
          break;
        case SubQuery.EMAIL:
          const eq = new EmailQuery();
          eq.setMethod(TextQueryMethod.TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE);
          eq.setEmailAddress('');

          const e_sq = new UserSearchQuery();
          e_sq.setEmailQuery(eq);

          this.searchQueries.push(e_sq);
          break;

        case SubQuery.USERNAME:
          const unq = new UserNameQuery();
          unq.setMethod(TextQueryMethod.TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE);
          unq.setUserName('');

          const un_sq = new UserSearchQuery();
          un_sq.setUserNameQuery(unq);

          this.searchQueries.push(un_sq);
          break;
      }
    } else {
      switch (subquery) {
        case SubQuery.DISPLAYNAME:
          const index_dn = this.searchQueries.findIndex(
            (q) => (q as UserSearchQuery).toObject().displayNameQuery !== undefined,
          );
          if (index_dn > -1) {
            this.searchQueries.splice(index_dn, 1);
          }
          break;
        case SubQuery.EMAIL:
          const index_e = this.searchQueries.findIndex((q) => (q as UserSearchQuery).toObject().emailQuery !== undefined);
          if (index_e > -1) {
            this.searchQueries.splice(index_e, 1);
          }
          break;
        case SubQuery.USERNAME:
          const index_un = this.searchQueries.findIndex(
            (q) => (q as UserSearchQuery).toObject().userNameQuery !== undefined,
          );
          if (index_un > -1) {
            this.searchQueries.splice(index_un, 1);
          }
          break;
      }
    }
  }

  public setValue(subquery: SubQuery, query: any, event: any) {
    switch (subquery) {
      case SubQuery.DISPLAYNAME:
        (query as DisplayNameQuery).setDisplayName(event?.target?.value);
        this.filterChanged.emit(this.filterCount ? this.searchQueries : undefined);
        break;
      case SubQuery.EMAIL:
        (query as EmailQuery).setEmailAddress(event?.target?.value);
        this.filterChanged.emit(this.filterCount ? this.searchQueries : undefined);
        break;
      case SubQuery.USERNAME:
        (query as UserNameQuery).setUserName(event?.target?.value);
        this.filterChanged.emit(this.filterCount ? this.searchQueries : undefined);
        break;
    }
  }

  public setMethod(query: any, event: any) {
    (query as UserNameQuery).setMethod(event.value);
    this.filterChanged.emit(this.filterCount ? this.searchQueries : undefined);
  }

  public reset() {
    this.searchQueries = [];
  }

  public toggleFilter(): void {
    this.showFilter = !this.showFilter;
    this.filterOpen.emit(this.showFilter);
  }

  public emitFilter(): void {
    this.filterChanged.emit(this.filterCount ? this.searchQueries : undefined);
    this.showFilter = false;
    this.filterOpen.emit(false);
  }

  public get filterCount(): number {
    return this.searchQueries.length;
  }

  public getSubFilter(subquery: SubQuery): any {
    switch (subquery) {
      case SubQuery.DISPLAYNAME:
        const dn = this.searchQueries.find((q) => (q as UserSearchQuery).toObject().displayNameQuery !== undefined);
        if (dn) {
          return (dn as UserSearchQuery).getDisplayNameQuery();
        } else {
          return undefined;
        }
      case SubQuery.EMAIL:
        const e = this.searchQueries.find((q) => (q as UserSearchQuery).toObject().emailQuery !== undefined);
        if (e) {
          return (e as UserSearchQuery).getEmailQuery();
        } else {
          return undefined;
        }
      case SubQuery.USERNAME:
        const un = this.searchQueries.find((q) => (q as UserSearchQuery).toObject().userNameQuery !== undefined);
        if (un) {
          return (un as UserSearchQuery).getUserNameQuery();
        } else {
          return undefined;
        }
    }
  }
}
