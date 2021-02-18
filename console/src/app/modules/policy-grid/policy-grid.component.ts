import { Component, Input, OnInit } from '@angular/core';
import { PolicyComponentType } from 'src/app/modules/policies/policy-component-types.enum';
import { PasswordComplexityPolicyView as MgmtPasswordComplexityPolicyView } from 'src/app/proto/generated/management_pb';
import { DefaultPasswordComplexityPolicyView as AdminPasswordComplexityPolicyView } from 'src/app/proto/generated/admin_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { AdminService } from 'src/app/services/admin.service';

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

    public complexityPolicy!: MgmtPasswordComplexityPolicyView.AsObject | AdminPasswordComplexityPolicyView.AsObject | any;
    constructor(private mgmtService: ManagementService, private adminService: AdminService) { }

    public ngOnInit(): void {
        if (this.type == PolicyGridType.ORG) {
            this.mgmtService.GetDefaultPasswordComplexityPolicy().then((policy) => {
                this.complexityPolicy = policy.toObject();
            });
        } else if (this.type == PolicyGridType.IAM) {
            this.adminService.GetDefaultPasswordComplexityPolicy().then((policy) => {
                this.complexityPolicy = policy.toObject();
            });
        }
    }
}
