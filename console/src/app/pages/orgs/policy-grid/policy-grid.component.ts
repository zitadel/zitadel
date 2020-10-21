import { Component } from '@angular/core';
import { PolicyComponentType } from 'src/app/modules/policies/policy-component-types.enum';
import {
    LoginPolicyView,
    OrgIamPolicyView,
    PasswordComplexityPolicyView,
    PolicyState,
} from 'src/app/proto/generated/management_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
    selector: 'app-policy-grid',
    templateUrl: './policy-grid.component.html',
    styleUrls: ['./policy-grid.component.scss'],
})
export class PolicyGridComponent {
    public complexityPolicy!: PasswordComplexityPolicyView.AsObject;
    public iamPolicy!: OrgIamPolicyView.AsObject;
    public loginPolicy!: LoginPolicyView.AsObject;

    public PolicyState: any = PolicyState;
    public PolicyComponentType: any = PolicyComponentType;

    constructor(
        public mgmtService: ManagementService,
        private toast: ToastService,
    ) {
        this.getData();
    }

    private getData(): void {
        this.mgmtService.GetPasswordComplexityPolicy().then(data => this.complexityPolicy = data.toObject()).catch(error => {
            this.toast.showError(error);
        });
        this.mgmtService.GetMyOrgIamPolicy().then(data => this.iamPolicy = data.toObject());
        this.mgmtService.GetLoginPolicy().then(data => {
            this.loginPolicy = data.toObject();
        });
    }
}
