import { OverlayModule } from '@angular/cdk/overlay';
import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonToggleModule } from '@angular/material/button-toggle';
import { MatExpansionModule } from '@angular/material/expansion';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyCheckboxModule as MatCheckboxModule } from '@angular/material/legacy-checkbox';
import { MatLegacyDialogModule as MatDialogModule } from '@angular/material/legacy-dialog';
import { MatLegacyProgressSpinnerModule as MatProgressSpinnerModule } from '@angular/material/legacy-progress-spinner';
import { MatLegacyTooltipModule as MatTooltipModule } from '@angular/material/legacy-tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { ColorChromeModule } from 'ngx-color/chrome';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';

import { DropzoneModule } from '../../../directives/dropzone/dropzone.module';
import { CardModule } from '../../card/card.module';
import { DetailLayoutModule } from '../../detail-layout/detail-layout.module';
import { InfoSectionModule } from '../../info-section/info-section.module';
import { InputModule } from '../../input/input.module';
import { WarnDialogModule } from '../../warn-dialog/warn-dialog.module';
import { ColorComponent } from './color/color.component';
import { PreviewComponent } from './preview/preview.component';
import { PrivateLabelingPolicyRoutingModule } from './private-labeling-policy-routing.module';
import { PrivateLabelingPolicyComponent } from './private-labeling-policy.component';

@NgModule({
  declarations: [PrivateLabelingPolicyComponent, PreviewComponent, ColorComponent],
  imports: [
    ColorChromeModule,
    PrivateLabelingPolicyRoutingModule,
    CommonModule,
    FormsModule,
    InputModule,
    MatButtonModule,
    MatButtonToggleModule,
    OverlayModule,
    CardModule,
    MatIconModule,
    HasRoleModule,
    MatTooltipModule,
    TranslateModule,
    DetailLayoutModule,
    MatCheckboxModule,
    DropzoneModule,
    MatDialogModule,
    WarnDialogModule,
    HasRolePipeModule,
    MatProgressSpinnerModule,
    MatExpansionModule,
    InfoSectionModule,
  ],
  exports: [PrivateLabelingPolicyComponent],
})
export class PrivateLabelingPolicyModule {}
