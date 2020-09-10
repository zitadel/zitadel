import { Component, Injector, OnDestroy, Type } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { ActivatedRoute, Router } from '@angular/router';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';
import { Subscription } from 'rxjs';
import { switchMap } from 'rxjs/operators';
import {
    DefaultLoginPolicy,
    DefaultLoginPolicyView,
    IdpProviderView as AdminIdpProviderView,
    IdpView as AdminIdpView,
} from 'src/app/proto/generated/admin_pb';
import {
    IdpProviderType,
    IdpProviderView as MgmtIdpProviderView,
    IdpView as MgmtIdpView,
    LoginPolicy,
    LoginPolicyView,
} from 'src/app/proto/generated/management_pb';
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
    public loginData!: LoginPolicy.AsObject | DefaultLoginPolicy.AsObject;

    private sub: Subscription = new Subscription();
    private service!: ManagementService | AdminService;
    PolicyComponentServiceType: any = PolicyComponentServiceType;
    public serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;
    public idps: MgmtIdpProviderView.AsObject[] | AdminIdpProviderView.AsObject[] = [];
    constructor(
        private route: ActivatedRoute,
        private router: Router,
        private toast: ToastService,
        private dialog: MatDialog,
        private injector: Injector,
    ) {
        this.sub = this.route.data.pipe(switchMap(data => {
            console.log(data.serviceType);
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
            this.getData().then(data => {
                if (data) {
                    this.loginData = data.toObject();
                }
            });
            this.getIdps().then(idps => {
                console.log(idps);
                this.idps = idps;
            });
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
                return (this.service as ManagementService).UpdateLoginPolicy(mgmtreq);
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
            switch (this.serviceType) {
                case PolicyComponentServiceType.MGMT:
                    this.router.navigate(['org']);
                    break;
                case PolicyComponentServiceType.ADMIN:
                    this.router.navigate(['iam']);
                    break;
            }
        }).catch(error => {
            this.toast.showError(error);
        });
    }

    public deletePolicy(): Promise<Empty> {
        switch (this.serviceType) {
            case PolicyComponentServiceType.MGMT:
                return (this.service as ManagementService).RemoveLoginPolicy();
            case PolicyComponentServiceType.ADMIN:
                return (this.service as AdminService).GetDefaultLoginPolicy();
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
                this.addIdp(resp.idp, resp.type);
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

    public removeIdp(idp: AdminIdpView.AsObject | MgmtIdpView.AsObject): void {
        switch (this.serviceType) {
            case PolicyComponentServiceType.MGMT:
                (this.service as ManagementService).RemoveIdpProviderFromLoginPolicy(idp.id);
                break;
            case PolicyComponentServiceType.ADMIN:
                (this.service as AdminService).RemoveIdpProviderFromDefaultLoginPolicy(idp.id);
                break;
        }
    }
}
