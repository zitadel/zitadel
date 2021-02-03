import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { OnboardingComponent } from './onboarding.component';

const routes: Routes = [
    {
        path: '',
        component: OnboardingComponent,
        data: { animation: 'AddPage' },
    },
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule],
})
export class OnboardingRoutingModule { }
