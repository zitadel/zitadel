import { Component, OnDestroy } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Subscription } from 'rxjs';
import { GetLabelPolicyResponse, UpdateLabelPolicyRequest } from 'src/app/proto/generated/zitadel/admin_pb';
import { LabelPolicy } from 'src/app/proto/generated/zitadel/policy_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ToastService } from 'src/app/services/toast.service';
import { CnslLinks } from '../../links/links.component';
import { IAM_COMPLEXITY_LINK, IAM_LABEL_LINK, IAM_LOGIN_POLICY_LINK, IAM_POLICY_LINK } from '../../policy-grid/policy-links';

import { CnslLinks } from '../../links/links.component';
import { IAM_COMPLEXITY_LINK, IAM_LOGIN_POLICY_LINK, IAM_POLICY_LINK } from '../../policy-grid/policy-links';
import { PolicyComponentServiceType } from '../policy-component-types.enum';


@Component({
    selector: 'app-label-policy',
    templateUrl: './label-policy.component.html',
    styleUrls: ['./label-policy.component.scss'],
})
export class LabelPolicyComponent implements OnDestroy {
    public labelData!: LabelPolicy.AsObject;

    private sub: Subscription = new Subscription();

    public PolicyComponentServiceType: any = PolicyComponentServiceType;
    public nextLinks: CnslLinks[] = [
        IAM_COMPLEXITY_LINK,
        IAM_POLICY_LINK,
        IAM_LOGIN_POLICY_LINK,
    ];
    constructor(
        private route: ActivatedRoute,
        private toast: ToastService,
        private adminService: AdminService,
    ) {
        this.route.params.subscribe(() => {
            this.getData().then(data => {
                if (data?.policy) {
                    this.labelData = data.policy;
                }
            });
        });
    }

    public ngOnDestroy(): void {
        this.sub.unsubscribe();
    }

    private async getData(): Promise<GetLabelPolicyResponse.AsObject> {
        return this.adminService.getLabelPolicy();
    }

    public savePolicy(): void {
        const req = new UpdateLabelPolicyRequest();
        req.setPrimaryColor(this.labelData.primaryColor);
        req.setSecondaryColor(this.labelData.secondaryColor);
        this.adminService.updateLabelPolicy(req).then(() => {
            this.toast.showInfo('POLICY.TOAST.SET', true);
        }).catch(error => {
            this.toast.showError(error);
        });
    }
}
