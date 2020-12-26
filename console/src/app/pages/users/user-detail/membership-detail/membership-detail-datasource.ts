import { DataSource } from '@angular/cdk/collections';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { BehaviorSubject, from, Observable, of } from 'rxjs';
import { catchError, finalize, map } from 'rxjs/operators';
import { UserMembershipView } from 'src/app/proto/generated/zitadel/management_pb';
import { ManagementService } from 'src/app/services/mgmt.service';

export class MembershipDetailDataSource extends DataSource<UserMembershipView.AsObject> {
    public totalResult: number = 0;
    public viewTimestamp!: Timestamp.AsObject;
    public membersSubject: BehaviorSubject<UserMembershipView.AsObject[]>
        = new BehaviorSubject<UserMembershipView.AsObject[]>([]);
    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();

    constructor(private mgmtUserService: ManagementService) {
        super();
    }

    public loadMemberships(userId: string, pageIndex: number, pageSize: number): void {
        const offset = pageIndex * pageSize;

        this.loadingSubject.next(true);
        from(this.mgmtUserService.SearchUserMemberships(userId, pageSize, offset)).pipe(
            map(resp => {
                const response = resp.toObject();
                this.totalResult = response.totalResult;
                if (response.viewTimestamp) {
                    this.viewTimestamp = response.viewTimestamp;
                }
                return response.resultList;
            }),
            catchError(() => of([])),
            finalize(() => this.loadingSubject.next(false)),
        ).subscribe(members => {
            this.membersSubject.next(members);
        });
    }


    /**
     * Connect this data source to the table. The table will only update when
     * the returned stream emits new items.
     * @returns A stream of the items to be rendered.
     */
    public connect(): Observable<UserMembershipView.AsObject[]> {
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
