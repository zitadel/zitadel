import { Component } from '@angular/core';
import { MatCheckboxChange } from '@angular/material/checkbox';
import { TextQueryMethod } from 'src/app/proto/generated/zitadel/object_pb';
import { OrgNameQuery, OrgQuery, OrgState } from 'src/app/proto/generated/zitadel/org_pb';
import { UserNameQuery } from 'src/app/proto/generated/zitadel/user_pb';

import { FilterComponent } from '../filter/filter.component';

enum SubQuery {
  NAME,
}

@Component({
  selector: 'cnsl-filter-org',
  templateUrl: './filter-org.component.html',
  styleUrls: ['./filter-org.component.scss'],
})
export class FilterOrgComponent extends FilterComponent {
  public SubQuery: any = SubQuery;
  public searchQueries: OrgQuery[] = [];

  public states: OrgState[] = [OrgState.ORG_STATE_ACTIVE, OrgState.ORG_STATE_INACTIVE];
  constructor() {
    super();
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
      }
    } else {
      switch (subquery) {
        case SubQuery.NAME:
          const index_dn = this.searchQueries.findIndex((q) => (q as OrgQuery).toObject().nameQuery !== undefined);
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
      case SubQuery.NAME:
        (query as OrgNameQuery).setName(value);
        this.filterChanged.emit(this.filterCount ? this.searchQueries : undefined);
        break;
    }

    this.filterCount$.next(this.filterCount);
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
    }
  }

  public setMethod(query: any, event: any) {
    (query as UserNameQuery).setMethod(event.value);
    this.filterChanged.emit(this.filterCount ? this.searchQueries : undefined);
  }

  public emitFilter(): void {
    this.filterChanged.emit(this.filterCount ? this.searchQueries : undefined);
    this.showFilter = false;
    this.filterOpen.emit(false);

    this.filterCount$.next(this.filterCount);
  }

  public resetFilter(): void {
    this.searchQueries = [];
    this.emitFilter();
  }

  public get filterCount(): number {
    return this.searchQueries.length;
  }
}
