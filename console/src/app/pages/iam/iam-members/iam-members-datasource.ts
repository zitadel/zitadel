import { DataSource } from '@angular/cdk/collections';
import { BehaviorSubject, from, Observable, of } from 'rxjs';
import { catchError, finalize, map } from 'rxjs/operators';
import { IamMemberSearchResponse } from 'src/app/proto/generated/admin_pb';
import { ProjectMember } from 'src/app/proto/generated/management_pb';
import { AdminService } from 'src/app/services/admin.service';

/**
 * Data source for the ProjectMembers view. This class should
 * encapsulate all logic for fetching and manipulating the displayed data
 * (including sorting, pagination, and filtering).
 */
export class IamMembersDataSource extends DataSource<ProjectMember.AsObject> {
    public totalResult: number = 0;
    public membersSubject: BehaviorSubject<ProjectMember.AsObject[]> = new BehaviorSubject<ProjectMember.AsObject[]>([]);
    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();

    constructor(private adminService: AdminService) {
        super();
    }

    public loadMembers(
        pageIndex: number, pageSize: number, grantId?: string, sortDirection?: string): void {
        const offset = pageIndex * pageSize;

        this.loadingSubject.next(true);
        // TODO
        const promise: Promise<IamMemberSearchResponse> =
            this.adminService.SearchIamMembers(pageSize, offset);
        if (promise) {
            from(promise).pipe(
                map(resp => {
                    this.totalResult = resp.toObject().totalResult;
                    return resp.toObject().resultList;
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
    public connect(): Observable<ProjectMember.AsObject[]> {
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
