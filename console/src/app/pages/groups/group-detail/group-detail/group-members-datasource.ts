import { DataSource } from '@angular/cdk/collections';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { BehaviorSubject, from, Observable, of } from 'rxjs';
import { catchError, finalize, map } from 'rxjs/operators';
import { ListGroupMembersResponse } from 'src/app/proto/generated/zitadel/management_pb';
import { Member } from 'src/app/proto/generated/zitadel/member_pb';
import { ManagementService } from 'src/app/services/mgmt.service';

/**
 * Data source for the ProjectMembers view. This class should
 * encapsulate all logic for fetching and manipulating the displayed data
 * (including sorting, pagination, and filtering).
 */
export class GroupMembersDataSource extends DataSource<Member.AsObject> {
  public totalResult: number = 0;
  public viewTimestamp!: Timestamp.AsObject;

  public membersSubject: BehaviorSubject<Member.AsObject[]> = new BehaviorSubject<Member.AsObject[]>([]);
  private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public loading$: Observable<boolean> = this.loadingSubject.asObservable();

  constructor(private service: ManagementService) {
    super();
  }

  public loadMembers(
    groupId: string,
    pageIndex: number,
    pageSize: number,
  ): void {
    const offset = pageIndex * pageSize;

    this.loadingSubject.next(true);

    const promise:
      | Promise<ListGroupMembersResponse.AsObject>
      | undefined =
      this.service.listGroupMembers(groupId, pageSize, offset);
    if (promise) {
      from(promise)
        .pipe(
          map((resp) => {
            if (resp.details?.totalResult) {
              this.totalResult = resp.details?.totalResult;
            } else {
              this.totalResult = 0;
            }
            if (resp.details?.viewTimestamp) {
              this.viewTimestamp = resp.details.viewTimestamp;
            }
            return resp.resultList.map((member) => {
              return {
                ...member,
                rolesList: [],
                userResourceOwner: "",
              };
            });
          }),
          catchError(() => of([])),
          finalize(() => this.loadingSubject.next(false)),
        )
        .subscribe((members) => {
          this.membersSubject.next(members);
        });
    }
  }

  /**
   * Connect this data source to the table. The table will only update when
   * the returned stream emits new items.
   * @returns A stream of the items to be rendered.
   */
  public connect(): Observable<Member.AsObject[]> {
    return this.membersSubject.asObservable();
  }

  /**
   *  Called when the table is being destroyed. Use this function, to clean up
   * any open connections or free any held resources that were set up during connect.
   */
  public disconnect(): void {
    this.membersSubject.complete();
    this.loadingSubject.complete();
  }
}
