import { Component, OnInit } from '@angular/core';
import { MatCheckboxChange } from '@angular/material/checkbox';
import { ActivatedRoute, Router } from '@angular/router';
import { take } from 'rxjs';
import { TextQueryMethod } from 'src/app/proto/generated/zitadel/object_pb';
import { OrgNameQuery, OrgQuery, OrgState, OrgStateQuery } from 'src/app/proto/generated/zitadel/org_pb';
import { UserNameQuery } from 'src/app/proto/generated/zitadel/user_pb';

import { FilterComponent } from '../filter/filter.component';

enum SubQuery {
  NAME,
  STATE,
}

@Component({
  selector: 'cnsl-filter-org',
  templateUrl: './filter-org.component.html',
  styleUrls: ['./filter-org.component.scss'],
})
export class FilterOrgComponent extends FilterComponent implements OnInit {
  public SubQuery: any = SubQuery;
  public searchQueries: OrgQuery[] = [];

  public states: OrgState[] = [OrgState.ORG_STATE_ACTIVE, OrgState.ORG_STATE_INACTIVE, OrgState.ORG_STATE_REMOVED];

  constructor(
    router: Router,
    protected override route: ActivatedRoute,
  ) {
    super(router, route);
  }

  ngOnInit(): void {
    this.route.queryParams.pipe(take(1)).subscribe((params) => {
      const { filter } = params;
      if (filter) {
        const stringifiedFilters = filter as string;
        const filters: OrgQuery.AsObject[] = JSON.parse(stringifiedFilters) as OrgQuery.AsObject[];

        const orgQueries = filters.map((filter) => {
          if (filter.nameQuery) {
            const orgQuery = new OrgQuery();
            const orgNameQuery = new OrgNameQuery();
            orgNameQuery.setName(filter.nameQuery.name);
            orgNameQuery.setMethod(filter.nameQuery.method);
            orgQuery.setNameQuery(orgNameQuery);
            return orgQuery;
          } else if (filter.stateQuery) {
            const orgQuery = new OrgQuery();
            const orgStateQuery = new OrgStateQuery();
            orgStateQuery.setState(filter.stateQuery.state);
            orgQuery.setStateQuery(orgStateQuery);
            return orgQuery;
          } else {
            return undefined;
          }
        });

        this.searchQueries = orgQueries.filter((q) => q !== undefined) as OrgQuery[];
        this.filterChanged.emit(this.searchQueries ? this.searchQueries : undefined);
        // this.showFilter = true;
        // this.filterOpen.emit(true);
      }
    });
  }

  public changeCheckbox(subquery: SubQuery, event: MatCheckboxChange) {
    if (event.checked) {
      switch (subquery) {
        case SubQuery.NAME:
          const nq = new OrgNameQuery();
          nq.setMethod(TextQueryMethod.TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE);
          nq.setName('');
          const oq = new OrgQuery();
          oq.setNameQuery(nq);
          this.searchQueries.push(oq);
          break;
        case SubQuery.STATE:
          const sq = new OrgStateQuery();
          sq.setState(OrgState.ORG_STATE_ACTIVE);
          const osq = new OrgQuery();
          osq.setStateQuery(sq);
          this.searchQueries.push(osq);
          break;
      }
    } else {
      switch (subquery) {
        case SubQuery.NAME:
          const index_dn = this.searchQueries.findIndex((q) => (q as OrgQuery).toObject().nameQuery !== undefined);
          if (index_dn > -1) {
            this.searchQueries.splice(index_dn, 1);
          }
          break;
        case SubQuery.STATE:
          const index_sn = this.searchQueries.findIndex((q) => (q as OrgQuery).toObject().stateQuery !== undefined);
          if (index_sn > -1) {
            this.searchQueries.splice(index_sn, 1);
          }
          break;
      }
    }
  }

  public setValue(subquery: SubQuery, query: any, event: any) {
    const value = event?.target?.value ?? event.value;
    switch (subquery) {
      case SubQuery.NAME:
        (query as OrgNameQuery).setName(value);
        this.filterChanged.emit(this.searchQueries ? this.searchQueries : []);
        break;
      case SubQuery.STATE:
        (query as OrgStateQuery).setState(value);
        this.filterChanged.emit(this.searchQueries ? this.searchQueries : []);
        break;
    }
  }

  public getSubFilter(subquery: SubQuery): any {
    switch (subquery) {
      case SubQuery.NAME:
        const dn = this.searchQueries.find((q) => (q as OrgQuery).toObject().nameQuery !== undefined);
        if (dn) {
          return (dn as OrgQuery).getNameQuery();
        } else {
          return undefined;
        }
      case SubQuery.STATE:
        const sn = this.searchQueries.find((q) => (q as OrgQuery).toObject().stateQuery !== undefined);
        if (sn) {
          return (sn as OrgQuery).getStateQuery();
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
}
