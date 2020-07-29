import { animate, animateChild, keyframes, query, stagger, style, transition, trigger } from '@angular/animations';
import { Component, Input, OnInit } from '@angular/core';
import { BehaviorSubject, from, Observable, of } from 'rxjs';
import { catchError, scan, take, tap } from 'rxjs/operators';
import { Change, Changes } from 'src/app/proto/generated/management_pb';
import { AuthUserService } from 'src/app/services/auth-user.service';
import { MgmtUserService } from 'src/app/services/mgmt-user.service';

export enum ChangeType {
    MYUSER = 'myuser',
    USER = 'user',
    ORG = 'org',
    PROJECT = 'project',
}

@Component({
    selector: 'app-changes',
    templateUrl: './changes.component.html',
    styleUrls: ['./changes.component.scss'],
    animations: [
        trigger('cardAnimation', [
            transition('* => *', [
                query('@animate', stagger('50ms', animateChild()), { optional: true }),
            ]),
        ]),
        trigger('animate', [
            transition(':enter', [
                animate('.2s ease-in', keyframes([
                    style({ opacity: 0 }),
                    style({ opacity: .5, transform: 'scale(1.02)' }),
                    style({ opacity: 1, transform: 'scale(1)' }),
                ])),
            ]),
        ]),
    ],
})
export class ChangesComponent implements OnInit {
    @Input() public changeType: ChangeType = ChangeType.USER;
    @Input() public id: string = '';
    @Input() public sortDirectionAsc: boolean = true;
    public bottom: boolean = false;

    private _done: BehaviorSubject<any> = new BehaviorSubject(false);
    private _loading: BehaviorSubject<any> = new BehaviorSubject(false);
    private _data: BehaviorSubject<any> = new BehaviorSubject([]);

    loading: Observable<boolean> = this._loading.asObservable();
    public data!: Observable<Change.AsObject[]>;
    public changes!: Changes.AsObject;
    constructor(private mgmtUserService: MgmtUserService, private authUserService: AuthUserService) { }

    ngOnInit(): void {
        this.init();
    }

    public scrollHandler(e: any): void {
        if (e === 'bottom') {
            this.more();
        }
    }

    private init(): void {
        let first: Promise<Changes>;
        switch (this.changeType) {
            case ChangeType.MYUSER: first = this.authUserService.GetMyUserChanges(10, 0);
                break;
            case ChangeType.USER: first = this.mgmtUserService.UserChanges(this.id, 10, 0);
                break;
            case ChangeType.PROJECT: first = this.mgmtUserService.ProjectChanges(this.id, 20, 0);
                break;
            case ChangeType.ORG: first = this.mgmtUserService.OrgChanges(this.id, 10, 0);
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
            case ChangeType.MYUSER: more = this.authUserService.GetMyUserChanges(10, cursor);
                break;
            case ChangeType.USER: more = this.mgmtUserService.UserChanges(this.id, 10, cursor);
                break;
            case ChangeType.PROJECT: more = this.mgmtUserService.ProjectChanges(this.id, 10, cursor);
                break;
            case ChangeType.ORG: more = this.mgmtUserService.OrgChanges(this.id, 10, cursor);
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
                tap((res: Changes) => {
                    let values = res.toObject().changesList;
                    // If prepending, reverse the batch order
                    values = false ? values.reverse() : values;

                    // update source with new values, done loading
                    this._data.next(values);

                    this._loading.next(false);

                    // no more values, mark done
                    if (!values.length) {
                        this._done.next(true);
                    }
                }),
                catchError(err => {
                    this._loading.next(false);
                    this.bottom = true;
                    return of([]);
                }),
                take(1),
            ).subscribe();
        }
    }
}
