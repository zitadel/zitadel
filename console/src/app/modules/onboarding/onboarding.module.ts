import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { TranslateModule } from '@ngx-translate/core';

import { OnboardingRoutingModule } from './onboarding-routing.module';
import { OnboardingComponent } from './onboarding.component';

@NgModule({
    declarations: [OnboardingComponent],
    imports: [
        CommonModule,
        TranslateModule,
        OnboardingRoutingModule,
        MatButtonModule,
    ],
})
export class OnboardingModule { }
