import { Component, Input, OnInit, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatPaginator } from '@angular/material/paginator';
import { TranslateService } from '@ngx-translate/core';
import { BehaviorSubject, Observable } from 'rxjs';
import {
    AddMultiFactorToLoginPolicyRequest as AdminAddMultiFactorToLoginPolicyRequest,
    AddSecondFactorToLoginPolicyRequest as AdminAddSecondFactorToLoginPolicyRequest,
    RemoveMultiFactorFromLoginPolicyRequest as AdminRemoveMultiFactorFromLoginPolicyRequest,
    RemoveSecondFactorFromLoginPolicyRequest as AdminRemoveSecondFactorFromLoginPolicyRequest,
} from 'src/app/proto/generated/zitadel/admin_pb';
import {
    AddMultiFactorToLoginPolicyRequest as MgmtAddMultiFactorToLoginPolicyRequest,
    AddSecondFactorToLoginPolicyRequest as MgmtAddSecondFactorToLoginPolicyRequest,
    RemoveMultiFactorFromLoginPolicyRequest as MgmtRemoveMultiFactorFromLoginPolicyRequest,
    RemoveSecondFactorFromLoginPolicyRequest as MgmtRemoveSecondFactorFromLoginPolicyRequest,
} from 'src/app/proto/generated/zitadel/management_pb';
import { MultiFactorType, SecondFactorType } from 'src/app/proto/generated/zitadel/policy_pb';
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
    public mfas: Array<MultiFactorType | SecondFactorType> = [];

    private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
    public loading$: Observable<boolean> = this.loadingSubject.asObservable();

    public PolicyComponentServiceType: any = PolicyComponentServiceType;

    constructor(public translate: TranslateService, private toast: ToastService, private dialog: MatDialog) { }

    public ngOnInit(): void {
        this.getData();
    }

    public removeMfa(type: MultiFactorType | SecondFactorType): void {
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
                        const req = new MgmtRemoveMultiFactorFromLoginPolicyRequest();
                        req.setType(type as MultiFactorType);
                        (this.service as ManagementService).removeMultiFactorFromLoginPolicy(req).then(() => {
                            this.toast.showInfo('MFA.TOAST.DELETED', true);
                            this.refreshPageAfterTimout(2000);
                        });
                    } else if (this.componentType === LoginMethodComponentType.SecondFactor) {
                        const req = new MgmtRemoveSecondFactorFromLoginPolicyRequest();
                        req.setType(type as SecondFactorType);
                        (this.service as ManagementService).removeSecondFactorFromLoginPolicy(req).then(() => {
                            this.toast.showInfo('MFA.TOAST.DELETED', true);
                            this.refreshPageAfterTimout(2000);
                        });
                    }
                } else if (this.serviceType === PolicyComponentServiceType.ADMIN) {
                    if (this.componentType === LoginMethodComponentType.MultiFactor) {
                        const req = new AdminRemoveMultiFactorFromLoginPolicyRequest();
                        req.setType(type as MultiFactorType);
                        (this.service as AdminService).removeMultiFactorFromLoginPolicy(req).then(() => {
                            this.toast.showInfo('MFA.TOAST.DELETED', true);
                            this.refreshPageAfterTimout(2000);
                        });
                    } else if (this.componentType === LoginMethodComponentType.SecondFactor) {
                        const req = new AdminRemoveSecondFactorFromLoginPolicyRequest();
                        req.setType(type as SecondFactorType);
                        (this.service as AdminService).removeSecondFactorFromLoginPolicy(req).then(() => {
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
            selection = [MultiFactorType.MULTI_FACTOR_TYPE_U2F_WITH_VERIFICATION];
        } else if (this.componentType === LoginMethodComponentType.SecondFactor) {
            selection = [SecondFactorType.SECOND_FACTOR_TYPE_U2F, SecondFactorType.SECOND_FACTOR_TYPE_OTP];
        }

        this.mfas.forEach(mfa => {
            const index = selection.findIndex(sel => sel === mfa);
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

        dialogRef.afterClosed().subscribe((mfaType: MultiFactorType | SecondFactorType) => {
            if (mfaType) {
                if (this.serviceType === PolicyComponentServiceType.MGMT) {
                    if (this.componentType === LoginMethodComponentType.MultiFactor) {
                        const req = new MgmtAddMultiFactorToLoginPolicyRequest();
                        req.setType(mfaType as MultiFactorType);
                        (this.service as ManagementService).addMultiFactorToLoginPolicy(req).then(() => {
                            this.refreshPageAfterTimout(2000);
                        }).catch(error => {
                            this.toast.showError(error);
                        });
                    } else if (this.componentType === LoginMethodComponentType.SecondFactor) {
                        const req = new MgmtAddSecondFactorToLoginPolicyRequest();
                        req.setType(mfaType as SecondFactorType);
                        (this.service as ManagementService).addSecondFactorToLoginPolicy(req).then(() => {
                            this.refreshPageAfterTimout(2000);
                        }).catch(error => {
                            this.toast.showError(error);
                        });
                    }
                } else if (this.serviceType === PolicyComponentServiceType.ADMIN) {
                    if (this.componentType === LoginMethodComponentType.MultiFactor) {
                        const req = new AdminAddMultiFactorToLoginPolicyRequest();
                        req.setType(mfaType as MultiFactorType);
                        (this.service as AdminService).addMultiFactorToLoginPolicy(req).then(() => {
                            this.refreshPageAfterTimout(2000);
                        }).catch(error => {
                            this.toast.showError(error);
                        });
                    } else if (this.componentType === LoginMethodComponentType.SecondFactor) {
                        const req = new AdminAddSecondFactorToLoginPolicyRequest();
                        req.setType(mfaType as SecondFactorType);
                        (this.service as AdminService).addSecondFactorToLoginPolicy(req).then(() => {
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
                (this.service as ManagementService).listLoginPolicyMultiFactors().then(resp => {
                    this.mfas = resp.resultList;
                    this.loadingSubject.next(false);
                }).catch(error => {
                    this.toast.showError(error);
                    this.loadingSubject.next(false);
                });
            } else if (this.componentType === LoginMethodComponentType.SecondFactor) {
                (this.service as ManagementService).listLoginPolicySecondFactors().then(resp => {
                    this.mfas = resp.resultList;
                    this.loadingSubject.next(false);
                }).catch(error => {
                    this.toast.showError(error);
                    this.loadingSubject.next(false);
                });
            }
        } else if (this.serviceType === PolicyComponentServiceType.ADMIN) {
            if (this.componentType === LoginMethodComponentType.MultiFactor) {
                (this.service as AdminService).listLoginPolicyMultiFactors().then(resp => {
                    this.mfas = resp.resultList;
                    this.loadingSubject.next(false);
                }).catch(error => {
                    this.toast.showError(error);
                    this.loadingSubject.next(false);
                });
            } else if (this.componentType === LoginMethodComponentType.SecondFactor) {
                (this.service as AdminService).listLoginPolicySecondFactors().then(resp => {
                    this.mfas = resp.resultList;
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
