import { DataSource } from '@angular/cdk/collections';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { BehaviorSubject, from, Observable, of } from 'rxjs';
import { catchError, finalize, map } from 'rxjs/operators';
import { Application } from 'src/app/proto/generated/zitadel/management_pb';
import { ManagementService } from 'src/app/services/mgmt.service';

/**
 * Data source for the ProjectMembers view. This class should
 * encapsulate all logic for fetching and manipulating the displayed data
 * (including sorting, pagination, and filtering).
 */
export class ProjectApplicationsDataSource extends DataSource<Application.AsObject> {
    public totalResult: number = 0;
    public viewTimestamp!: Timestamp.AsObject;

    public appsSubject: BehaviorSubject<Application.AsObject[]> = new BehaviorSubject<Application.AsObject[]>([]);
    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();

    constructor(private mgmtService: ManagementService) {
        super();
    }

    public loadApps(projectId: string, pageIndex: number, pageSize: number): void {
        const offset = pageIndex * pageSize;

        this.loadingSubject.next(true);
        from(this.mgmtService.SearchApplications(projectId, pageSize, offset)).pipe(
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
        ).subscribe(apps => {
            this.appsSubject.next(apps);
        });
    }


    /**
     * Connect this data source to the table. The table will only update when
     * the returned stream emits new items.
     * @returns A stream of the items to be rendered.
     */
    public connect(): Observable<Application.AsObject[]> {
        return this.appsSubject.asObservable();
    }

    /**
     *  Called when the table is being destroyed. Use this function, to clean up
     * any open connections or free any held resources that were set up during connect.
     */
    public disconnect(): void {
        this.appsSubject.complete();
        this.loadingSubject.complete();
    }
}
