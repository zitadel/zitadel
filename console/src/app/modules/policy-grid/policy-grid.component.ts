import { Component, Input, OnInit } from '@angular/core';
import { PolicyComponentType } from 'src/app/modules/policies/policy-component-types.enum';
import { PasswordComplexityPolicy } from 'src/app/proto/generated/zitadel/policy_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ManagementService } from 'src/app/services/mgmt.service';

export enum PolicyGridType {
    ORG,
    IAM,
}

@Component({
    selector: 'app-policy-grid',
    templateUrl: './policy-grid.component.html',
    styleUrls: ['./policy-grid.component.scss'],
})
export class PolicyGridComponent implements OnInit {
    @Input() public type!: PolicyGridType;
    public PolicyComponentType: any = PolicyComponentType;
    public PolicyGridType: any = PolicyGridType;

    public complexityPolicy!: PasswordComplexityPolicy | any;

    constructor(private mgmtService: ManagementService, private adminService: AdminService) { }

    public ngOnInit(): void {
        if (this.type == PolicyGridType.ORG) {
            this.mgmtService.GetDefaultPasswordComplexityPolicy().then((policy) => {
                this.complexityPolicy = policy;
            });
        } else if (this.type == PolicyGridType.IAM) {
            this.adminService.getPasswordComplexityPolicy().then((resp) => {
                this.complexityPolicy = resp.policy;
            });
        }
    }
}
