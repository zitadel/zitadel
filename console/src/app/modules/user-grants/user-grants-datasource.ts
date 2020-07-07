import { DataSource } from '@angular/cdk/collections';
import { BehaviorSubject, from, Observable, of } from 'rxjs';
import { catchError, finalize, map } from 'rxjs/operators';
import {
    UserGrant,
    UserGrantSearchKey,
    UserGrantSearchQuery,
    UserGrantSearchResponse,
} from 'src/app/proto/generated/management_pb';
import { MgmtUserService } from 'src/app/services/mgmt-user.service';

export enum UserGrantContext {
    // AUTHUSER = 'authuser',
    USER = 'user',
    OWNED_PROJECT = 'owned',
    GRANTED_PROJECT = 'granted',
}

export class UserGrantsDataSource extends DataSource<UserGrant.AsObject> {
    public totalResult: number = 0;
    public grantsSubject: BehaviorSubject<UserGrant.AsObject[]> = new BehaviorSubject<UserGrant.AsObject[]>([]);
    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();

    constructor(private userService: MgmtUserService) {
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
        const offset = pageIndex * pageSize;

        switch (context) {
            case UserGrantContext.USER:
                if (data && data.userId) {
                    this.loadingSubject.next(true);

                    const userfilter = new UserGrantSearchQuery();
                    userfilter.setKey(UserGrantSearchKey.USERGRANTSEARCHKEY_USER_ID);
                    userfilter.setValue(data.userId);
                    if (queries) {
                        queries.push(userfilter);
                    } else {
                        queries = [userfilter];
                    }

                    const promise = this.userService.SearchUserGrants(10, 0, queries);
                    this.loadResponse(promise);
                }
                break;
            case UserGrantContext.OWNED_PROJECT:
                if (data && data.projectId) {
                    this.loadingSubject.next(true);

                    const promise1 = this.userService.SearchProjectUserGrants(data.projectId, 10, 0, queries);
                    this.loadResponse(promise1);
                }
                break;
            case UserGrantContext.GRANTED_PROJECT:
                if (data && data.grantId) {
                    this.loadingSubject.next(true);

                    const promise2 = this.userService.SearchProjectGrantUserGrants(data.grantId, 10, 0, queries);
                    this.loadResponse(promise2);
                }
                break;
        }
    }

    private loadResponse(promise: Promise<UserGrantSearchResponse>): void {
        from(promise).pipe(
            map(resp => {
                this.totalResult = resp.toObject().totalResult;
                console.log(resp.toObject().resultList);
                return resp.toObject().resultList;
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
