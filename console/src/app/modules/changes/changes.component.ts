import { Component, Input, OnDestroy, OnInit } from '@angular/core';
import { BehaviorSubject, from, Observable, of, Subject } from 'rxjs';
import { catchError, debounceTime, scan, take, takeUntil, tap } from 'rxjs/operators';
import { Change, Changes } from 'src/app/proto/generated/management_pb';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { KeyValue } from '@angular/common';

export enum ChangeType {
    MYUSER = 'myuser',
    USER = 'user',
    ORG = 'org',
    PROJECT = 'project',
    APP = 'app',
}

export interface MappedChange {
    key: string,
    values: Array<{
        data: any[];
        dates: Timestamp.AsObject[];
        editorId: string;
        editorName: string;
        eventTypes: Array<{ key: string; localizedMessage: string; }>;
        sequences: number[];
    }>;
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
    public data!: Observable<MappedChange[]>;
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
        console.log('cursor' + cursor);

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
        console.log(current);

        if (current.length) {
            const lastElementValues = current[current.length - 1].values;
            const seq = lastElementValues[lastElementValues.length - 1].sequences;
            console.log(seq);
            return seq[seq.length - 1];
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
                    const mapped = this.mapChanges(values);
                    // update source with new values, done loading
                    // this._data.next(values);
                    this._data.next(mapped);

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

    mapChanges(changes: Change.AsObject[]) {
        const splitted: { [editorId: string]: any[]; } = {};
        changes.forEach((change) => {
            if (change.changeDate) {
                const index = this.getDateString(change.changeDate);//change.changeDate?.seconds;//this.getDateString(change.changeDate);

                if (index) {
                    if (splitted[index]) {
                        const userData: any = {
                            editor: change.editor,
                            editorId: change.editorId,

                            dates: [change.changeDate],
                            data: [change.data],
                            eventTypes: [change.eventType],
                            sequences: [change.sequence],
                        };
                        const lastIndex = splitted[index].length - 1;
                        if (lastIndex > -1 && splitted[index][lastIndex].editor === change.editor) {
                            splitted[index][lastIndex].dates.push(change.changeDate);
                            splitted[index][lastIndex].data.push(change.data);
                            splitted[index][lastIndex].eventTypes.push(change.eventType);
                            splitted[index][lastIndex].sequences.push(change.sequence);
                        } else {
                            splitted[index].push(userData);
                        }
                    } else {
                        splitted[index] = [
                            {
                                editorName: change.editor,
                                editorId: change.editorId,
                                dates: [change.changeDate],
                                data: [change.data],
                                eventTypes: [change.eventType],
                                sequences: [change.sequence],
                            }
                        ];
                    }
                }
            }
        });
        const arr = Object.keys(splitted).map(key => {
            return { key: key, values: splitted[key] };
        });

        arr.sort((a, b) => {
            return parseFloat(b.key) - parseFloat(a.key);
        });

        return arr;
    }

    getDateString(ts: Timestamp.AsObject) {
        const date = new Date(ts.seconds * 1000 + ts.nanos / 1000 / 1000);
        return date.getUTCFullYear() + this.pad(date.getUTCMonth() + 1) + this.pad(date.getUTCDate());
    }

    getTimestampIndex(date: any): number {
        const ts: Date = new Date(date.seconds * 1000 + date.nanos / 1000 / 1000);
        console.log(ts);
        return ts.getTime();
    }

    pad(n: number): string {
        return n < 10 ? '0' + n : n.toString();
    }

    // Order by ascending property value
    valueAscOrder = (a: KeyValue<number, string>, b: KeyValue<number, string>): number => {
        return a.value.localeCompare(b.value);
    };

    // Order by descending property key
    keyDescOrder = (a: KeyValue<number, string>, b: KeyValue<number, string>): number => {
        return a.key > b.key ? -1 : (b.key > a.key ? 1 : 0);
    };
}
