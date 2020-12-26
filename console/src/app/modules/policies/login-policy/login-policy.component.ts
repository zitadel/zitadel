import { Component, Injector, OnDestroy, Type } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { ActivatedRoute } from '@angular/router';
import { Subscription } from 'rxjs';
import { switchMap } from 'rxjs/operators';
import {
    DefaultLoginPolicy,
    DefaultLoginPolicyView,
    IdpProviderView as AdminIdpProviderView,
    IdpView as AdminIdpView,
} from 'src/app/proto/generated/zitadel/admin_pb';
import {
    IdpProviderType,
    IdpProviderView as MgmtIdpProviderView,
    IdpView as MgmtIdpView,
    LoginPolicy,
    LoginPolicyView,
} from 'src/app/proto/generated/zitadel/management_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { PolicyComponentServiceType } from '../policy-component-types.enum';
import { AddIdpDialogComponent } from './add-idp-dialog/add-idp-dialog.component';

@Component({
    selector: 'app-login-policy',
    templateUrl: './login-policy.component.html',
    styleUrls: ['./login-policy.component.scss'],
})
export class LoginPolicyComponent implements OnDestroy {
    public loginData!: LoginPolicyView.AsObject | DefaultLoginPolicyView.AsObject;

    private sub: Subscription = new Subscription();
    public service!: ManagementService | AdminService;
    public PolicyComponentServiceType: any = PolicyComponentServiceType;
    public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;
    public idps: MgmtIdpProviderView.AsObject[] | AdminIdpProviderView.AsObject[] = [];

    public loading: boolean = false;
    public disabled: boolean = true;
    constructor(
        private route: ActivatedRoute,
        private toast: ToastService,
        private dialog: MatDialog,
        private injector: Injector,
    ) {
        this.sub = this.route.data.pipe(switchMap(data => {
            this.serviceType = data.serviceType;
            switch (this.serviceType) {
                case PolicyComponentServiceType.MGMT:
                    this.service = this.injector.get(ManagementService as Type<ManagementService>);
                    break;
                case PolicyComponentServiceType.ADMIN:
                    this.service = this.injector.get(AdminService as Type<AdminService>);
                    break;
            }

            return this.route.params;
        })).subscribe(() => {
            this.fetchData();
        });
    }

    private fetchData(): void {
        this.getData().then(data => {
            if (data) {
                this.loginData = data.toObject();
                this.loading = false;
                this.disabled = ((this.loginData as LoginPolicyView.AsObject)?.pb_default) ?? false;
            }
        });
        this.getIdps().then(idps => {
            this.idps = idps;
        });
    }

    public ngOnDestroy(): void {
        this.sub.unsubscribe();
    }

    private async getData():
        Promise<LoginPolicyView | DefaultLoginPolicyView> {
        switch (this.serviceType) {
            case PolicyComponentServiceType.MGMT:
                return (this.service as ManagementService).GetLoginPolicy();
            case PolicyComponentServiceType.ADMIN:
                return (this.service as AdminService).GetDefaultLoginPolicy();
        }
    }

    private async getIdps(): Promise<MgmtIdpProviderView.AsObject[] | AdminIdpProviderView.AsObject[]> {
        switch (this.serviceType) {
            case PolicyComponentServiceType.MGMT:
                return (this.service as ManagementService).GetLoginPolicyIdpProviders()
                    .then((providers) => {
                        return providers.toObject().resultList;
                    });
            case PolicyComponentServiceType.ADMIN:
                return (this.service as AdminService).GetDefaultLoginPolicyIdpProviders()
                    .then((providers) => {
                        return providers.toObject().resultList;
                    });
        }
    }

    private async updateData():
        Promise<LoginPolicy | DefaultLoginPolicy> {
        switch (this.serviceType) {
            case PolicyComponentServiceType.MGMT:
                const mgmtreq = new LoginPolicy();
                mgmtreq.setAllowExternalIdp(this.loginData.allowExternalIdp);
                mgmtreq.setAllowRegister(this.loginData.allowRegister);
                mgmtreq.setAllowUsernamePassword(this.loginData.allowUsernamePassword);
                if ((this.loginData as LoginPolicyView.AsObject).pb_default) {
                    return (this.service as ManagementService).CreateLoginPolicy(mgmtreq);
                } else {
                    return (this.service as ManagementService).UpdateLoginPolicy(mgmtreq);
                }
            case PolicyComponentServiceType.ADMIN:
                const adminreq = new DefaultLoginPolicy();
                adminreq.setAllowExternalIdp(this.loginData.allowExternalIdp);
                adminreq.setAllowRegister(this.loginData.allowRegister);
                adminreq.setAllowUsernamePassword(this.loginData.allowUsernamePassword);
                return (this.service as AdminService).UpdateDefaultLoginPolicy(adminreq);
        }
    }

    public savePolicy(): void {
        this.updateData().then(() => {
            this.toast.showInfo('POLICY.LOGIN_POLICY.SAVED', true);
            this.loading = true;
            setTimeout(() => {
                this.fetchData();
            }, 2000);
        }).catch(error => {
            this.toast.showError(error);
        });
    }

    public removePolicy(): void {
        if (this.serviceType === PolicyComponentServiceType.MGMT) {
            (this.service as ManagementService).RemoveLoginPolicy().then(() => {
                this.toast.showInfo('POLICY.TOAST.RESETSUCCESS', true);
                this.loading = true;
                setTimeout(() => {
                    this.fetchData();
                }, 2000);
            }).catch(error => {
                this.toast.showError(error);
            });
        }
    }

    public openDialog(): void {
        const dialogRef = this.dialog.open(AddIdpDialogComponent, {
            data: {
                serviceType: this.serviceType,
            },
            width: '400px',
        });

        dialogRef.afterClosed().subscribe(resp => {
            if (resp && resp.idp && resp.type) {
                this.addIdp(resp.idp, resp.type).then(() => {
                    this.loading = true;
                    setTimeout(() => {
                        this.fetchData();
                    }, 2000);
                }).catch(error => {
                    this.toast.showError(error);
                });
            }
        });
    }

    private addIdp(idp: AdminIdpView.AsObject | MgmtIdpView.AsObject,
        type: IdpProviderType = IdpProviderType.IDPPROVIDERTYPE_SYSTEM): Promise<any> {
        switch (this.serviceType) {
            case PolicyComponentServiceType.MGMT:
                return (this.service as ManagementService).addIdpProviderToLoginPolicy(idp.id, type);
            case PolicyComponentServiceType.ADMIN:
                return (this.service as AdminService).AddIdpProviderToDefaultLoginPolicy(idp.id);
        }
    }

    public removeIdp(idp: AdminIdpProviderView.AsObject | MgmtIdpProviderView.AsObject): void {
        switch (this.serviceType) {
            case PolicyComponentServiceType.MGMT:
                (this.service as ManagementService).RemoveIdpProviderFromLoginPolicy(idp.idpConfigId).then(() => {
                    const index = this.idps.findIndex(temp => temp === idp);
                    if (index > -1) {
                        this.idps.splice(index, 1);
                    }
                });
                break;
            case PolicyComponentServiceType.ADMIN:
                (this.service as AdminService).RemoveIdpProviderFromDefaultLoginPolicy(idp.idpConfigId).then(() => {
                    const index = this.idps.findIndex(temp => temp === idp);
                    if (index > -1) {
                        this.idps.splice(index, 1);
                    }
                });
                break;
        }
    }

    public get isDefault(): boolean {
        if (this.loginData && this.serviceType === PolicyComponentServiceType.MGMT) {
            return (this.loginData as LoginPolicyView.AsObject).pb_default;
        } else {
            return false;
        }
    }
}
