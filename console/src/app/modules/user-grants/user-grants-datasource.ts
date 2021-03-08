import { DataSource } from '@angular/cdk/collections';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { BehaviorSubject, from, Observable, of } from 'rxjs';
import { catchError, finalize, map } from 'rxjs/operators';
import { ListUserGrantResponse } from 'src/app/proto/generated/zitadel/management_pb';
import {
    UserGrant,
    UserGrantProjectGrantIDQuery,
    UserGrantProjectIDQuery,
    UserGrantQuery,
    UserGrantUserIDQuery,
} from 'src/app/proto/generated/zitadel/user_pb';
import { ManagementService } from 'src/app/services/mgmt.service';

export enum UserGrantContext {
    NONE = 'none',
    USER = 'user',
    OWNED_PROJECT = 'owned',
    GRANTED_PROJECT = 'granted',
}

export class UserGrantsDataSource extends DataSource<UserGrant.AsObject> {
    public totalResult: number = 0;
    public viewTimestamp!: Timestamp.AsObject;

    public grantsSubject: BehaviorSubject<UserGrant.AsObject[]> = new BehaviorSubject<UserGrant.AsObject[]>([]);
    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();

    constructor(private userService: ManagementService) {
        super();
    }

    public loadGrants(
        context: UserGrantContext,
        pageIndex: number,
        pageSize: number,
        data: {
            projectId?: string;
            grantId?: string;
            userId?: string;
        },
        queries?: UserGrantQuery[],
    ): void {
        switch (context) {
            case UserGrantContext.USER:
                if (data && data.userId) {
                    this.loadingSubject.next(true);

                    const userfilter = new UserGrantQuery();
                    const ugUiq = new UserGrantUserIDQuery();
                    ugUiq.setUserId(data.userId);
                    userfilter.setUserIdQuery(ugUiq);

                    if (queries) {
                        queries.push(userfilter);
                    } else {
                        queries = [userfilter];
                    }

                    const promise = this.userService.listUserGrants(pageSize, pageSize * pageIndex, queries);
                    this.loadResponse(promise);
                }
                break;
            case UserGrantContext.OWNED_PROJECT:
                if (data && data.projectId) {
                    this.loadingSubject.next(true);

                    const projectfilter = new UserGrantQuery();
                    const ugPfq = new UserGrantProjectIDQuery();
                    ugPfq.setProjectId(data.projectId);
                    projectfilter.setProjectIdQuery(ugPfq);

                    if (queries) {
                        queries.push(projectfilter);
                    } else {
                        queries = [projectfilter];
                    }

                    const promise1 = this.userService.listUserGrants(pageSize, pageSize * pageIndex, queries);
                    this.loadResponse(promise1);
                }
                break;
            case UserGrantContext.GRANTED_PROJECT:
                if (data && data.grantId && data.projectId) {
                    this.loadingSubject.next(true);

                    const grantfilter = new UserGrantQuery();

                    const uggiq = new UserGrantProjectGrantIDQuery();
                    uggiq.setProjectGrantId(data.grantId);
                    grantfilter.setProjectGrantIdQuery(uggiq);

                    const projectfilter = new UserGrantQuery();
                    const ugPfq = new UserGrantProjectIDQuery();
                    ugPfq.setProjectId(data.projectId);
                    projectfilter.setProjectIdQuery(ugPfq);

                    if (queries) {
                        queries.push(grantfilter);
                    } else {
                        queries = [grantfilter];
                    }

                    const promise2 = this.userService.listUserGrants(pageSize, pageSize * pageIndex, queries);
                    this.loadResponse(promise2);
                }
                break;
            default:
                this.loadingSubject.next(true);
                const promise3 = this.userService.listUserGrants(pageSize, pageSize * pageIndex, queries ?? []);
                this.loadResponse(promise3);
                break;
        }
    }

    private loadResponse(promise: Promise<ListUserGrantResponse.AsObject>): void {
        from(promise).pipe(
            map(resp => {
                if (resp.details?.totalResult) {
                    this.totalResult = resp.details.totalResult;
                }
                if (resp.details?.viewTimestamp) {
                    this.viewTimestamp = resp.details.viewTimestamp;
                }
                return resp.resultList;
            }),
            catchError(() => of([])),
            finalize(() => this.loadingSubject.next(false)),
        ).subscribe(grants => {
            this.grantsSubject.next(grants);
        });
    }


    /**
     * Connect this data source to the table. The table will only update when
     * the returned stream emits new items.
     * @returns A stream of the items to be rendered.
     */
    public connect(): Observable<UserGrant.AsObject[]> {
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
