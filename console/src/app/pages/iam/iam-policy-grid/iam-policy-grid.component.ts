import { Component } from '@angular/core';
import { PolicyComponentType } from 'src/app/modules/policies/policy-component-types.enum';
import { DefaultLoginPolicy, DefaultPasswordComplexityPolicyView, OrgIamPolicyView } from 'src/app/proto/generated/admin_pb';
import { PolicyState } from 'src/app/proto/generated/management_pb';
import { AdminService } from 'src/app/services/admin.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';

@Component({
    selector: 'app-iam-policy-grid',
    templateUrl: './iam-policy-grid.component.html',
    styleUrls: ['./iam-policy-grid.component.scss'],
})
export class IamPolicyGridComponent {
    public complexityPolicy!: DefaultPasswordComplexityPolicyView.AsObject;
    public loginPolicy!: DefaultLoginPolicy.AsObject;
    public iamPolicy!: OrgIamPolicyView.AsObject;

    public PolicyState: any = PolicyState;
    public PolicyComponentType: any = PolicyComponentType;

    constructor(
        private authService: GrpcAuthService,
        private adminService: AdminService,
    ) {
        this.getData();
    }

    private getData(): void {

        this.authService.isAllowed(['policy.read']).subscribe(allowed => {
            if (allowed) {
                this.adminService.GetDefaultLoginPolicy().then(data => this.loginPolicy = data.toObject());
                this.adminService.GetDefaultPasswordComplexityPolicy().then(data => this.complexityPolicy = data.toObject());
            }
        });

        this.authService.isAllowed(['iam.policy.read']).subscribe(allowed => {
            if (allowed) {
                this.adminService.GetDefaultOrgIamPolicy().then(data => this.iamPolicy = data.toObject());
            }
        });
    }
}
