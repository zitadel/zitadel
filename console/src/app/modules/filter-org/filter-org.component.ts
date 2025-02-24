import { Component, DestroyRef, OnInit } from '@angular/core';
import { MatCheckboxChange } from '@angular/material/checkbox';
import { ActivatedRoute, Router } from '@angular/router';
import { take } from 'rxjs';
import { TextQueryMethod } from 'src/app/proto/generated/zitadel/object_pb';
import { OrgDomainQuery, OrgNameQuery, OrgQuery, OrgState, OrgStateQuery } from 'src/app/proto/generated/zitadel/org_pb';
import { UserNameQuery } from 'src/app/proto/generated/zitadel/user_pb';

import { FilterComponent } from '../filter/filter.component';

enum SubQuery {
  NAME,
  STATE,
  DOMAIN,
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
    destroyRef: DestroyRef,
    protected override route: ActivatedRoute,
  ) {
    super(router, route, destroyRef);
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
          } else if (filter.domainQuery) {
            const orgQuery = new OrgQuery();
            const orgDomainQuery = new OrgDomainQuery();
            orgDomainQuery.setDomain(filter.domainQuery.domain);
            orgDomainQuery.setMethod(filter.domainQuery.method);
            orgQuery.setDomainQuery(orgDomainQuery);
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
        case SubQuery.DOMAIN:
          const dq = new OrgDomainQuery();
          dq.setMethod(TextQueryMethod.TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE);
          dq.setDomain('');
          const odq = new OrgQuery();
          odq.setDomainQuery(dq);
          this.searchQueries.push(odq);
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
        case SubQuery.DOMAIN:
          const index_pdn = this.searchQueries.findIndex((q) => (q as OrgQuery).toObject().domainQuery !== undefined);
          if (index_pdn > -1) {
            this.searchQueries.splice(index_pdn, 1);
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
      case SubQuery.DOMAIN:
        (query as OrgDomainQuery).setDomain(value);
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
      case SubQuery.DOMAIN:
        const pdn = this.searchQueries.find((q) => (q as OrgQuery).toObject().domainQuery !== undefined);
        if (pdn) {
          return (pdn as OrgQuery).getDomainQuery();
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
