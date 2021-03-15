import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { BehaviorSubject, from, Observable, of } from 'rxjs';
import { catchError, finalize, map } from 'rxjs/operators';
import { App, OIDCAppType } from 'src/app/proto/generated/zitadel/app_pb';
import { ManagementService } from 'src/app/services/mgmt.service';

@Component({
    selector: 'app-application-grid',
    templateUrl: './application-grid.component.html',
    styleUrls: ['./application-grid.component.scss'],
})
export class ApplicationGridComponent implements OnInit {
    @Input() public projectId: string = '';
    @Input() public disabled: boolean = false;
    @Output() public changeView: EventEmitter<void> = new EventEmitter();
    public appsSubject: BehaviorSubject<App.AsObject[]> = new BehaviorSubject<App.AsObject[]>([]);
    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(true);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();
    public OIDCAppType: any = OIDCAppType;

    constructor(private mgmtService: ManagementService) { }

    public ngOnInit(): void {
        this.loadApps();
    }

    public loadApps(): void {
        from(this.mgmtService.listApps(this.projectId, 100, 0)).pipe(
            map(resp => {
                console.log(resp.resultList);
                return resp.resultList;
            }),
            catchError(() => of([])),
            finalize(() => this.loadingSubject.next(false)),
        ).subscribe((apps) => {
            console.log(apps);
            this.appsSubject.next(apps as App.AsObject[]);
        });
    }

    public closeView(): void {
        this.changeView.emit();
    }
}
