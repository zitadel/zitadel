import { Component, OnInit } from '@angular/core';
import { MatCheckboxChange } from '@angular/material/checkbox';
import { ActivatedRoute, Router } from '@angular/router';
import { take } from 'rxjs';
import { TextQueryMethod } from 'src/app/proto/generated/zitadel/object_pb';
import {
    DisplayNameQuery,
    UserGrantOrgNameQuery,
    UserGrantProjectNameQuery,
    UserGrantQuery,
    UserNameQuery,
} from 'src/app/proto/generated/zitadel/user_pb';

import { FilterComponent } from '../filter/filter.component';

enum SubQuery {
  DISPLAYNAME,
  USERNAME,
  ORGNAME,
  PROJECTNAME,
}

@Component({
  selector: 'cnsl-filter-user-grants',
  templateUrl: './filter-user-grants.component.html',
  styleUrls: ['./filter-user-grants.component.scss'],
})
export class FilterUserGrantsComponent extends FilterComponent implements OnInit {
  public SubQuery: any = SubQuery;
  public searchQueries: UserGrantQuery[] = [];

  constructor(router: Router, route: ActivatedRoute) {
    super(router, route);
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

            userGrantQuery.setProjectNameQuery(projectNameQuery);
            return userGrantQuery;
          } else {
            return undefined;
          }
        });

        this.searchQueries = userQueries.filter((q) => q !== undefined) as UserGrantQuery[];
        this.filterChanged.emit(this.searchQueries ? this.searchQueries : []);
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
      }
    }
  }

  public setValue(subquery: SubQuery, query: any, event: any) {
    switch (subquery) {
      case SubQuery.DISPLAYNAME:
        (query as DisplayNameQuery).setDisplayName(event?.target?.value);
        this.filterChanged.emit(this.searchQueries ? this.searchQueries : []);
        break;
      case SubQuery.USERNAME:
        (query as UserNameQuery).setUserName(event?.target?.value);
        this.filterChanged.emit(this.searchQueries ? this.searchQueries : []);
        break;
      case SubQuery.ORGNAME:
        (query as UserGrantOrgNameQuery).setOrgName(event?.target?.value);
        this.filterChanged.emit(this.searchQueries ? this.searchQueries : []);
        break;
      case SubQuery.PROJECTNAME:
        (query as UserGrantProjectNameQuery).setProjectName(event?.target?.value);
        this.filterChanged.emit(this.searchQueries ? this.searchQueries : []);
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
