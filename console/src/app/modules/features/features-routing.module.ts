import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { FeaturesComponent } from './features.component';

const routes: Routes = [
    {
        path: '',
        component: FeaturesComponent,
        data: {
            animation: 'DetailPage',
        },
    },
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule],
})
export class FeaturesRoutingModule { }
