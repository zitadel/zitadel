import { Component, DestroyRef, inject, OnInit } from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { MatCheckboxChange } from '@angular/material/checkbox';
import { ActivatedRoute, Router } from '@angular/router';
import { from, of, Subject, take } from 'rxjs';
import { catchError, debounceTime, distinctUntilChanged, switchMap } from 'rxjs/operators';
import { ManagementService } from 'src/app/services/mgmt.service';
import { User } from 'src/app/proto/generated/zitadel/user_pb';
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
import { filter, map } from 'rxjs/operators';

/** Kept small: the list is a hint while typing, not a browsable result set. */
const SUGGESTION_LIMIT = 10;

export enum SubQuery {
  STATE,
  DISPLAYNAME,
  EMAIL,
  USERNAME,
}

@Component({
  selector: 'cnsl-filter-user',
  templateUrl: './filter-user.component.html',
  styleUrls: ['./filter-user.component.scss'],
  standalone: false,
})
export class FilterUserComponent extends FilterComponent implements OnInit {
  public SubQuery: any = SubQuery;
  private searchQueries: UserSearchQuery[] = [];

  public states: UserState[] = [
    UserState.USER_STATE_ACTIVE,
    UserState.USER_STATE_INACTIVE,
    UserState.USER_STATE_DELETED,
    UserState.USER_STATE_LOCKED,
    UserState.USER_STATE_INITIAL,
  ];
  private readonly mgmtService = inject(ManagementService);

  /**
   * Value suggestions for the text filters. Typing an exact match by hand is close to
   * impossible with the "equals" method, so the field offers existing values while
   * still accepting free text for the "contains" and "ends with" methods.
   */
  protected suggestions: string[] = [];
  private suggestionSubQuery: SubQuery | undefined;
  private readonly suggest$ = new Subject<{ subquery: SubQuery; value: string }>();

  constructor(router: Router, route: ActivatedRoute, destroyRef: DestroyRef) {
    super(router, route, destroyRef);

    this.suggest$
      .pipe(
        debounceTime(250),
        distinctUntilChanged((a, b) => a.subquery === b.subquery && a.value === b.value),
        switchMap(({ subquery, value }) =>
          from(this.fetchSuggestions(subquery, value)).pipe(
            // a failed lookup must not break typing, it just means no suggestions
            catchError(() => of([] as string[])),
          ),
        ),
        takeUntilDestroyed(destroyRef),
      )
      .subscribe((values) => (this.suggestions = values));
  }

  protected onSuggestionInput(subquery: SubQuery, event: Event): void {
    this.suggestionSubQuery = subquery;
    this.suggest$.next({ subquery, value: (event.target as HTMLInputElement).value });
  }

  /** Scoped per field so an open second filter never shows the first one's values. */
  protected suggestionsFor(subquery: SubQuery): string[] {
    return this.suggestionSubQuery === subquery ? this.suggestions : [];
  }

  protected selectSuggestion(subquery: SubQuery, query: any, value: string): void {
    this.setValue(subquery, query, { value });
  }

  private async fetchSuggestions(subquery: SubQuery, value: string): Promise<string[]> {
    const term = value.trim();
    if (!term) {
      return [];
    }

    const query = new UserSearchQuery();
    switch (subquery) {
      case SubQuery.DISPLAYNAME:
        const dnq = new DisplayNameQuery();
        dnq.setMethod(TextQueryMethod.TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE);
        dnq.setDisplayName(term);
        query.setDisplayNameQuery(dnq);
        break;
      case SubQuery.EMAIL:
        const eq = new EmailQuery();
        eq.setMethod(TextQueryMethod.TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE);
        eq.setEmailAddress(term);
        query.setEmailQuery(eq);
        break;
      case SubQuery.USERNAME:
        const unq = new UserNameQuery();
        unq.setMethod(TextQueryMethod.TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE);
        unq.setUserName(term);
        query.setUserNameQuery(unq);
        break;
      default:
        return [];
    }

    const response = await this.mgmtService.listUsers(SUGGESTION_LIMIT, 0, [query]);
    const values = response.resultList
      .map((user) => FilterUserComponent.suggestionValue(subquery, user))
      .filter((v): v is string => !!v);

    // Several users can share an address, so the raw list would repeat entries.
    return Array.from(new Set(values));
  }

  private static suggestionValue(subquery: SubQuery, user: User.AsObject): string | undefined {
    switch (subquery) {
      case SubQuery.DISPLAYNAME:
        return user.human?.profile?.displayName ?? user.machine?.name;
      case SubQuery.EMAIL:
        return user.human?.email?.email;
      case SubQuery.USERNAME:
        return user.userName;
      default:
        return undefined;
    }
  }

  ngOnInit(): void {
    this.route.queryParamMap
      .pipe(
        take(1),
        map((params) => params.get('filter')),
        filter(Boolean),
      )
      .subscribe((stringifiedFilters) => {
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
        this.emitQueries();
        // this.showFilter = true;
        // this.filterOpen.emit(true);
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
        this.emitQueries();
        break;
      case SubQuery.DISPLAYNAME:
        (query as DisplayNameQuery).setDisplayName(value);
        this.emitQueries();
        break;
      case SubQuery.EMAIL:
        (query as EmailQuery).setEmailAddress(value);
        this.emitQueries();
        break;
      case SubQuery.USERNAME:
        (query as UserNameQuery).setUserName(value);
        this.emitQueries();
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
    this.emitQueries();
  }

  public override emitFilter(): void {
    this.emitQueries();
    this.showFilter = false;
    this.filterOpen.emit(false);
  }

  public resetFilter(): void {
    this.searchQueries = [];
    this.emitFilter();
  }

  /**
   * Emits only queries that carry a value.
   *
   * Ticking a checkbox creates the query with an empty string, but the API rejects text
   * queries shorter than one character. Sending those produced an "invalid argument"
   * toast and, worse, persisted a broken filter in the URL that failed again on reload.
   * The incomplete query stays in the local list so its input keeps rendering.
   */
  private emitQueries(): void {
    this.filterChanged.emit(this.searchQueries.filter((query) => FilterUserComponent.hasValue(query)));
  }

  private static hasValue(query: UserSearchQuery): boolean {
    const q = query.toObject();
    if (q.displayNameQuery) {
      return !!q.displayNameQuery.displayName.trim();
    }
    if (q.emailQuery) {
      return !!q.emailQuery.emailAddress.trim();
    }
    if (q.userNameQuery) {
      return !!q.userNameQuery.userName.trim();
    }
    // state queries always carry a valid enum value
    return true;
  }
}
