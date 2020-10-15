import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { OrgIamPolicyComponent } from './org-iam-policy.component';

const routes: Routes = [
    {
        path: '',
        component: OrgIamPolicyComponent,
        data: {
            animation: 'DetailPage',
        },
    },
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule],
})
export class OrgIamPolicyRoutingModule { }
