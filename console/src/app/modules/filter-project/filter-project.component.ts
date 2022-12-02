import { Component, OnInit } from '@angular/core';
import { MatCheckboxChange } from '@angular/material/checkbox';
import { ActivatedRoute, Router } from '@angular/router';
import { take } from 'rxjs';
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
export class FilterProjectComponent extends FilterComponent implements OnInit {
  public SubQuery: any = SubQuery;
  public searchQueries: ProjectQuery[] = [];

  public states: ProjectState[] = [ProjectState.PROJECT_STATE_ACTIVE, ProjectState.PROJECT_STATE_INACTIVE];
  constructor(router: Router, route: ActivatedRoute) {
    super(router, route);
  }

  ngOnInit(): void {
    this.route.queryParams.pipe(take(1)).subscribe((params) => {
      const { filter } = params;
      if (filter) {
        const stringifiedFilters = filter as string;
        const filters: ProjectQuery.AsObject[] = JSON.parse(stringifiedFilters) as ProjectQuery.AsObject[];

        const projectQueries = filters.map((filter) => {
          if (filter.nameQuery) {
            const nameQuery = new ProjectNameQuery();

            const projectQuery = new ProjectQuery();
            nameQuery.setName(filter.nameQuery.name);
            nameQuery.setMethod(filter.nameQuery.method);

            projectQuery.setNameQuery(nameQuery);
            return projectQuery;
          } else {
            return undefined;
          }
        });

        this.searchQueries = projectQueries.filter((q) => q !== undefined) as ProjectQuery[];
        this.filterChanged.emit(this.searchQueries ? this.searchQueries : []);
        // this.showFilter = true;
        // this.filterOpen.emit(true);
      }
    });
  }

  public changeCheckbox(subquery: SubQuery, event: MatCheckboxChange) {
    if (event.checked) {
      switch (subquery) {
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
      case SubQuery.NAME:
        (query as ProjectNameQuery).setName(value);
        this.filterChanged.emit(this.searchQueries ? this.searchQueries : []);
        break;
    }
  }

  public getSubFilter(subquery: SubQuery): any {
    switch (subquery) {
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

  public get filterCounter(): number {
    return this.searchQueries.length;
  }
}
