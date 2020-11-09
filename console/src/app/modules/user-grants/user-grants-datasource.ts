import { DataSource } from '@angular/cdk/collections';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { BehaviorSubject, from, Observable, of } from 'rxjs';
import { catchError, finalize, map } from 'rxjs/operators';
import {
    SearchMethod,
    UserGrant,
    UserGrantSearchKey,
    UserGrantSearchQuery,
    UserGrantSearchResponse,
    UserGrantView,
} from 'src/app/proto/generated/management_pb';
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

    public grantsSubject: BehaviorSubject<UserGrantView.AsObject[]> = new BehaviorSubject<UserGrantView.AsObject[]>([]);
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
        queries?: UserGrantSearchQuery[],
    ): void {
        switch (context) {
            case UserGrantContext.USER:
                if (data && data.userId) {
                    this.loadingSubject.next(true);
                    const userfilter = new UserGrantSearchQuery();
                    userfilter.setKey(UserGrantSearchKey.USERGRANTSEARCHKEY_USER_ID);
                    userfilter.setMethod(SearchMethod.SEARCHMETHOD_EQUALS);
                    userfilter.setValue(data.userId);
                    if (queries) {
                        queries.push(userfilter);
                    } else {
                        queries = [userfilter];
                    }

                    const promise = this.userService.SearchUserGrants(pageSize, pageSize * pageIndex, queries);
                    this.loadResponse(promise);
                }
                break;
            case UserGrantContext.OWNED_PROJECT:
                if (data && data.projectId) {
                    this.loadingSubject.next(true);
                    const projectfilter = new UserGrantSearchQuery();
                    projectfilter.setKey(UserGrantSearchKey.USERGRANTSEARCHKEY_PROJECT_ID);
                    projectfilter.setMethod(SearchMethod.SEARCHMETHOD_EQUALS);
                    projectfilter.setValue(data.projectId);
                    if (queries) {
                        queries.push(projectfilter);
                    } else {
                        queries = [projectfilter];
                    }

                    const promise1 = this.userService.SearchUserGrants(pageSize, pageSize * pageIndex, queries);
                    this.loadResponse(promise1);
                }
                break;
            case UserGrantContext.GRANTED_PROJECT:
                if (data && data.grantId && data.projectId) {
                    this.loadingSubject.next(true);

                    const grantquery: UserGrantSearchQuery = new UserGrantSearchQuery();
                    grantquery.setKey(UserGrantSearchKey.USERGRANTSEARCHKEY_GRANT_ID);
                    grantquery.setMethod(SearchMethod.SEARCHMETHOD_EQUALS);
                    grantquery.setValue(data.grantId);

                    const projectfilter = new UserGrantSearchQuery();
                    projectfilter.setKey(UserGrantSearchKey.USERGRANTSEARCHKEY_PROJECT_ID);
                    projectfilter.setValue(data.projectId);

                    if (queries) {
                        queries.push(projectfilter);
                        queries.push(grantquery);
                    } else {
                        queries = [projectfilter, grantquery];
                    }

                    const promise2 = this.userService.SearchUserGrants(pageSize, pageSize * pageIndex, queries);
                    this.loadResponse(promise2);
                }
                break;
            default:
                this.loadingSubject.next(true);
                const promise3 = this.userService.SearchUserGrants(pageSize, pageSize * pageIndex, queries ?? []);
                this.loadResponse(promise3);
                break;
        }
    }

    private loadResponse(promise: Promise<UserGrantSearchResponse>): void {
        from(promise).pipe(
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
        ).subscribe(grants => {
            this.grantsSubject.next(grants);
        });
    }


    /**
     * Connect this data source to the table. The table will only update when
     * the returned stream emits new items.
     * @returns A stream of the items to be rendered.
     */
    public connect(): Observable<UserGrantView.AsObject[]> {
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
