import { DataSource } from '@angular/cdk/collections';
import { BehaviorSubject, from, Observable, of } from 'rxjs';
import { catchError, finalize, map } from 'rxjs/operators';
import { ProjectRole } from 'src/app/proto/generated/management_pb';
import { ProjectService } from 'src/app/services/project.service';

/**
 * Data source for the ProjectMembers view. This class should
 * encapsulate all logic for fetching and manipulating the displayed data
 * (including sorting, pagination, and filtering).
 */
export class ProjectRolesDataSource extends DataSource<ProjectRole.AsObject> {
    public totalResult: number = 0;
    public rolesSubject: BehaviorSubject<ProjectRole.AsObject[]> = new BehaviorSubject<ProjectRole.AsObject[]>([]);
    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();

    constructor(private projectService: ProjectService) {
        super();
    }

    public loadRoles(projectId: string, pageIndex: number, pageSize: number, sortDirection?: string): void {
        const offset = pageIndex * pageSize;

        this.loadingSubject.next(true);
        from(this.projectService.SearchProjectRoles(projectId, pageSize, offset)).pipe(
            map(resp => {
                this.totalResult = resp.toObject().totalResult;
                return resp.toObject().resultList;
            }),
            catchError(() => of([])),
            finalize(() => this.loadingSubject.next(false)),
        ).subscribe(roles => {
            console.log(roles);
            this.rolesSubject.next(roles);
        });
    }


    /**
     * Connect this data source to the table. The table will only update when
     * the returned stream emits new items.
     * @returns A stream of the items to be rendered.
     */
    public connect(): Observable<ProjectRole.AsObject[]> {
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
