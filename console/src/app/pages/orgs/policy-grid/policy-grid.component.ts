import { Component } from '@angular/core';
import { PolicyComponentType } from 'src/app/modules/policies/policy-component-types.enum';
import { LoginPolicy, OrgIamPolicy, PasswordComplexityPolicy, PolicyState } from 'src/app/proto/generated/management_pb';
import { ManagementService } from 'src/app/services/mgmt.service';

@Component({
    selector: 'app-policy-grid',
    templateUrl: './policy-grid.component.html',
    styleUrls: ['./policy-grid.component.scss'],
})
export class PolicyGridComponent {
    public complexityPolicy!: PasswordComplexityPolicy.AsObject;
    public iamPolicy!: OrgIamPolicy.AsObject;
    public loginPolicy!: LoginPolicy.AsObject;

    public PolicyState: any = PolicyState;
    public PolicyComponentType: any = PolicyComponentType;

    constructor(
        private mgmtService: ManagementService,
    ) {
        this.getData();
    }

    private getData(): void {
        this.mgmtService.GetPasswordComplexityPolicy().then(data => this.complexityPolicy = data.toObject());
        this.mgmtService.GetMyOrgIamPolicy().then(data => this.iamPolicy = data.toObject());
        this.mgmtService.GetLoginPolicy().then(data => {
            this.loginPolicy = data.toObject();
        });
    }
}
