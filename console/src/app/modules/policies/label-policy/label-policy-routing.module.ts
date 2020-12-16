import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { LabelPolicyComponent } from './label-policy.component';

const routes: Routes = [
    {
        path: '',
        component: LabelPolicyComponent,
        data: {
            animation: 'DetailPage',
        },
    },
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule],
})
export class LabelPolicyRoutingModule { }
