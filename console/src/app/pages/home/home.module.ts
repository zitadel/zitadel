import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { AvatarModule } from 'src/app/modules/avatar/avatar.module';
import { CardModule } from 'src/app/modules/card/card.module';
import { OnboardingModule } from 'src/app/modules/onboarding/onboarding.module';
import { SharedModule } from 'src/app/modules/shared/shared.module';

import { HomeRoutingModule } from './home-routing.module';
import { HomeComponent } from './home.component';

@NgModule({
    declarations: [HomeComponent],
    imports: [
        CommonModule,
        MatIconModule,
        HasRoleModule,
        HomeRoutingModule,
        MatButtonModule,
        TranslateModule,
        AvatarModule,
        MatTooltipModule,
        SharedModule,
        MatProgressSpinnerModule,
        CardModule,
        OnboardingModule,
    ],
})
export class HomeModule { }
