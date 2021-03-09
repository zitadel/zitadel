import { DataSource } from '@angular/cdk/collections';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { BehaviorSubject, from, Observable, of } from 'rxjs';
import { catchError, finalize, map } from 'rxjs/operators';
import { Role } from 'src/app/proto/generated/zitadel/project_pb';
import { ManagementService } from 'src/app/services/mgmt.service';

/**
 * Data source for the ProjectMembers view. This class should
 * encapsulate all logic for fetching and manipulating the displayed data
 * (including sorting, pagination, and filtering).
 */
export class ProjectRolesDataSource extends DataSource<Role.AsObject> {
    public totalResult: number = 0;
    public viewTimestamp!: Timestamp.AsObject;

    public rolesSubject: BehaviorSubject<Role.AsObject[]> = new BehaviorSubject<Role.AsObject[]>([]);
    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();

    constructor(private mgmtService: ManagementService) {
        super();
    }

    public loadRoles(projectId: string, pageIndex: number, pageSize: number, sortDirection?: string): void {
        const offset = pageIndex * pageSize;

        this.loadingSubject.next(true);
        from(this.mgmtService.listProjectRoles(projectId, pageSize, offset)).pipe(
            map(resp => {
                if (resp.details?.totalResult !== undefined) {
                    this.totalResult = resp.details.totalResult;
                }
                if (resp.details?.viewTimestamp) {
                    this.viewTimestamp = resp.details.viewTimestamp;
                }
                return resp.resultList;
            }),
            catchError(() => of([])),
            finalize(() => this.loadingSubject.next(false)),
        ).subscribe(roles => {
            this.rolesSubject.next(roles);
        });
    }


    /**
     * Connect this data source to the table. The table will only update when
     * the returned stream emits new items.
     * @returns A stream of the items to be rendered.
     */
    public connect(): Observable<Role.AsObject[]> {
        return this.rolesSubject.asObservable();
    }

    /**
     *  Called when the table is being destroyed. Use this function, to clean up
     * any open connections or free any held resources that were set up during connect.
     */
    public disconnect(): void {
        this.rolesSubject.complete();
        this.loadingSubject.complete();
    }
}
