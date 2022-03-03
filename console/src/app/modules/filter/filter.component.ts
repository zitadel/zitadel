import { Component, EventEmitter, OnDestroy, Output } from '@angular/core';
import { BehaviorSubject } from 'rxjs';
import { SearchQuery as MemberSearchQuery } from 'src/app/proto/generated/zitadel/member_pb';
import { TextQueryMethod } from 'src/app/proto/generated/zitadel/object_pb';
import { OrgQuery } from 'src/app/proto/generated/zitadel/org_pb';
import { ProjectQuery } from 'src/app/proto/generated/zitadel/project_pb';
import { SearchQuery as UserSearchQuery, UserGrantQuery } from 'src/app/proto/generated/zitadel/user_pb';

import { ActionKeysType } from '../action-keys/action-keys.component';

type FilterSearchQuery = UserSearchQuery | MemberSearchQuery | UserGrantQuery | ProjectQuery | OrgQuery;

@Component({
  selector: 'cnsl-filter',
  templateUrl: './filter.component.html',
  styleUrls: ['./filter.component.scss'],
})
export class FilterComponent implements OnDestroy {
  @Output() public filterChanged: EventEmitter<FilterSearchQuery[] | undefined> = new EventEmitter();
  @Output() public filterOpen: EventEmitter<boolean> = new EventEmitter<boolean>(false);

  @Output() public resetted: EventEmitter<void> = new EventEmitter();
  @Output() public trigger: EventEmitter<void> = new EventEmitter();

  public filterCount$: BehaviorSubject<number> = new BehaviorSubject(0);

  public showFilter: boolean = false;
  public methods: TextQueryMethod[] = [
    TextQueryMethod.TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE,
    TextQueryMethod.TEXT_QUERY_METHOD_ENDS_WITH_IGNORE_CASE,
    TextQueryMethod.TEXT_QUERY_METHOD_EQUALS_IGNORE_CASE,
  ];
  ActionKeysType: any = ActionKeysType;

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
    this.filterCount$.complete();
  }
}
