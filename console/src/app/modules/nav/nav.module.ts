import { OverlayModule } from '@angular/cdk/overlay';
import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTooltipModule } from '@angular/material/tooltip';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { HasFeatureModule } from 'src/app/directives/has-feature/has-feature.module';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { HasFeaturePipeModule } from 'src/app/pipes/has-feature-pipe/has-feature-pipe.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';

import { InfoOverlayModule } from '../info-overlay/info-overlay.module';
import { SharedModule } from '../shared/shared.module';
import { NavComponent } from './nav.component';

@NgModule({
  declarations: [NavComponent],
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    InfoOverlayModule,
    TranslateModule,
    MatIconModule,
    RouterModule,
    MatTooltipModule,
    HasRolePipeModule,
    HasFeaturePipeModule,
    HasFeatureModule,
    HasRoleModule,
    MatMenuModule,
    MatButtonModule,
    MatProgressSpinnerModule,
    SharedModule,
    OverlayModule,
  ],
  exports: [NavComponent],
})
export class NavModule {}
