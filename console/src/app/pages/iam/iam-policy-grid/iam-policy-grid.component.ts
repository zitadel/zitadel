import { Component } from '@angular/core';
import { PolicyComponentType } from 'src/app/modules/policies/policy-component-types.enum';
import { DefaultLoginPolicy } from 'src/app/proto/generated/admin_pb';
import { PolicyState } from 'src/app/proto/generated/management_pb';
import { AdminService } from 'src/app/services/admin.service';

@Component({
    selector: 'app-iam-policy-grid',
    templateUrl: './iam-policy-grid.component.html',
    styleUrls: ['./iam-policy-grid.component.scss'],
})
export class IamPolicyGridComponent {
    public loginPolicy!: DefaultLoginPolicy.AsObject;

    public PolicyState: any = PolicyState;
    public PolicyComponentType: any = PolicyComponentType;

    constructor(
        private adminService: AdminService,
    ) {
        this.getData();
    }

    private getData(): void {
        this.adminService.GetDefaultLoginPolicy().then(data => this.loginPolicy = data.toObject());
    }
}
