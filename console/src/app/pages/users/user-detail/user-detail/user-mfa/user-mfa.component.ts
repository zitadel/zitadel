import { Component, Input, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { MatSort } from '@angular/material/sort';
import { MatTable, MatTableDataSource } from '@angular/material/table';
import { BehaviorSubject, Observable } from 'rxjs';
import { MFAState, MfaType, UserMultiFactor, UserView } from 'src/app/proto/generated/management_pb';
import { ManagementService } from 'src/app/services/mgmt.service';


export interface MFAItem {
    name: string;
    verified: boolean;
}

@Component({
    selector: 'app-user-mfa',
    templateUrl: './user-mfa.component.html',
    styleUrls: ['./user-mfa.component.scss'],
})
export class UserMfaComponent implements OnInit, OnDestroy {
    public displayedColumns: string[] = ['type', 'state'];
    @Input() private user!: UserView.AsObject;
    public mfaSubject: BehaviorSubject<UserMultiFactor.AsObject[]> = new BehaviorSubject<UserMultiFactor.AsObject[]>([]);
    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();

    @ViewChild(MatTable) public table!: MatTable<UserMultiFactor.AsObject>;
    @ViewChild(MatSort) public sort!: MatSort;
    public dataSource!: MatTableDataSource<UserMultiFactor.AsObject>;

    public MfaType: any = MfaType;
    public MFAState: any = MFAState;

    public error: string = '';
    constructor(private mgmtUserService: ManagementService) { }

    public ngOnInit(): void {
        this.getOTP();
    }

    public ngOnDestroy(): void {
        this.mfaSubject.complete();
        this.loadingSubject.complete();
    }

    public getOTP(): void {
        this.mgmtUserService.getUserMfas(this.user.id).then(mfas => {
            this.dataSource = new MatTableDataSource(mfas.toObject().mfasList);
            this.dataSource.sort = this.sort;
        }).catch(error => {
            this.error = error.message;
        });
    }
}
