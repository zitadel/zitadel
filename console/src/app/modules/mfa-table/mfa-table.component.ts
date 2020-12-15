import { Component, Input, OnInit, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatPaginator } from '@angular/material/paginator';
import { TranslateService } from '@ngx-translate/core';
import { BehaviorSubject, Observable } from 'rxjs';
import {
    MultiFactor as AdminMultiFactor,
    MultiFactorType as AdminMultiFactorType,
    SecondFactor as AdminSecondFactor,
    SecondFactorType as AdminSecondFactorType,
} from 'src/app/proto/generated/admin_pb';
import {
    MultiFactor as MgmtMultiFactor,
    MultiFactorType as MgmtMultiFactorType,
    SecondFactor as MgmtSecondFactor,
    SecondFactorType as MgmtSecondFactorType,
} from 'src/app/proto/generated/management_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { PolicyComponentServiceType } from '../policies/policy-component-types.enum';
import { WarnDialogComponent } from '../warn-dialog/warn-dialog.component';
import { DialogAddTypeComponent } from './dialog-add-type/dialog-add-type.component';

export enum LoginMethodComponentType {
    MultiFactor = 1,
    SecondFactor = 2,
}

@Component({
    selector: 'app-mfa-table',
    templateUrl: './mfa-table.component.html',
    styleUrls: ['./mfa-table.component.scss'],
})
export class MfaTableComponent implements OnInit {
    public LoginMethodComponentType: any = LoginMethodComponentType;
    @Input() componentType!: LoginMethodComponentType;
    @Input() public serviceType!: PolicyComponentServiceType;
    @Input() service!: AdminService | ManagementService;
    @Input() disabled: boolean = false;
    @ViewChild(MatPaginator) public paginator!: MatPaginator;
    public mfas: Array<AdminMultiFactorType | MgmtMultiFactorType | MgmtSecondFactorType | AdminSecondFactorType> = [];

    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();

    public PolicyComponentServiceType: any = PolicyComponentServiceType;

    constructor(public translate: TranslateService, private toast: ToastService, private dialog: MatDialog) { }

    public ngOnInit(): void {
        this.getData();
    }

    public removeMfa(type: MgmtMultiFactorType | AdminMultiFactorType | MgmtSecondFactorType | AdminSecondFactorType): void {
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
                    if (this.componentType === LoginMethodComponentType.MultiFactor) {
                        const req = new MgmtMultiFactor();
                        req.setMultiFactor(type as MgmtMultiFactorType);
                        (this.service as ManagementService).RemoveMultiFactorFromLoginPolicy(req).then(() => {
                            this.toast.showInfo('MFA.TOAST.DELETED', true);
                            this.refreshPageAfterTimout(2000);
                        });
                    } else if (this.componentType === LoginMethodComponentType.SecondFactor) {
                        const req = new MgmtSecondFactor();
                        req.setSecondFactor(type as MgmtSecondFactorType);
                        (this.service as ManagementService).RemoveSecondFactorFromLoginPolicy(req).then(() => {
                            this.toast.showInfo('MFA.TOAST.DELETED', true);
                            this.refreshPageAfterTimout(2000);
                        });
                    }
                } else if (this.serviceType === PolicyComponentServiceType.ADMIN) {
                    if (this.componentType === LoginMethodComponentType.MultiFactor) {
                        const req = new AdminMultiFactor();
                        req.setMultiFactor(type as AdminMultiFactorType);
                        (this.service as AdminService).RemoveMultiFactorFromDefaultLoginPolicy(req).then(() => {
                            this.toast.showInfo('MFA.TOAST.DELETED', true);
                            this.refreshPageAfterTimout(2000);
                        });
                    } else if (this.componentType === LoginMethodComponentType.SecondFactor) {
                        const req = new AdminSecondFactor();
                        req.setSecondFactor(type as AdminSecondFactorType);
                        (this.service as AdminService).RemoveSecondFactorFromDefaultLoginPolicy(req).then(() => {
                            this.toast.showInfo('MFA.TOAST.DELETED', true);
                            this.refreshPageAfterTimout(2000);
                        });
                    }
                }
            }
        });
    }

    public addMfa(): void {

        let selection: any[] = [];

        if (this.componentType === LoginMethodComponentType.MultiFactor) {
            selection = this.serviceType === PolicyComponentServiceType.MGMT ?
                [MgmtMultiFactorType.MULTIFACTORTYPE_U2F_WITH_PIN] :
                this.serviceType === PolicyComponentServiceType.ADMIN ?
                    [AdminMultiFactorType.MULTIFACTORTYPE_U2F_WITH_PIN] :
                    [];
        } else if (this.componentType === LoginMethodComponentType.SecondFactor) {
            selection = this.serviceType === PolicyComponentServiceType.MGMT ?
                [MgmtSecondFactorType.SECONDFACTORTYPE_U2F, MgmtSecondFactorType.SECONDFACTORTYPE_OTP] :
                this.serviceType === PolicyComponentServiceType.ADMIN ?
                    [AdminSecondFactorType.SECONDFACTORTYPE_OTP, AdminSecondFactorType.SECONDFACTORTYPE_U2F] :
                    [];
        }

        this.mfas.forEach(mfa => {
            const index = selection.findIndex(sel => sel == mfa);
            if (index > -1) {
                selection.splice(index, 1);
            }
        });

        const dialogRef = this.dialog.open(DialogAddTypeComponent, {
            data: {
                title: 'MFA.CREATE.TITLE',
                desc: 'MFA.CREATE.DESCRIPTION',
                componentType: this.componentType,
                types: selection,
            },
            width: '400px',
        });

        dialogRef.afterClosed().subscribe((mfaType: AdminMultiFactorType | MgmtMultiFactorType |
            AdminSecondFactorType | MgmtSecondFactorType) => {
            if (mfaType) {
                if (this.serviceType === PolicyComponentServiceType.MGMT) {
                    if (this.componentType === LoginMethodComponentType.MultiFactor) {
                        const req = new MgmtMultiFactor();
                        req.setMultiFactor(mfaType as MgmtMultiFactorType);
                        (this.service as ManagementService).AddMultiFactorToLoginPolicy(req).then(() => {
                            this.refreshPageAfterTimout(2000);
                        }).catch(error => {
                            this.toast.showError(error);
                        });
                    } else if (this.componentType === LoginMethodComponentType.SecondFactor) {
                        const req = new MgmtSecondFactor();
                        req.setSecondFactor(mfaType as MgmtSecondFactorType);
                        (this.service as ManagementService).AddSecondFactorToLoginPolicy(req).then(() => {
                            this.refreshPageAfterTimout(2000);
                        }).catch(error => {
                            this.toast.showError(error);
                        });
                    }
                } else if (this.serviceType === PolicyComponentServiceType.ADMIN) {
                    if (this.componentType === LoginMethodComponentType.MultiFactor) {
                        const req = new AdminMultiFactor();
                        req.setMultiFactor(mfaType as AdminMultiFactorType);
                        (this.service as AdminService).addMultiFactorToDefaultLoginPolicy(req).then(() => {
                            this.refreshPageAfterTimout(2000);
                        }).catch(error => {
                            this.toast.showError(error);
                        });
                    } else if (this.componentType === LoginMethodComponentType.SecondFactor) {
                        const req = new AdminSecondFactor();
                        req.setSecondFactor(mfaType as AdminSecondFactorType);
                        (this.service as AdminService).AddSecondFactorToDefaultLoginPolicy(req).then(() => {
                            this.refreshPageAfterTimout(2000);
                        }).catch(error => {
                            this.toast.showError(error);
                        });
                    }
                }
            }
        });
    }

    private async getData(): Promise<void> {
        this.loadingSubject.next(true);

        if (this.serviceType === PolicyComponentServiceType.MGMT) {
            if (this.componentType === LoginMethodComponentType.MultiFactor) {
                (this.service as ManagementService).GetLoginPolicyMultiFactors().then(resp => {
                    this.mfas = resp.toObject().multiFactorsList;
                    this.loadingSubject.next(false);
                }).catch(error => {
                    this.toast.showError(error);
                    this.loadingSubject.next(false);
                });
            } else if (this.componentType === LoginMethodComponentType.SecondFactor) {
                (this.service as ManagementService).GetLoginPolicySecondFactors().then(resp => {
                    this.mfas = resp.toObject().secondFactorsList;
                    this.loadingSubject.next(false);
                }).catch(error => {
                    this.toast.showError(error);
                    this.loadingSubject.next(false);
                });
            }
        } else if (this.serviceType === PolicyComponentServiceType.ADMIN) {
            if (this.componentType === LoginMethodComponentType.MultiFactor) {
                (this.service as AdminService).getDefaultLoginPolicyMultiFactors().then(resp => {
                    this.mfas = resp.toObject().multiFactorsList;
                    this.loadingSubject.next(false);
                }).catch(error => {
                    this.toast.showError(error);
                    this.loadingSubject.next(false);
                });
            } else if (this.componentType === LoginMethodComponentType.SecondFactor) {
                (this.service as AdminService).GetDefaultLoginPolicySecondFactors().then(resp => {
                    this.mfas = resp.toObject().secondFactorsList;
                    this.loadingSubject.next(false);
                }).catch(error => {
                    this.toast.showError(error);
                    this.loadingSubject.next(false);
                });
            }
        }
    }

    public refreshPageAfterTimout(to: number): void {
        setTimeout(() => {
            this.getData();
        }, to);
    }
}
