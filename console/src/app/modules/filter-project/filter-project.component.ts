import { Component } from '@angular/core';
import { MatCheckboxChange } from '@angular/material/checkbox';
import { TextQueryMethod } from 'src/app/proto/generated/zitadel/object_pb';
import { ProjectNameQuery, ProjectQuery, ProjectState } from 'src/app/proto/generated/zitadel/project_pb';
import { UserNameQuery } from 'src/app/proto/generated/zitadel/user_pb';

import { FilterComponent } from '../filter/filter.component';

enum SubQuery {
  NAME,
  // RESOURCEOWNER,
}

@Component({
  selector: 'cnsl-filter-project',
  templateUrl: './filter-project.component.html',
  styleUrls: ['./filter-project.component.scss'],
})
export class FilterProjectComponent extends FilterComponent {
  public SubQuery: any = SubQuery;
  public searchQueries: ProjectQuery[] = [];

  public states: ProjectState[] = [ProjectState.PROJECT_STATE_ACTIVE, ProjectState.PROJECT_STATE_INACTIVE];
  constructor() {
    super();
  }

  public changeCheckbox(subquery: SubQuery, event: MatCheckboxChange) {
    if (event.checked) {
      switch (subquery) {
        // case SubQuery.RESOURCEOWNER:
        //   const ronq = new ProjectResourceOwnerQuery();
        //   ronq.setResourceOwner('');

        //   const ro_sq = new ProjectQuery();
        //   ro_sq.setProjectResourceOwnerQuery(ronq);

        //   this.searchQueries.push(ro_sq);
        //   break;
        case SubQuery.NAME:
          const nq = new ProjectNameQuery();
          nq.setMethod(TextQueryMethod.TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE);
          nq.setName('');

          const e_sq = new ProjectQuery();
          e_sq.setNameQuery(nq);

          this.searchQueries.push(e_sq);
          break;
      }
    } else {
      switch (subquery) {
        // case SubQuery.RESOURCEOWNER:
        //   const index_s = this.searchQueries.findIndex(
        //     (q) => (q as ProjectQuery).toObject().projectResourceOwnerQuery !== undefined,
        //   );
        //   if (index_s > -1) {
        //     this.searchQueries.splice(index_s, 1);
        //   }
        //   break;
        case SubQuery.NAME:
          const index_dn = this.searchQueries.findIndex((q) => (q as ProjectQuery).toObject().nameQuery !== undefined);
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
      // case SubQuery.RESOURCEOWNER:
      //   (query as ProjectResourceOwnerQuery).setResourceOwner(value);
      //   this.filterChanged.emit(this.filterCount ? this.searchQueries : undefined);
      //   break;
      case SubQuery.NAME:
        (query as ProjectNameQuery).setName(value);
        this.filterChanged.emit(this.filterCount ? this.searchQueries : undefined);
        break;
    }

    this.filterCount$.next(this.filterCount);
  }

  public getSubFilter(subquery: SubQuery): any {
    switch (subquery) {
      // case SubQuery.RESOURCEOWNER:
      //   const s = this.searchQueries.find((q) => (q as ProjectQuery).toObject().projectResourceOwnerQuery !== undefined);
      //   if (s) {
      //     return (s as ProjectQuery).getProjectResourceOwnerQuery();
      //   } else {
      //     return undefined;
      //   }
      case SubQuery.NAME:
        const dn = this.searchQueries.find((q) => (q as ProjectQuery).toObject().nameQuery !== undefined);
        if (dn) {
          return (dn as ProjectQuery).getNameQuery();
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
