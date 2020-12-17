import { Component, Input, OnDestroy, OnInit } from '@angular/core';
import { BehaviorSubject, from, Observable, of, Subject } from 'rxjs';
import { catchError, debounceTime, scan, take, takeUntil, tap } from 'rxjs/operators';
import { Change, Changes } from 'src/app/proto/generated/management_pb';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ManagementService } from 'src/app/services/mgmt.service';

export enum ChangeType {
    MYUSER = 'myuser',
    USER = 'user',
    ORG = 'org',
    PROJECT = 'project',
    APP = 'app',
}

@Component({
    selector: 'app-changes',
    templateUrl: './changes.component.html',
    styleUrls: ['./changes.component.scss'],
})
export class ChangesComponent implements OnInit, OnDestroy {
    @Input() public changeType: ChangeType = ChangeType.USER;
    @Input() public id: string = '';
    @Input() public secId: string = '';
    @Input() public sortDirectionAsc: boolean = true;
    @Input() public refresh!: Observable<void>;
    public bottom: boolean = false;

    private _done: BehaviorSubject<any> = new BehaviorSubject(false);
    private _loading: BehaviorSubject<any> = new BehaviorSubject(false);
    private _data: BehaviorSubject<any> = new BehaviorSubject([]);

    loading: Observable<boolean> = this._loading.asObservable();
    public data!: Observable<Change.AsObject[]>;
    public changes!: Changes.AsObject;
    private destroyed$: Subject<void> = new Subject();
    constructor(private mgmtUserService: ManagementService, private authUserService: GrpcAuthService) {

    }

    ngOnInit(): void {
        this.init();
        if (this.refresh) {
            this.refresh.pipe(takeUntil(this.destroyed$), debounceTime(2000)).subscribe(() => {
                this.init();
            });
        }
    }

    ngOnDestroy(): void {
        this.destroyed$.next();
    }

    public scrollHandler(e: any): void {
        if (e === 'bottom') {
            this.more();
        }
    }

    public init(): void {
        let first: Promise<Changes>;
        switch (this.changeType) {
            case ChangeType.MYUSER: first = this.authUserService.GetMyUserChanges(20, 0);
                break;
            case ChangeType.USER: first = this.mgmtUserService.UserChanges(this.id, 20, 0);
                break;
            case ChangeType.PROJECT: first = this.mgmtUserService.ProjectChanges(this.id, 20, 0);
                break;
            case ChangeType.ORG: first = this.mgmtUserService.OrgChanges(this.id, 20, 0);
                break;
            case ChangeType.APP: first = this.mgmtUserService.ApplicationChanges(this.id, this.secId, 20, 0);
                break;
        }

        this.mapAndUpdate(first);

        // Create the observable array for consumption in components
        this.data = this._data.asObservable().pipe(
            scan((acc, val) => {
                return false ? val.concat(acc) : acc.concat(val);
            }));
    }

    private more(): void {
        const cursor = this.getCursor();

        let more: Promise<Changes>;

        switch (this.changeType) {
            case ChangeType.MYUSER: more = this.authUserService.GetMyUserChanges(20, cursor);
                break;
            case ChangeType.USER: more = this.mgmtUserService.UserChanges(this.id, 20, cursor);
                break;
            case ChangeType.PROJECT: more = this.mgmtUserService.ProjectChanges(this.id, 20, cursor);
                break;
            case ChangeType.ORG: more = this.mgmtUserService.OrgChanges(this.id, 20, cursor);
                break;
            case ChangeType.APP: more = this.mgmtUserService.ApplicationChanges(this.id, this.secId, 20, cursor);
                break;
        }

        this.mapAndUpdate(more);
    }

    // Determines the snapshot to paginate query
    private getCursor(): number {
        const current = this._data.value;
        if (current.length) {
            return !this.sortDirectionAsc ? current[0].sequence :
                current[current.length - 1].sequence;
        }
        return 0;
    }

    // Maps the snapshot to usable format the updates source
    private mapAndUpdate(col: Promise<Changes>): any {
        if (this._done.value || this._loading.value) { return; }

        // Map snapshot with doc ref (needed for cursor)
        if (!this.bottom) {
            // loading
            this._loading.next(true);

            return from(col).pipe(
                take(1),
                tap((res: Changes) => {
                    const values = res.toObject().changesList;
                    // update source with new values, done loading
                    this._data.next(values);

                    this._loading.next(false);

                    // no more values, mark done
                    if (!values.length) {
                        this._done.next(true);
                    }
                }),
                catchError(_ => {
                    this._loading.next(false);
                    this.bottom = true;
                    return of([]);
                }),
            ).subscribe();
        }
    }
}
