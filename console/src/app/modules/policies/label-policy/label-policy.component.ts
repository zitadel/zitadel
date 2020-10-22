import { Component, OnDestroy } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Subscription } from 'rxjs';
import { DefaultLabelPolicyUpdate, DefaultLabelPolicyView } from 'src/app/proto/generated/admin_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ToastService } from 'src/app/services/toast.service';

import { PolicyComponentServiceType } from '../policy-component-types.enum';


@Component({
    selector: 'app-label-policy',
    templateUrl: './label-policy.component.html',
    styleUrls: ['./label-policy.component.scss'],
})
export class LabelPolicyComponent implements OnDestroy {
    public labelData!: DefaultLabelPolicyView.AsObject;

    private sub: Subscription = new Subscription();

    public PolicyComponentServiceType: any = PolicyComponentServiceType;
    constructor(
        private route: ActivatedRoute,
        private toast: ToastService,
        private adminService: AdminService,
    ) {
        this.route.params.subscribe(() => {
            this.getData().then(data => {
                if (data) {
                    this.labelData = data.toObject();
                }
            });
        });
    }

    public ngOnDestroy(): void {
        this.sub.unsubscribe();
    }

    private async getData(): Promise<DefaultLabelPolicyView> {
        return this.adminService.GetDefaultLabelPolicy();
    }

    public savePolicy(): void {
        const req = new DefaultLabelPolicyUpdate();
        req.setPrimaryColor(this.labelData.primaryColor);
        req.setSecondaryColor(this.labelData.secondaryColor);
        this.adminService.UpdateDefaultLabelPolicy(req).then(() => {
            this.toast.showInfo('POLICY.TOAST.SET', true);
        }).catch(error => {
            this.toast.showError(error);
        });
    }
}
