import { ConnectedPosition, ConnectionPositionPair } from '@angular/cdk/overlay';
import { Component, EventEmitter, Input, OnDestroy, Output } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { Observable, Subject, takeUntil } from 'rxjs';
import { SearchQuery as MemberSearchQuery } from 'src/app/proto/generated/zitadel/member_pb';
import { TextQueryMethod } from 'src/app/proto/generated/zitadel/object_pb';
import { OrgQuery } from 'src/app/proto/generated/zitadel/org_pb';
import { ProjectQuery } from 'src/app/proto/generated/zitadel/project_pb';
import { SearchQuery as UserSearchQuery, UserGrantQuery } from 'src/app/proto/generated/zitadel/user_pb';

import { ActionKeysType } from '../action-keys/action-keys.component';

type FilterSearchQuery = UserSearchQuery | MemberSearchQuery | UserGrantQuery | ProjectQuery | OrgQuery;
type FilterSearchQueryAsObject =
  | UserSearchQuery.AsObject
  | MemberSearchQuery.AsObject
  | UserGrantQuery.AsObject
  | ProjectQuery.AsObject
  | OrgQuery.AsObject;

@Component({
  selector: 'cnsl-filter',
  templateUrl: './filter.component.html',
  styleUrls: ['./filter.component.scss'],
})
export class FilterComponent implements OnDestroy {
  @Output() public filterChanged: EventEmitter<FilterSearchQuery[]> = new EventEmitter();
  @Output() public filterOpen: EventEmitter<boolean> = new EventEmitter<boolean>(false);

  @Output() public resetted: EventEmitter<void> = new EventEmitter();
  @Output() public trigger: EventEmitter<void> = new EventEmitter();

  @Input() public queryCount: number = 0;

  private destroy$: Subject<void> = new Subject();
  public filterChanged$: Observable<FilterSearchQuery[]> = this.filterChanged.asObservable();

  public showFilter: boolean = false;
  public methods: TextQueryMethod[] = [
    TextQueryMethod.TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE,
    TextQueryMethod.TEXT_QUERY_METHOD_ENDS_WITH_IGNORE_CASE,
    TextQueryMethod.TEXT_QUERY_METHOD_EQUALS_IGNORE_CASE,
  ];
  ActionKeysType: any = ActionKeysType;

  public positions: ConnectedPosition[] = [
    new ConnectionPositionPair({ originX: 'start', originY: 'bottom' }, { overlayX: 'start', overlayY: 'top' }, 0, 10),
    new ConnectionPositionPair({ originX: 'end', originY: 'bottom' }, { overlayX: 'end', overlayY: 'top' }, 0, 10),
  ];

  public toggleFilter(): void {
    this.showFilter = !this.showFilter;
    this.filterOpen.emit(this.showFilter);
  }

  public emitFilter(): void {
    this.showFilter = false;
    this.filterOpen.emit(false);
    this.trigger.emit();
  }

  public ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  constructor(private router: Router, protected route: ActivatedRoute) {
    const changes$ = this.filterChanged.asObservable();
    changes$.pipe(takeUntil(this.destroy$)).subscribe((queries) => {
      const filters: Array<FilterSearchQueryAsObject | {}> | undefined = queries
        ?.map((q) => q.toObject())
        .map((query) =>
          Object.keys(query).reduce((acc, key) => {
            const _acc = acc;
            if ((query as any)[key] !== undefined) (_acc as any)[key] = (query as any)[key];
            return _acc as FilterSearchQueryAsObject;
          }, {}),
        );

      if (filters && Object.keys(filters)) {
        this.router.navigate([], {
          relativeTo: this.route,
          queryParams: {
            ['filter']: JSON.stringify(filters),
          },
          replaceUrl: true,
          queryParamsHandling: 'merge',
          skipLocationChange: false,
        });
      }
    });
  }
}
