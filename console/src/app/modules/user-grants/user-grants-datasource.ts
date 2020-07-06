import { DataSource } from '@angular/cdk/collections';
import { BehaviorSubject, from, Observable, of } from 'rxjs';
import { catchError, finalize, map } from 'rxjs/operators';
import { UserGrant, UserGrantSearchKey, UserGrantSearchQuery } from 'src/app/proto/generated/management_pb';
import { MgmtUserService } from 'src/app/services/mgmt-user.service';

export class UserGrantsDataSource extends DataSource<UserGrant.AsObject> {
    public totalResult: number = 0;
    public grantsSubject: BehaviorSubject<UserGrant.AsObject[]> = new BehaviorSubject<UserGrant.AsObject[]>([]);
    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();

    constructor(private userService: MgmtUserService) {
        super();
    }

    public loadGrants(filter: UserGrantSearchKey, userId: string, pageIndex: number, pageSize: number): void {
        const offset = pageIndex * pageSize;

        this.loadingSubject.next(true);

        const query = new UserGrantSearchQuery();
        query.setKey(filter);
        query.setValue(userId);

        const queries: UserGrantSearchQuery[] = [query];
        from(this.userService.SearchUserGrants(10, 0, queries)).pipe(
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
