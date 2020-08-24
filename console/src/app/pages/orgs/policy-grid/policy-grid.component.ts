import { Component } from '@angular/core';
import {
    OrgIamPolicy,
    PasswordAgePolicy,
    PasswordComplexityPolicy,
    PasswordLockoutPolicy,
    PolicyState,
} from 'src/app/proto/generated/management_pb';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ManagementService } from 'src/app/services/mgmt.service';

import { PolicyComponentType } from '../password-policy/password-policy.component';

@Component({
    selector: 'app-policy-grid',
    templateUrl: './policy-grid.component.html',
    styleUrls: ['./policy-grid.component.scss'],
})
export class PolicyGridComponent {
    public lockoutPolicy!: PasswordLockoutPolicy.AsObject;
    public agePolicy!: PasswordAgePolicy.AsObject;
    public complexityPolicy!: PasswordComplexityPolicy.AsObject;
    public iamPolicy!: OrgIamPolicy.AsObject;

    public PolicyState: any = PolicyState;
    public PolicyComponentType: any = PolicyComponentType;

    constructor(
        private orgService: ManagementService,
        public authUserService: GrpcAuthService,
    ) {
        this.getData();
    }

    private getData(): void {
        this.orgService.GetPasswordComplexityPolicy().then(data => this.complexityPolicy = data.toObject());
        this.orgService.GetMyOrgIamPolicy().then(data => this.iamPolicy = data.toObject());
    }
}
