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

export enum PolicyComponentType {
    LOCKOUT = 'lockout',
    AGE = 'age',
    COMPLEXITY = 'complexity',
    IAM_POLICY = 'iam_policy',
}

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
        private mgmtService: ManagementService,
        public authUserService: GrpcAuthService,
    ) {
        this.getData();
    }

    private getData(): void {
        this.mgmtService.GetPasswordComplexityPolicy().then(data => this.complexityPolicy = data.toObject());
        this.mgmtService.GetMyOrgIamPolicy().then(data => this.iamPolicy = data.toObject());
    }
}
