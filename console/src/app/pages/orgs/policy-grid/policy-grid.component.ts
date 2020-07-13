import { Component, OnInit } from '@angular/core';
import {
    OrgIamPolicy,
    PasswordAgePolicy,
    PasswordComplexityPolicy,
    PasswordLockoutPolicy,
    PolicyState,
} from 'src/app/proto/generated/management_pb';
import { AdminService } from 'src/app/services/admin.service';
import { AuthUserService } from 'src/app/services/auth-user.service';
import { OrgService } from 'src/app/services/org.service';
import { ToastService } from 'src/app/services/toast.service';

import { PolicyComponentType } from '../password-policy/password-policy.component';

@Component({
    selector: 'app-policy-grid',
    templateUrl: './policy-grid.component.html',
    styleUrls: ['./policy-grid.component.scss'],
})
export class PolicyGridComponent implements OnInit {
    public lockoutPolicy!: PasswordLockoutPolicy.AsObject;
    public agePolicy!: PasswordAgePolicy.AsObject;
    public complexityPolicy!: PasswordComplexityPolicy.AsObject;
    public iamPolicy!: OrgIamPolicy.AsObject;

    public PolicyState: any = PolicyState;
    public PolicyComponentType: any = PolicyComponentType;

    constructor(
        private orgService: OrgService,
        private adminService: AdminService,
        public authUserService: AuthUserService,
        private toast: ToastService,
    ) {
        this.getData();
    }

    ngOnInit(): void {
    }

    private getData(): void {
        // this.orgService.GetPasswordLockoutPolicy().then(data => this.lockoutPolicy = data.toObject()).catch(error => { });
        // this.orgService.GetPasswordAgePolicy().then(data => this.agePolicy = data.toObject()).catch(error => { });
        this.orgService.GetPasswordComplexityPolicy().then(data => this.complexityPolicy = data.toObject())
            .catch(error => { });
        this.orgService.GetMyOrgIamPolicy().then(data => this.iamPolicy = data.toObject())
            .catch(error => { });
    }
}
