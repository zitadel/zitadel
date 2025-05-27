import { DataSource } from '@angular/cdk/collections';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { BehaviorSubject, Observable } from 'rxjs';
import { ListGroupGrantResponse } from 'src/app/proto/generated/zitadel/management_pb';
import { GroupGrantQuery, GroupGrant, GroupGrantGroupIDQuery, GroupGrantProjectIDQuery, GroupGrantProjectGrantIDQuery } from 'src/app/proto/generated/zitadel/group_pb';
import { ManagementService } from 'src/app/services/mgmt.service';

type GroupGrantAsObject = GroupGrant.AsObject;

export enum GroupGrantContext {
  NONE = 'none',
  GROUP = 'group',
  OWNED_PROJECT = 'owned',
  GRANTED_PROJECT = 'granted',
}


export class GroupGrantsDataSource extends DataSource<GroupGrantAsObject> {
  public totalResult: number = 0;
  public viewTimestamp!: Timestamp.AsObject;

  public grantsSubject: BehaviorSubject<Array<GroupGrantAsObject>> = new BehaviorSubject<Array<GroupGrantAsObject>>([]);
  private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public loading$: Observable<boolean> = this.loadingSubject.asObservable();

  constructor(private groupService: ManagementService) {
    super();
  }

  public loadGrants(
    context: GroupGrantContext,
    pageIndex: number,
    pageSize: number,
    data: {
      projectId?: string;
      grantId?: string;
      groupId?: string;
    },
    queries?: GroupGrantQuery[],
  ): void {
    switch (context) {
      case GroupGrantContext.GROUP:
        if (data && data.groupId) {
          this.loadingSubject.next(true);

          const groupfilter = new GroupGrantQuery();
          const ugUiq = new GroupGrantGroupIDQuery();
          ugUiq.setGroupId(data.groupId);
          groupfilter.setGroupIdQuery(ugUiq);

          if (queries) {
            queries.push(groupfilter);
          } else {
            queries = [groupfilter];
          }

          const promise = this.groupService.listGroupGrants(pageSize, pageSize * pageIndex, queries);
          this.loadResponse(promise);
        }
        break;
      case GroupGrantContext.OWNED_PROJECT:
        if (data && data.projectId) {
          this.loadingSubject.next(true);

          const projectfilter = new GroupGrantQuery();
          const ugPfq = new GroupGrantProjectIDQuery();
          ugPfq.setProjectId(data.projectId);
          projectfilter.setProjectIdQuery(ugPfq);

          if (queries) {
            queries.push(projectfilter);
          } else {
            queries = [projectfilter];
          }

          const promise1 = this.groupService.listGroupGrants(pageSize, pageSize * pageIndex, queries);
          this.loadResponse(promise1);
        }
        break;
      case GroupGrantContext.GRANTED_PROJECT:
        if (data && data.grantId && data.projectId) {
          this.loadingSubject.next(true);

          const grantfilter = new GroupGrantQuery();

          const uggiq = new GroupGrantProjectGrantIDQuery();
          uggiq.setProjectGrantId(data.grantId);
          grantfilter.setProjectGrantIdQuery(uggiq);

          const projectfilter = new GroupGrantQuery();
          const ugPfq = new GroupGrantProjectIDQuery();
          ugPfq.setProjectId(data.projectId);
          projectfilter.setProjectIdQuery(ugPfq);

          if (queries) {
            queries.push(grantfilter);
          } else {
            queries = [grantfilter];
          }

          const promise2 = this.groupService.listGroupGrants(pageSize, pageSize * pageIndex, queries);
          this.loadResponse(promise2);
        }
        break;
      default:
        this.loadingSubject.next(true);
        const promise3 = this.groupService.listGroupGrants(pageSize, pageSize * pageIndex, queries ?? []);
        this.loadResponse(promise3);
        break;
    }
  }

  private loadResponse(promise: Promise<ListGroupGrantResponse.AsObject>): void {
    promise
      .then((resp) => {
        this.loadingSubject.next(false);
        if (resp.resultList) {
          this.grantsSubject.next(resp.resultList);
        }
        if (resp.details) {
          this.totalResult = resp.details.totalResult;
          if (resp.details.viewTimestamp) {
            this.viewTimestamp = resp.details.viewTimestamp;
          }
        }
      })
      .catch((error) => {
        this.grantsSubject.next([]);
        this.loadingSubject.next(false);
      });
  }

  /**
   * Connect this data source to the table. The table will only update when
   * the returned stream emits new lists of items.
   * @returns A stream of item lists to be rendered.
   */
  public connect(): Observable<Array<GroupGrantAsObject>> {
    return this.grantsSubject.asObservable();
  }

  /**
   *  Called when the table is being destroyed. Use this function, to clean up
   * any open connections or free any held resources that were set up during connect.
   */
  public disconnect(): void {
    this.grantsSubject.complete();
    this.loadingSubject.complete();
  }
}
