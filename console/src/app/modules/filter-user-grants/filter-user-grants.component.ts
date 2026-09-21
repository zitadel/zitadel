import { Component, DestroyRef, inject, OnInit } from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { MatCheckboxChange } from '@angular/material/checkbox';
import { ActivatedRoute, Router } from '@angular/router';
import { from, of, Subject, take } from 'rxjs';
import { catchError, debounceTime, distinctUntilChanged, switchMap } from 'rxjs/operators';
import { SuggestionField, UserSuggestionService } from 'src/app/services/user-suggestion.service';
import { TextQueryMethod } from 'src/app/proto/generated/zitadel/object_pb';
import {
  DisplayNameQuery,
  UserGrantOrgNameQuery,
  UserGrantProjectNameQuery,
  UserGrantQuery,
  UserGrantRoleKeyQuery,
  UserGrantWithGrantedQuery,
  UserNameQuery,
} from 'src/app/proto/generated/zitadel/user_pb';

import { FilterComponent } from '../filter/filter.component';

export enum SubQuery {
  DISPLAYNAME,
  USERNAME,
  ORGNAME,
  PROJECTNAME,
  ROLEKEY,
  WITHGRANTED,
}

@Component({
  selector: 'cnsl-filter-user-grants',
  templateUrl: './filter-user-grants.component.html',
  styleUrls: ['./filter-user-grants.component.scss'],
  standalone: false,
})
export class FilterUserGrantsComponent extends FilterComponent implements OnInit {
  public SubQuery: any = SubQuery;
  public searchQueries: UserGrantQuery[] = [];

  private readonly suggestionService = inject(UserSuggestionService);

  /**
   * Value suggestions for the user related text filters. Role keys are deliberately left
   * out: suggesting them needs a project, and this list spans every project at once.
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
        switchMap(({ subquery, value }) => {
          const field = FilterUserGrantsComponent.suggestionField(subquery);
          if (!field) {
            return of([] as string[]);
          }
          // a failed lookup must not break typing, it just means no suggestions
          return from(this.suggestionService.suggest(field, value)).pipe(catchError(() => of([] as string[])));
        }),
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
    this.setValue(subquery, query, { target: { value } });
  }

  private static suggestionField(subquery: SubQuery): SuggestionField | undefined {
    switch (subquery) {
      case SubQuery.DISPLAYNAME:
        return 'displayName';
      case SubQuery.USERNAME:
        return 'userName';
      default:
        return undefined;
    }
  }

  ngOnInit(): void {
    this.route.queryParams.pipe(take(1)).subscribe((params) => {
      const { filter } = params;
      if (filter) {
        const stringifiedFilters = filter as string;
        const filters: UserGrantQuery.AsObject[] = JSON.parse(stringifiedFilters) as UserGrantQuery.AsObject[];

        const userQueries = filters.map((filter) => {
          if (filter.userNameQuery) {
            const userGrantQuery = new UserGrantQuery();

            const userNameQuery = new UserNameQuery();
            userNameQuery.setUserName(filter.userNameQuery.userName);
            userNameQuery.setMethod(filter.userNameQuery.method);

            userGrantQuery.setUserNameQuery(userNameQuery);
            return userGrantQuery;
          } else if (filter.displayNameQuery) {
            const userGrantQuery = new UserGrantQuery();

            const displayNameQuery = new DisplayNameQuery();
            displayNameQuery.setDisplayName(filter.displayNameQuery.displayName);
            displayNameQuery.setMethod(filter.displayNameQuery.method);

            userGrantQuery.setDisplayNameQuery(displayNameQuery);
            return userGrantQuery;
          } else if (filter.orgNameQuery) {
            const userGrantQuery = new UserGrantQuery();

            const orgNameQuery = new UserGrantOrgNameQuery();
            orgNameQuery.setOrgName(filter.orgNameQuery.orgName);
            orgNameQuery.setMethod(filter.orgNameQuery.method);

            userGrantQuery.setOrgNameQuery(orgNameQuery);
            return userGrantQuery;
          } else if (filter.projectNameQuery) {
            const userGrantQuery = new UserGrantQuery();

            const projectNameQuery = new UserGrantProjectNameQuery();
            projectNameQuery.setProjectName(filter.projectNameQuery.projectName);
            projectNameQuery.setMethod(filter.projectNameQuery.method);

            userGrantQuery.setProjectNameQuery(projectNameQuery);
            return userGrantQuery;
          } else if (filter.roleKeyQuery) {
            const userGrantQuery = new UserGrantQuery();

            const roleKeyQuery = new UserGrantRoleKeyQuery();
            roleKeyQuery.setRoleKey(filter.roleKeyQuery.roleKey);
            roleKeyQuery.setMethod(filter.roleKeyQuery.method);

            userGrantQuery.setRoleKeyQuery(roleKeyQuery);
            return userGrantQuery;
          } else if (filter.withGrantedQuery) {
            const userGrantQuery = new UserGrantQuery();

            const withGrantedQuery = new UserGrantWithGrantedQuery();
            withGrantedQuery.setWithGranted(filter.withGrantedQuery.withGranted);

            userGrantQuery.setWithGrantedQuery(withGrantedQuery);
            return userGrantQuery;
          } else {
            return undefined;
          }
        });

        this.searchQueries = userQueries.filter((q) => q !== undefined) as UserGrantQuery[];
        this.emitQueries();
        // this.showFilter = true;
        // this.filterOpen.emit(true);
      }
    });
  }

  public changeCheckbox(subquery: SubQuery, event: MatCheckboxChange) {
    if (event.checked) {
      switch (subquery) {
        case SubQuery.DISPLAYNAME:
          const dnq = new DisplayNameQuery();
          dnq.setMethod(TextQueryMethod.TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE);
          dnq.setDisplayName('');

          const dn_sq = new UserGrantQuery();
          dn_sq.setDisplayNameQuery(dnq);

          this.searchQueries.push(dn_sq);
          break;

        case SubQuery.USERNAME:
          const unq = new UserNameQuery();
          unq.setMethod(TextQueryMethod.TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE);
          unq.setUserName('');

          const un_sq = new UserGrantQuery();
          un_sq.setUserNameQuery(unq);

          this.searchQueries.push(un_sq);
          break;

        case SubQuery.ORGNAME:
          const onq = new UserGrantOrgNameQuery();
          onq.setMethod(TextQueryMethod.TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE);
          onq.setOrgName('');

          const on_sq = new UserGrantQuery();
          on_sq.setOrgNameQuery(onq);

          this.searchQueries.push(on_sq);
          break;

        case SubQuery.PROJECTNAME:
          const pnq = new UserGrantProjectNameQuery();
          pnq.setMethod(TextQueryMethod.TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE);
          pnq.setProjectName('');

          const pn_sq = new UserGrantQuery();
          pn_sq.setProjectNameQuery(pnq);

          this.searchQueries.push(pn_sq);
          break;

        case SubQuery.ROLEKEY:
          const rkq = new UserGrantRoleKeyQuery();
          rkq.setMethod(TextQueryMethod.TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE);
          rkq.setRoleKey('');

          const rk_sq = new UserGrantQuery();
          rk_sq.setRoleKeyQuery(rkq);

          this.searchQueries.push(rk_sq);
          break;

        case SubQuery.WITHGRANTED:
          // A plain toggle: the API only knows "include grants of granted projects",
          // there is no value to type and nothing to compare against.
          const wgq = new UserGrantWithGrantedQuery();
          wgq.setWithGranted(true);

          const wg_sq = new UserGrantQuery();
          wg_sq.setWithGrantedQuery(wgq);

          this.searchQueries.push(wg_sq);
          break;
      }
    } else {
      switch (subquery) {
        case SubQuery.DISPLAYNAME:
          const index_dn = this.searchQueries.findIndex((q) => q.toObject().displayNameQuery !== undefined);
          if (index_dn > -1) {
            this.searchQueries.splice(index_dn, 1);
          }
          break;
        case SubQuery.USERNAME:
          const index_un = this.searchQueries.findIndex((q) => q.toObject().userNameQuery !== undefined);
          if (index_un > -1) {
            this.searchQueries.splice(index_un, 1);
          }
          break;
        case SubQuery.ORGNAME:
          const index_on = this.searchQueries.findIndex((q) => q.toObject().orgNameQuery !== undefined);
          if (index_on > -1) {
            this.searchQueries.splice(index_on, 1);
          }
          break;
        case SubQuery.PROJECTNAME:
          const index_pn = this.searchQueries.findIndex((q) => q.toObject().projectNameQuery !== undefined);
          if (index_pn > -1) {
            this.searchQueries.splice(index_pn, 1);
          }
          break;
        case SubQuery.ROLEKEY:
          const index_rk = this.searchQueries.findIndex((q) => q.toObject().roleKeyQuery !== undefined);
          if (index_rk > -1) {
            this.searchQueries.splice(index_rk, 1);
          }
          break;
        case SubQuery.WITHGRANTED:
          const index_wg = this.searchQueries.findIndex((q) => q.toObject().withGrantedQuery !== undefined);
          if (index_wg > -1) {
            this.searchQueries.splice(index_wg, 1);
          }
          break;
      }
    }
  }

  public setValue(subquery: SubQuery, query: any, event: any) {
    switch (subquery) {
      case SubQuery.DISPLAYNAME:
        (query as DisplayNameQuery).setDisplayName(event?.target?.value);
        this.emitQueries();
        break;
      case SubQuery.USERNAME:
        (query as UserNameQuery).setUserName(event?.target?.value);
        this.emitQueries();
        break;
      case SubQuery.ORGNAME:
        (query as UserGrantOrgNameQuery).setOrgName(event?.target?.value);
        this.emitQueries();
        break;
      case SubQuery.PROJECTNAME:
        (query as UserGrantProjectNameQuery).setProjectName(event?.target?.value);
        this.emitQueries();
        break;
      case SubQuery.ROLEKEY:
        (query as UserGrantRoleKeyQuery).setRoleKey(event?.target?.value);
        this.emitQueries();
        break;
    }
  }

  public getSubFilter(subquery: SubQuery): any {
    switch (subquery) {
      case SubQuery.DISPLAYNAME:
        const dn = this.searchQueries.find((q) => q.toObject().displayNameQuery !== undefined);
        if (dn) {
          return dn.getDisplayNameQuery();
        } else {
          return undefined;
        }

      case SubQuery.USERNAME:
        const un = this.searchQueries.find((q) => q.toObject().userNameQuery !== undefined);
        if (un) {
          return un.getUserNameQuery();
        } else {
          return undefined;
        }
      case SubQuery.ORGNAME:
        const e = this.searchQueries.find((q) => q.toObject().orgNameQuery !== undefined);
        if (e) {
          return e.getOrgNameQuery();
        } else {
          return undefined;
        }
      case SubQuery.PROJECTNAME:
        const pn = this.searchQueries.find((q) => q.toObject().projectNameQuery !== undefined);
        if (pn) {
          return pn.getProjectNameQuery();
        } else {
          return undefined;
        }
      case SubQuery.ROLEKEY:
        const rk = this.searchQueries.find((q) => q.toObject().roleKeyQuery !== undefined);
        if (rk) {
          return rk.getRoleKeyQuery();
        } else {
          return undefined;
        }
      case SubQuery.WITHGRANTED:
        const wg = this.searchQueries.find((q) => q.toObject().withGrantedQuery !== undefined);
        if (wg) {
          return wg.getWithGrantedQuery();
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

  /**
   * Emits only queries that carry a value. Ticking a checkbox creates the query with an
   * empty string, which the API rejects because text queries need at least one character.
   */
  private emitQueries(): void {
    this.filterChanged.emit(this.activeQueries);
  }

  private get activeQueries(): UserGrantQuery[] {
    return this.searchQueries.filter((query) => FilterUserGrantsComponent.hasValue(query));
  }

  /** Badge count, so a checkbox without a value does not look like an active filter. */
  public get activeQueryCount(): number {
    return this.activeQueries.length;
  }

  private static hasValue(query: UserGrantQuery): boolean {
    const q = query.toObject();
    if (q.displayNameQuery) {
      return !!q.displayNameQuery.displayName.trim();
    }
    if (q.userNameQuery) {
      return !!q.userNameQuery.userName.trim();
    }
    if (q.orgNameQuery) {
      return !!q.orgNameQuery.orgName.trim();
    }
    if (q.projectNameQuery) {
      return !!q.projectNameQuery.projectName.trim();
    }
    if (q.roleKeyQuery) {
      return !!q.roleKeyQuery.roleKey.trim();
    }
    // withGrantedQuery carries a boolean, it is complete as soon as it exists
    return true;
  }

  public resetFilter(): void {
    this.searchQueries = [];
    this.emitFilter();
  }
}
