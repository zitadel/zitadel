import { Component, Input, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatSort } from '@angular/material/sort';
import { MatTable, MatTableDataSource } from '@angular/material/table';
import { BehaviorSubject, Observable } from 'rxjs';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import { MultiFactorType } from 'src/app/proto/generated/zitadel/policy_pb';
import { MultiFactor, MultiFactorState, User } from 'src/app/proto/generated/zitadel/user_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';


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
    public displayedColumns: string[] = ['type', 'attr', 'state', 'actions'];
    @Input() private user!: User.AsObject;
    public mfaSubject: BehaviorSubject<MultiFactor.AsObject[]> = new BehaviorSubject<MultiFactor.AsObject[]>([]);
    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();

    @ViewChild(MatTable) public table!: MatTable<MultiFactor.AsObject>;
    @ViewChild(MatSort) public sort!: MatSort;
    public dataSource!: MatTableDataSource<MultiFactor.AsObject>;

    public MultiFactorType: any = MultiFactorType;
    public MultiFactorState: any = MultiFactorState;

    public error: string = '';
    constructor(private mgmtUserService: ManagementService, private dialog: MatDialog, private toast: ToastService) { }

    public ngOnInit(): void {
        this.getMFAs();
    }

    public ngOnDestroy(): void {
        this.mfaSubject.complete();
        this.loadingSubject.complete();
    }

    public getMFAs(): void {
        this.mgmtUserService.listHumanMultiFactors(this.user.id).then(mfas => {
            this.dataSource = new MatTableDataSource(mfas.resultList);
            this.dataSource.sort = this.sort;
        }).catch(error => {
            this.error = error.message;
        });
    }

    public deleteMFA(type: MultiFactorType, id?: string): void {
        const dialogRef = this.dialog.open(WarnDialogComponent, {
            data: {
                confirmKey: 'ACTIONS.DELETE',
                cancelKey: 'ACTIONS.CANCEL',
                titleKey: 'USER.MFA.DIALOG.MFA_DELETE_TITLE',
                descriptionKey: 'USER.MFA.DIALOG.MFA_DELETE_DESCRIPTION',
            },
            width: '400px',
        });

        dialogRef.afterClosed().subscribe(resp => {
            if (resp) {
                if (type === MultiFactorType.otp) {
                    this.mgmtUserService.removeHumanMultiFactorOTP(this.user.id).then(() => {
                        this.toast.showInfo('USER.TOAST.OTPREMOVED', true);

                        const index = this.dataSource.data.findIndex(mfa => !!mfa.otp);
                        if (index > -1) {
                            this.dataSource.data.splice(index, 1);
                        }
                        this.getMFAs();
                    }).catch(error => {
                        this.toast.showError(error);
                    });
                } else if (type === MultiFactorType.MULTI_FACTOR_TYPE_U2F_WITH_VERIFICATION && id) {
                    this.mgmtUserService.removeHumanMultiFactorU2F(this.user.id, id).then(() => {
                        this.toast.showInfo('USER.TOAST.U2FREMOVED', true);

                        const index = this.dataSource.data.findIndex(mfa => !!mfa.u2f);
                        if (index > -1) {
                            this.dataSource.data.splice(index, 1);
                        }
                        this.getMFAs();
                    }).catch(error => {
                        this.toast.showError(error);
                    });
                }
            }
        });
    }
}
