import { Component, OnInit } from '@angular/core';
import { MatCheckboxChange } from '@angular/material/checkbox';
import { ActivatedRoute, Router } from '@angular/router';
import { take } from 'rxjs';
import { TextQueryMethod } from 'src/app/proto/generated/zitadel/object_pb';
import { RoleDisplayNameQuery, RoleKeyQuery, RoleQuery } from 'src/app/proto/generated/zitadel/project_pb';
import { UserNameQuery } from 'src/app/proto/generated/zitadel/user_pb';

import { FilterComponent } from '../filter/filter.component';

enum SubQuery {
  KEY,
  DISPLAYNAME,
}

@Component({
  selector: 'cnsl-filter-project-roles',
  templateUrl: './filter-project-roles.component.html',
  styleUrls: ['./filter-project-roles.component.scss'],
})
export class FilterProjectRolesComponent extends FilterComponent implements OnInit {
  public SubQuery: any = SubQuery;
  public searchQueries: RoleQuery[] = [];

  constructor(router: Router, route: ActivatedRoute) {
    super(router, route);
  }

  ngOnInit(): void {
    this.route.queryParams.pipe(take(1)).subscribe((params) => {
      const { filter } = params;
      if (filter) {
        const stringifiedFilters = filter as string;
        const filters: RoleQuery.AsObject[] = JSON.parse(stringifiedFilters) as RoleQuery.AsObject[];

        const projectQueries = filters.map((filter) => {
          if (filter.keyQuery) {
            const keyQuery = new RoleKeyQuery();

            const projectQuery = new RoleQuery();
            keyQuery.setKey(filter.keyQuery.key);
            keyQuery.setMethod(filter.keyQuery.method);

            projectQuery.setKeyQuery(keyQuery);
            return projectQuery;
          } else if (filter.displayNameQuery) {
            const displayNameQuery = new RoleDisplayNameQuery();

            const projectQuery = new RoleQuery();
            displayNameQuery.setDisplayName(filter.displayNameQuery.displayName);
            displayNameQuery.setMethod(filter.displayNameQuery.method);

            projectQuery.setDisplayNameQuery(displayNameQuery);
            return projectQuery;
          } else {
            return undefined;
          }
        });

        this.searchQueries = projectQueries.filter((q) => q !== undefined) as RoleQuery[];
        this.filterChanged.emit(this.searchQueries ? this.searchQueries : []);
      }
    });
  }

  public changeCheckbox(subquery: SubQuery, event: MatCheckboxChange) {
    if (event.checked) {
      switch (subquery) {
        case SubQuery.KEY:
          const kq = new RoleKeyQuery();
          kq.setMethod(TextQueryMethod.TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE);
          kq.setKey('');

          const krq = new RoleQuery();
          krq.setKeyQuery(kq);

          this.searchQueries.push(krq);
          break;
        case SubQuery.DISPLAYNAME:
          const dq = new RoleDisplayNameQuery();
          dq.setMethod(TextQueryMethod.TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE);
          dq.setDisplayName('');

          const drq = new RoleQuery();
          drq.setDisplayNameQuery(dq);

          this.searchQueries.push(drq);
          break;
      }
    } else {
      switch (subquery) {
        case SubQuery.KEY:
          const index_kn = this.searchQueries.findIndex((q) => (q as RoleQuery).toObject().keyQuery !== undefined);
          if (index_kn > -1) {
            this.searchQueries.splice(index_kn, 1);
          }
          break;
        case SubQuery.DISPLAYNAME:
          const index_dn = this.searchQueries.findIndex((q) => (q as RoleQuery).toObject().displayNameQuery !== undefined);
          if (index_dn > -1) {
            this.searchQueries.splice(index_dn, 1);
          }
          break;
      }
    }
  }

  public setValue(subquery: SubQuery, query: any, event: any) {
    const value = event?.target?.value ?? event.value;
    switch (subquery) {
      case SubQuery.KEY:
        (query as RoleKeyQuery).setKey(value);
        this.filterChanged.emit(this.searchQueries ? this.searchQueries : []);
        break;
      case SubQuery.DISPLAYNAME:
        (query as RoleDisplayNameQuery).setDisplayName(value);
        this.filterChanged.emit(this.searchQueries ? this.searchQueries : []);
        break;
    }
  }

  public getSubFilter(subquery: SubQuery): any {
    switch (subquery) {
      case SubQuery.KEY:
        const ksf = this.searchQueries.find((q) => (q as RoleQuery).toObject().keyQuery !== undefined);
        if (ksf) {
          return (ksf as RoleQuery).getKeyQuery();
        } else {
          return undefined;
        }
      case SubQuery.DISPLAYNAME:
        const dsf = this.searchQueries.find((q) => (q as RoleQuery).toObject().displayNameQuery !== undefined);
        if (dsf) {
          return (dsf as RoleQuery).getDisplayNameQuery();
        } else {
          return undefined;
        }
    }
  }

  public setMethod(query: any, event: any) {
    (query as UserNameQuery).setMethod(event.value);
    this.filterChanged.emit(this.searchQueries ? this.searchQueries : []);
  }

  public override emitFilter(): void {
    this.filterChanged.emit(this.searchQueries ? this.searchQueries : []);
    this.showFilter = false;
    this.filterOpen.emit(false);
  }

  public resetFilter(): void {
    this.searchQueries = [];
    this.emitFilter();
  }

  public get filterCounter(): number {
    return this.searchQueries.length;
  }
}
