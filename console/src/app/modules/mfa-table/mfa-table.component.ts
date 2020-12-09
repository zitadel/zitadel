import { SelectionModel } from '@angular/cdk/collections';
import { Component, EventEmitter, Input, OnInit, Output, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatPaginator, PageEvent } from '@angular/material/paginator';
import { MatTableDataSource } from '@angular/material/table';
import { RouterLink } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { MultiFactor as AdminMultiFactor, MultiFactorType as AdminMultiFactorType } from 'src/app/proto/generated/admin_pb';
import { MultiFactor as MgmtMultiFactor, MultiFactorType as MgmtMultiFactorType } from 'src/app/proto/generated/management_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { PolicyComponentServiceType } from '../policies/policy-component-types.enum';
import { WarnDialogComponent } from '../warn-dialog/warn-dialog.component';
import { DialogAddTypeComponent } from './dialog-add-type/dialog-add-type.component';

@Component({
    selector: 'app-mfa-table',
    templateUrl: './mfa-table.component.html',
    styleUrls: ['./mfa-table.component.scss'],
})
export class MfaTableComponent implements OnInit {
    @Input() public serviceType!: PolicyComponentServiceType;
    @Input() service!: AdminService | ManagementService;
    @Input() disabled: boolean = false;
    @ViewChild(MatPaginator) public paginator!: MatPaginator;
    public mfas: Array<AdminMultiFactorType | MgmtMultiFactorType> = [];

    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();

    public PolicyComponentServiceType: any = PolicyComponentServiceType;
    @Input() public displayedColumns: string[] = ['type'];

    constructor(public translate: TranslateService, private toast: ToastService, private dialog: MatDialog) { }

    ngOnInit(): void {
        this.getData();
        if (this.serviceType === PolicyComponentServiceType.MGMT) {
            this.displayedColumns = ['type'];
        } else if (this.serviceType === PolicyComponentServiceType.ADMIN) {
            this.displayedColumns = ['type'];
        }

        if (!this.disabled) {
            this.displayedColumns.push('actions');
        }
    }

    public removeMfa(type: MgmtMultiFactorType | AdminMultiFactorType): void {
        const dialogRef = this.dialog.open(WarnDialogComponent, {
            data: {
                confirmKey: 'ACTIONS.DELETE',
                cancelKey: 'ACTIONS.CANCEL',
                titleKey: 'MFA.DELETE.TITLE',
                descriptionKey: 'MFA.DELETE.DESCRIPTION',
            },
            width: '400px',
        });

        dialogRef.afterClosed().subscribe(resp => {
            if (resp) {
                if (this.serviceType === PolicyComponentServiceType.MGMT) {
                    const req = new MgmtMultiFactor();
                    req.setMultiFactor(type);
                    (this.service as ManagementService).RemoveMultiFactorFromLoginPolicy(req).then(() => {
                        this.toast.showInfo('MFA.TOAST.DELETED', true);
                        this.refreshPageAfterTimout(2000);
                    });
                } else if (this.serviceType === PolicyComponentServiceType.ADMIN) {
                    const req = new AdminMultiFactor();
                    req.setMultiFactor(type);
                    (this.service as AdminService).RemoveMultiFactorFromDefaultLoginPolicy(req).then(() => {
                        this.toast.showInfo('MFA.TOAST.DELETED', true);
                        this.refreshPageAfterTimout(2000);
                    });
                }
            }
        });
    }

    public addMfa(): void {
        const dialogRef = this.dialog.open(DialogAddTypeComponent, {
            data: {
                title: 'MFA.CREATE.TITLE',
                desc: 'MFA.CREATE.DESCRIPTION',
                types:
                    this.serviceType === PolicyComponentServiceType.MGMT ?
                        [MgmtMultiFactorType.MULTIFACTORTYPE_U2F_WITH_PIN] :
                        this.serviceType === PolicyComponentServiceType.ADMIN ?
                            [AdminMultiFactorType.MULTIFACTORTYPE_U2F_WITH_PIN] :
                            [],
            },
            width: '400px',
        });

        dialogRef.afterClosed().subscribe((mfaType: AdminMultiFactorType | MgmtMultiFactorType) => {
            if (mfaType) {
                if (this.serviceType === PolicyComponentServiceType.MGMT) {
                    const req = new MgmtMultiFactor();
                    req.setMultiFactor(mfaType);
                    (this.service as ManagementService).AddMultiFactorToLoginPolicy(req).then(() => {
                        this.refreshPageAfterTimout(2000);
                    }).catch(error => {
                        this.toast.showError(error);
                    });
                } else if (this.serviceType === PolicyComponentServiceType.ADMIN) {
                    const req = new AdminMultiFactor();
                    req.setMultiFactor(mfaType);
                    (this.service as AdminService).addMultiFactorToDefaultLoginPolicy(req).then(() => {
                        this.refreshPageAfterTimout(2000);
                    }).catch(error => {
                        this.toast.showError(error);
                    });
                }
            }
        });
    }

    private async getData(): Promise<void> {
        this.loadingSubject.next(true);

        if (this.serviceType === PolicyComponentServiceType.MGMT) {
            (this.service as ManagementService).GetLoginPolicyMultiFactors().then(resp => {
                this.mfas = resp.toObject().multiFactorsList;
                this.loadingSubject.next(false);
            }).catch(error => {
                this.toast.showError(error);
                this.loadingSubject.next(false);
            });
        } else if (this.serviceType === PolicyComponentServiceType.ADMIN) {
            (this.service as AdminService).getDefaultLoginPolicyMultiFactors().then(resp => {
                this.mfas = resp.toObject().multiFactorsList;
                this.loadingSubject.next(false);
            }).catch(error => {
                this.toast.showError(error);
                this.loadingSubject.next(false);
            });
        }
    }

    public refreshPageAfterTimout(to: number): void {
        setTimeout(() => {
            this.getData();
        }, to);
    }
}
