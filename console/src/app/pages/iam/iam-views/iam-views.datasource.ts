import { DataSource } from '@angular/cdk/collections';
import { BehaviorSubject, from, Observable, of } from 'rxjs';
import { catchError, finalize, map } from 'rxjs/operators';
import { View } from 'src/app/proto/generated/admin_pb';
import { AdminService } from 'src/app/services/admin.service';

/**
 * Data source for the ProjectMembers view. This class should
 * encapsulate all logic for fetching and manipulating the displayed data
 * (including sorting, pagination, and filtering).
 */
export class IamViewsDataSource extends DataSource<View.AsObject> {
    public viewsSubject: BehaviorSubject<View.AsObject[]> = new BehaviorSubject<View.AsObject[]>([]);
    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();

    constructor(private adminService: AdminService) {
        super();
    }

    public loadViews(): void {
        this.loadingSubject.next(true);
        from(this.adminService.GetViews()).pipe(
            map(resp => {
                return resp.toObject().viewsList;
            }),
            catchError(() => of([])),
            finalize(() => this.loadingSubject.next(false)),
        ).subscribe(views => {
            this.viewsSubject.next(views);
        });
    }


    /**
     * Connect this data source to the table. The table will only update when
     * the returned stream emits new items.
     * @returns A stream of the items to be rendered.
     */
    public connect(): Observable<View.AsObject[]> {
        return this.viewsSubject.asObservable();
    }

    /**
     *  Called when the table is being destroyed. Use this function, to clean up
     * any open connections or free any held resources that were set up during connect.
     */
    public disconnect(): void {
        this.viewsSubject.complete();
        this.loadingSubject.complete();
    }
}
