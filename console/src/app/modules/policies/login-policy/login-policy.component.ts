import { Component, Injector, OnDestroy, Type } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';
import { Subscription } from 'rxjs';
import { switchMap } from 'rxjs/operators';
import { DefaultLoginPolicy, DefaultLoginPolicyView } from 'src/app/proto/generated/admin_pb';
import { LoginPolicy, LoginPolicyView } from 'src/app/proto/generated/management_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { StorageService } from 'src/app/services/storage.service';
import { ToastService } from 'src/app/services/toast.service';

import { PolicyComponentAction } from '../policy-component-action.enum';
import { PolicyComponentServiceType } from '../policy-component-types.enum';

@Component({
    selector: 'app-login-policy',
    templateUrl: './login-policy.component.html',
    styleUrls: ['./login-policy.component.scss'],
})
export class LoginPolicyComponent implements OnDestroy {
    public title: string = '';
    public desc: string = '';

    componentAction: PolicyComponentAction = PolicyComponentAction.CREATE;

    public PolicyComponentAction: any = PolicyComponentAction;

    public loginData!: LoginPolicy.AsObject | DefaultLoginPolicy.AsObject;

    private sub: Subscription = new Subscription();
    private service!: ManagementService | AdminService;
    private serviceType: PolicyComponentServiceType = PolicyComponentServiceType.MGMT;
    constructor(
        private route: ActivatedRoute,
        private router: Router,
        private toast: ToastService,
        private sessionStorage: StorageService,
        private injector: Injector,
    ) {
        this.sub = this.route.data.pipe(switchMap(data => {
            this.componentAction = data.action;
            console.log(data.serviceType);
            this.serviceType = data.serviceType;
            switch (this.serviceType) {
                case PolicyComponentServiceType.MGMT:
                    console.log('mgmt');
                    this.service = this.injector.get(ManagementService as Type<ManagementService>);
                    break;
                case PolicyComponentServiceType.ADMIN:
                    console.log('admin');
                    this.service = this.injector.get(AdminService as Type<AdminService>);
                    break;
            }

            return this.route.params;
        })).subscribe(() => {
            this.title = 'ORG.POLICY.LOGIN_POLICY.TITLECREATE';
            this.desc = 'ORG.POLICY.LOGIN_POLICY.DESCRIPTIONCREATE';

            if (this.componentAction === PolicyComponentAction.MODIFY) {
                this.getData().then(data => {
                    if (data) {
                        this.loginData = data.toObject();
                    }
                });
            }
        });
    }

    public ngOnDestroy(): void {
        this.sub.unsubscribe();
    }

    private async getData():
        Promise<LoginPolicyView | DefaultLoginPolicyView> {
        this.title = 'ORG.POLICY.LOGIN_POLICY.TITLECREATE';
        this.desc = 'ORG.POLICY.LOGIN_POLICY.DESCRIPTIONCREATE';
        switch (this.serviceType) {
            case PolicyComponentServiceType.MGMT:
                return (this.service as ManagementService).GetLoginPolicy();
            case PolicyComponentServiceType.ADMIN:
                return (this.service as AdminService).GetDefaultLoginPolicy();
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
        if (this.componentAction === PolicyComponentAction.CREATE) {
            // const orgId = this.sessionStorage.getItem('organization');
            // if (orgId) {
            //     this.service.CreateOrgIamPolicy(
            //         orgId,
            //         this.iamData.description,
            //         this.iamData.userLoginMustBeDomain,
            //     ).then(() => {
            //         this.router.navigate(['org']);
            //     }).catch(error => {
            //         this.toast.showError(error);
            //     });
            // }
        } else if (this.componentAction === PolicyComponentAction.MODIFY) {
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
    }

    public deletePolicy(): Promise<Empty> {
        switch (this.serviceType) {
            case PolicyComponentServiceType.MGMT:
                return (this.service as ManagementService).RemoveLoginPolicy();
            case PolicyComponentServiceType.ADMIN:
                return (this.service as AdminService).GetDefaultLoginPolicy();
        }
    }
}
