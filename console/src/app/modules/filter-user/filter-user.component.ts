import { Component, OnInit } from '@angular/core';
import { MatCheckboxChange } from '@angular/material/checkbox';
import { ActivatedRoute, Router } from '@angular/router';
import { take } from 'rxjs';
import { TextQueryMethod } from 'src/app/proto/generated/zitadel/object_pb';
import {
    DisplayNameQuery,
    EmailQuery,
    SearchQuery as UserSearchQuery,
    StateQuery,
    UserNameQuery,
    UserState,
} from 'src/app/proto/generated/zitadel/user_pb';

import { FilterComponent } from '../filter/filter.component';

enum SubQuery {
  STATE,
  DISPLAYNAME,
  EMAIL,
  USERNAME,
}

@Component({
  selector: 'cnsl-filter-user',
  templateUrl: './filter-user.component.html',
  styleUrls: ['./filter-user.component.scss'],
})
export class FilterUserComponent extends FilterComponent implements OnInit {
  public SubQuery: any = SubQuery;
  public searchQueries: UserSearchQuery[] = [];

  public states: UserState[] = [
    UserState.USER_STATE_ACTIVE,
    UserState.USER_STATE_INACTIVE,
    UserState.USER_STATE_DELETED,
    UserState.USER_STATE_INITIAL,
    UserState.USER_STATE_LOCKED,
    UserState.USER_STATE_SUSPEND,
  ];
  constructor(router: Router, route: ActivatedRoute) {
    super(router, route);
  }

  ngOnInit(): void {
    this.route.queryParams.pipe(take(1)).subscribe((params) => {
      const { filter } = params;
      if (filter) {
        const stringifiedFilters = filter as string;
        const filters: UserSearchQuery.AsObject[] = JSON.parse(stringifiedFilters) as UserSearchQuery.AsObject[];

        const userQueries = filters.map((filter) => {
          if (filter.userNameQuery) {
            const userQuery = new UserSearchQuery();

            const userNameQuery = new UserNameQuery();
            userNameQuery.setUserName(filter.userNameQuery.userName);
            userNameQuery.setMethod(filter.userNameQuery.method);

            userQuery.setUserNameQuery(userNameQuery);
            return userQuery;
          } else if (filter.displayNameQuery) {
            const userQuery = new UserSearchQuery();

            const displayNameQuery = new DisplayNameQuery();
            displayNameQuery.setDisplayName(filter.displayNameQuery.displayName);
            displayNameQuery.setMethod(filter.displayNameQuery.method);

            userQuery.setDisplayNameQuery(displayNameQuery);
            return userQuery;
          } else if (filter.emailQuery) {
            const userQuery = new UserSearchQuery();

            const emailQuery = new EmailQuery();
            emailQuery.setEmailAddress(filter.emailQuery.emailAddress);
            emailQuery.setMethod(filter.emailQuery.method);

            userQuery.setEmailQuery(emailQuery);
            return userQuery;
          } else if (filter.stateQuery) {
            const userQuery = new UserSearchQuery();

            const stateQuery = new StateQuery();
            stateQuery.setState(filter.stateQuery.state);

            userQuery.setStateQuery(stateQuery);
            return userQuery;
          } else {
            return undefined;
          }
        });

        this.searchQueries = userQueries.filter((q) => q !== undefined) as UserSearchQuery[];
        this.filterChanged.emit(this.searchQueries ? this.searchQueries : []);
        // this.showFilter = true;
        // this.filterOpen.emit(true);
      }
    });
  }

  public changeCheckbox(subquery: SubQuery, event: MatCheckboxChange) {
    if (event.checked) {
      switch (subquery) {
        case SubQuery.STATE:
          const sq = new StateQuery();
          sq.setState(UserState.USER_STATE_ACTIVE);

          const s_sq = new UserSearchQuery();
          s_sq.setStateQuery(sq);

          this.searchQueries.push(s_sq);
          break;
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
        case SubQuery.STATE:
          const index_s = this.searchQueries.findIndex((q) => (q as UserSearchQuery).toObject().stateQuery !== undefined);
          if (index_s > -1) {
            this.searchQueries.splice(index_s, 1);
          }
          break;
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
    const value = event?.target?.value ?? event.value;
    switch (subquery) {
      case SubQuery.STATE:
        (query as StateQuery).setState(value);
        this.filterChanged.emit(this.searchQueries ? this.searchQueries : []);
        break;
      case SubQuery.DISPLAYNAME:
        (query as DisplayNameQuery).setDisplayName(value);
        this.filterChanged.emit(this.searchQueries ? this.searchQueries : []);
        break;
      case SubQuery.EMAIL:
        (query as EmailQuery).setEmailAddress(value);
        this.filterChanged.emit(this.searchQueries ? this.searchQueries : []);
        break;
      case SubQuery.USERNAME:
        (query as UserNameQuery).setUserName(value);
        this.filterChanged.emit(this.searchQueries ? this.searchQueries : []);
        break;
    }
  }

  public getSubFilter(subquery: SubQuery): any {
    switch (subquery) {
      case SubQuery.STATE:
        const s = this.searchQueries.find((q) => (q as UserSearchQuery).toObject().stateQuery !== undefined);
        if (s) {
          return (s as UserSearchQuery).getStateQuery();
        } else {
          return undefined;
        }
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

  public setMethod(query: any, event: any) {
    (query as UserNameQuery).setMethod(event.value);
    this.filterChanged.emit(this.searchQueries ? this.searchQueries : []);
  }

  public emitFilter(): void {
    this.filterChanged.emit(this.searchQueries ? this.searchQueries : []);
    this.showFilter = false;
    this.filterOpen.emit(false);
  }

  public resetFilter(): void {
    this.searchQueries = [];
    this.emitFilter();
  }
}
