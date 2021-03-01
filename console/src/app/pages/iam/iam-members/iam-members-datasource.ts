import { DataSource } from '@angular/cdk/collections';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { BehaviorSubject, from, Observable, of } from 'rxjs';
import { catchError, finalize, map } from 'rxjs/operators';
import { IamMemberView } from 'src/app/proto/generated/admin_pb';
import { AdminService } from 'src/app/services/admin.service';

/**
 * Data source for the ProjectMembers view. This class should
 * encapsulate all logic for fetching and manipulating the displayed data
 * (including sorting, pagination, and filtering).
 */
export class IamMembersDataSource extends DataSource<IamMemberView.AsObject> {
    public totalResult: number = 0;
    public viewTimestamp!: Timestamp.AsObject;
    public membersSubject: BehaviorSubject<IamMemberView.AsObject[]> = new BehaviorSubject<IamMemberView.AsObject[]>([]);
    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();

    constructor(private adminService: AdminService) {
        super();
    }

    public loadMembers(
        pageIndex: number, pageSize: number): void {
        const offset = pageIndex * pageSize;

        this.loadingSubject.next(true);

        from(this.adminService.listIAMMembers(pageSize, offset)).pipe(
            map(resp => {
                this.totalResult = resp.metaData?.totalResult || 0;
                if (resp.metaData?.viewTimestamp) {
                    this.viewTimestamp = resp.metaData?.viewTimestamp;
                }
                return resp.resultList;
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
    public connect(): Observable<IamMemberView.AsObject[]> {
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
