import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { DetailLayoutModule } from 'src/app/modules/detail-layout/detail-layout.module';
import { InputModule } from 'src/app/modules/input/input.module';

import { InfoSectionModule } from '../info-section/info-section.module';
import { FeaturesRoutingModule } from './features-routing.module';
import { FeaturesComponent } from './features.component';

@NgModule({
    declarations: [
        FeaturesComponent
    ],
    imports: [
        FeaturesRoutingModule,
        CommonModule,
        FormsModule,
        InputModule,
        MatButtonModule,
        MatSlideToggleModule,
        MatIconModule,
        HasRoleModule,
        MatTooltipModule,
        InfoSectionModule,
        TranslateModule,
        DetailLayoutModule,
    ],
    exports: [
        FeaturesComponent,
    ]
})
export class FeaturesModule { }
