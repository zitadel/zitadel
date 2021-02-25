import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { BehaviorSubject, from, Observable, of } from 'rxjs';
import { catchError, finalize, map } from 'rxjs/operators';
import { Application, OIDCApplicationType, OIDCResponseType } from 'src/app/proto/generated/management_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { NATIVE_TYPE, USER_AGENT_TYPE, WEB_TYPE } from '../../../apps/authtypes';

@Component({
    selector: 'app-application-grid',
    templateUrl: './application-grid.component.html',
    styleUrls: ['./application-grid.component.scss'],
})
export class ApplicationGridComponent implements OnInit {
    @Input() public projectId: string = '';
    @Input() public disabled: boolean = false;
    @Output() public changeView: EventEmitter<void> = new EventEmitter();
    public appsSubject: BehaviorSubject<Application.AsObject[]> = new BehaviorSubject<Application.AsObject[]>([]);
    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(true);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();
    public OIDCApplicationType: any = OIDCApplicationType;

    public NATIVE_TYPE: any = NATIVE_TYPE;
    public WEB_TYPE: any = WEB_TYPE;
    public USER_AGENT_TYPE: any = USER_AGENT_TYPE;

    constructor(private mgmtService: ManagementService) { }

    public ngOnInit(): void {
        this.loadApps();
    }

    public loadApps(): void {
        from(this.mgmtService.SearchApplications(this.projectId, 100, 0)).pipe(
            map(resp => {
                return resp.toObject().resultList;
            }),
            catchError(() => of([])),
            finalize(() => this.loadingSubject.next(false)),
        ).subscribe((apps) => {
            this.appsSubject.next(apps as Application.AsObject[]);
            console.log(apps);
        });
    }

    public closeView(): void {
        this.changeView.emit();
    }
}
