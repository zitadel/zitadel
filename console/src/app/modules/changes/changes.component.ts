import { KeyValue } from '@angular/common';
import { Component, DestroyRef, Input, OnDestroy, OnInit } from '@angular/core';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { BehaviorSubject, from, Observable, of, Subject } from 'rxjs';
import { catchError, debounceTime, scan, take, takeUntil, tap } from 'rxjs/operators';
import { ListMyUserChangesResponse } from 'src/app/proto/generated/zitadel/auth_pb';
import { Change } from 'src/app/proto/generated/zitadel/change_pb';
import {
  ListAppChangesResponse,
  ListOrgChangesResponse,
  ListProjectChangesResponse,
  ListUserChangesResponse,
} from 'src/app/proto/generated/zitadel/management_pb';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';

export enum ChangeType {
  MYUSER = 'myuser',
  USER = 'user',
  ORG = 'org',
  PROJECT = 'project',
  PROJECT_GRANT = 'project-grant',
  APP = 'app',
}

export interface MappedChange {
  key: string;
  values: Array<{
    data: any[];
    dates: Timestamp.AsObject[];
    editorId: string;
    editorName: string;
    editorDisplayName: string;
    editorAvatarUrl: string;
    editorPreferredLoginName: string;
    eventTypes: Array<{ key: string; localizedMessage: string }>;
    sequences: number[];
  }>;
}

type ListChanges =
  | ListMyUserChangesResponse.AsObject
  | ListUserChangesResponse.AsObject
  | ListProjectChangesResponse.AsObject
  | ListOrgChangesResponse.AsObject
  | ListAppChangesResponse.AsObject;

// todo: update this component to react to input changes
@Component({
  selector: 'cnsl-changes',
  templateUrl: './changes.component.html',
  styleUrls: ['./changes.component.scss'],
})
export class ChangesComponent implements OnInit {
  @Input({ required: true }) public changeType!: ChangeType;
  @Input() public id: string = '';
  @Input() public secId: string = '';
  @Input() public sortDirectionAsc: boolean = true;
  @Input() public refresh?: Observable<void>;
  public bottom: boolean = false;

  private _done: BehaviorSubject<any> = new BehaviorSubject(false);
  private _loading: BehaviorSubject<any> = new BehaviorSubject(false);
  private _data: BehaviorSubject<any> = new BehaviorSubject([]);

  loading: Observable<boolean> = this._loading.asObservable();
  public data: Observable<MappedChange[]> = this._data.asObservable().pipe(
    scan((acc, val) => {
      return acc.concat(val);
    }),
  );
  public changes!: ListChanges;
  constructor(
    private readonly mgmtUserService: ManagementService,
    private readonly authUserService: GrpcAuthService,
    private readonly destroyRef: DestroyRef,
  ) {}

  ngOnInit(): void {
    this.init();
    if (this.refresh) {
      this.refresh.pipe(takeUntilDestroyed(this.destroyRef), debounceTime(2000)).subscribe(() => {
        this._data = new BehaviorSubject([]);
        this.init();
      });
    }
  }

  public init(): void {
    let first: Promise<ListChanges>;
    switch (this.changeType) {
      case ChangeType.MYUSER:
        first = this.authUserService.listMyUserChanges(30, 0);
        break;
      case ChangeType.USER:
        first = this.mgmtUserService.listUserChanges(this.id, 30, 0);
        break;
      case ChangeType.PROJECT:
        first = this.mgmtUserService.listProjectChanges(this.id, 30, 0);
        break;
      case ChangeType.PROJECT_GRANT:
        first = this.mgmtUserService.listProjectGrantChanges(this.id, this.secId, 30, 0);
        break;
      case ChangeType.ORG:
        first = this.mgmtUserService.listOrgChanges(30, 0);
        break;
      case ChangeType.APP:
        first = this.mgmtUserService.listAppChanges(this.id, this.secId, 30, 0);
        break;
    }

    this.mapAndUpdate(first);
  }

  public more(): void {
    const cursor = this.getCursor();

    let more: Promise<ListChanges>;

    switch (this.changeType) {
      case ChangeType.MYUSER:
        more = this.authUserService.listMyUserChanges(20, cursor);
        break;
      case ChangeType.USER:
        more = this.mgmtUserService.listUserChanges(this.id, 20, cursor);
        break;
      case ChangeType.PROJECT:
        more = this.mgmtUserService.listProjectChanges(this.id, 20, cursor);
        break;
      case ChangeType.PROJECT_GRANT:
        more = this.mgmtUserService.listProjectGrantChanges(this.id, this.secId, 20, cursor);
        break;
      case ChangeType.ORG:
        more = this.mgmtUserService.listOrgChanges(20, cursor);
        break;
      case ChangeType.APP:
        more = this.mgmtUserService.listAppChanges(this.id, this.secId, 20, cursor);
        break;
    }

    this.mapAndUpdate(more);
  }

  // Determines the snapshot to paginate query
  private getCursor(): number {
    const current = this._data.value;

    if (current.length) {
      const lastElementValues = current[current.length - 1].values;
      const seq = lastElementValues[lastElementValues.length - 1].sequences;
      return seq[seq.length - 1];
    }
    return 0;
  }

  // Maps the snapshot to usable format the updates source
  private mapAndUpdate(col: Promise<ListChanges>): any {
    if (this._done.value || this._loading.value) {
      return;
    }

    // Map snapshot with doc ref (needed for cursor)
    if (!this.bottom) {
      // loading
      this._loading.next(true);

      return from(col)
        .pipe(
          take(1),
          tap((res: ListChanges) => {
            const values = res.resultList;
            const mapped = this.mapChanges(values);

            this._data.next(mapped);

            this._loading.next(false);

            if (!values.length) {
              this._done.next(true);
            }
          }),
          catchError((_) => {
            this._loading.next(false);
            this.bottom = true;
            return of([]);
          }),
        )
        .subscribe();
    }
  }

  private mapChanges(changes: Change.AsObject[]): {
    key: string;
    values: any[];
  }[] {
    const splitted: { [editorId: string]: any[] } = {};
    changes.forEach((change) => {
      if (change.changeDate) {
        const index = `${this.getDateString(change.changeDate)}`;
        // `${this.getDateString(change.changeDate)}:${change.editorId}`;

        if (index) {
          if (splitted[index]) {
            const userData: any = {
              editor: change.editorDisplayName,
              editorId: change.editorId,
              editorDisplayName: change.editorDisplayName,
              editorPreferredLoginName: change.editorPreferredLoginName,
              editorAvatarUrl: change.editorAvatarUrl,

              dates: [change.changeDate],
              // data: [change.data],
              eventTypes: [change.eventType],
              sequences: [change.sequence],
            };
            const lastIndex = splitted[index].length - 1;
            if (lastIndex > -1 && splitted[index][lastIndex].editor === change.editorDisplayName) {
              splitted[index][lastIndex].dates.push(change.changeDate);
              // splitted[index][lastIndex].data.push(change.data);
              splitted[index][lastIndex].eventTypes.push(change.eventType);
              splitted[index][lastIndex].sequences.push(change.sequence);
            } else {
              splitted[index].push(userData);
            }
          } else {
            splitted[index] = [
              {
                editor: change.editorDisplayName,
                editorId: change.editorId,
                editorDisplayName: change.editorDisplayName,
                editorPreferredLoginName: change.editorPreferredLoginName,
                editorAvatarUrl: change.editorAvatarUrl,

                dates: [change.changeDate],
                // data: [change.data],
                eventTypes: [change.eventType],
                sequences: [change.sequence],
              },
            ];
          }
        }
      }
    });
    const arr = Object.keys(splitted).map((key) => {
      return { key: key, values: splitted[key] };
    });

    arr.sort((a, b) => {
      return parseFloat(b.key) - parseFloat(a.key);
    });

    return arr;
  }

  getDateString(ts: Timestamp.AsObject): string {
    const date = new Date(ts.seconds * 1000 + ts.nanos / 1000 / 1000);
    return date.getUTCFullYear() + this.pad(date.getUTCMonth() + 1) + this.pad(date.getUTCDate());
  }

  getTimestampIndex(date: any): number {
    const ts: Date = new Date(date.seconds * 1000 + date.nanos / 1000 / 1000);
    return ts.getTime();
  }

  pad(n: number): string {
    return n < 10 ? '0' + n : n.toString();
  }

  // Order by ascending property value
  /* eslint-disable */
  valueAscOrder = (a: KeyValue<number, string>, b: KeyValue<number, string>): number => {
    return a.value.localeCompare(b.value);
  };

  // Order by descending property key
  keyDescOrder = (a: KeyValue<number, string>, b: KeyValue<number, string>): number => {
    return a.key > b.key ? -1 : b.key > a.key ? 1 : 0;
  };
  /* eslint-enable */
}
