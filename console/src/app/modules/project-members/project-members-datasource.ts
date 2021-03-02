import { DataSource } from '@angular/cdk/collections';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { BehaviorSubject, from, Observable, of } from 'rxjs';
import { catchError, finalize, map } from 'rxjs/operators';
import { ListProjectGrantMembersResponse, ListProjectMembersResponse } from 'src/app/proto/generated/zitadel/management_pb';
import { Member } from 'src/app/proto/generated/zitadel/member_pb';
import { ManagementService } from 'src/app/services/mgmt.service';

import { ProjectType } from './project-members.component';

/**
 * Data source for the ProjectMembers view. This class should
 * encapsulate all logic for fetching and manipulating the displayed data
 * (including sorting, pagination, and filtering).
 */
export class ProjectMembersDataSource extends DataSource<Member.AsObject> {
    public totalResult: number = 0;
    public viewTimestamp!: Timestamp.AsObject;

    public membersSubject: BehaviorSubject<Member.AsObject[]> = new BehaviorSubject<Member.AsObject[]>([]);
    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();

    constructor(private mgmtService: ManagementService) {
        super();
    }

    public loadMembers(projectId: string,
        projectType: ProjectType,
        pageIndex: number, pageSize: number, grantId?: string): void {
        const offset = pageIndex * pageSize;

        this.loadingSubject.next(true);

        const promise: Promise<ListProjectMembersResponse.AsObject> | Promise<ListProjectGrantMembersResponse.AsObject> | undefined =
            projectType === ProjectType.PROJECTTYPE_OWNED ?
                this.mgmtService.listProjectMembers(projectId, pageSize, offset) :
                projectType === ProjectType.PROJECTTYPE_GRANTED && grantId ?
                    this.mgmtService.listProjectGrantMembers(projectId,
                        grantId, pageSize, offset) : undefined;
        if (promise) {
            from(promise).pipe(
                map(resp => {
                    if (resp.metaData?.totalResult) {
                        this.totalResult = resp.metaData?.totalResult;
                    }
                    if (resp.metaData?.viewTimestamp) {
                        this.viewTimestamp = resp.metaData.viewTimestamp;
                    }
                    return resp.resultList;
                }),
                catchError(() => of([])),
                finalize(() => this.loadingSubject.next(false)),
            ).subscribe(members => {
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
