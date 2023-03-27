import { A11yModule } from '@angular/cdk/a11y';
import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyCheckboxModule as MatCheckboxModule } from '@angular/material/legacy-checkbox';
import { MatLegacyChipsModule as MatChipsModule } from '@angular/material/legacy-chips';
import { MatLegacyDialogModule as MatDialogModule } from '@angular/material/legacy-dialog';
import { MatLegacyMenuModule as MatMenuModule } from '@angular/material/legacy-menu';
import { MatLegacyProgressBarModule as MatProgressBarModule } from '@angular/material/legacy-progress-bar';
import { MatLegacyProgressSpinnerModule as MatProgressSpinnerModule } from '@angular/material/legacy-progress-spinner';
import { MatLegacyRadioModule as MatRadioModule } from '@angular/material/legacy-radio';
import { MatLegacySelectModule as MatSelectModule } from '@angular/material/legacy-select';
import { MatLegacySlideToggleModule as MatSlideToggleModule } from '@angular/material/legacy-slide-toggle';
import { MatLegacySliderModule as MatSliderModule } from '@angular/material/legacy-slider';
import { MatLegacyTooltipModule as MatTooltipModule } from '@angular/material/legacy-tooltip';
import { MatStepperModule } from '@angular/material/stepper';
import { CodemirrorModule } from '@ctrl/ngx-codemirror';
import { TranslateModule } from '@ngx-translate/core';
import { CopyToClipboardModule } from 'src/app/directives/copy-to-clipboard/copy-to-clipboard.module';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { AppRadioModule } from 'src/app/modules/app-radio/app-radio.module';
import { CardModule } from 'src/app/modules/card/card.module';
import { ChangesModule } from 'src/app/modules/changes/changes.module';
import { ClientKeysModule } from 'src/app/modules/client-keys/client-keys.module';
import { CreateLayoutModule } from 'src/app/modules/create-layout/create-layout.module';
import { InfoRowModule } from 'src/app/modules/info-row/info-row.module';
import { InfoSectionModule } from 'src/app/modules/info-section/info-section.module';
import { InputModule } from 'src/app/modules/input/input.module';
import { MetaLayoutModule } from 'src/app/modules/meta-layout/meta-layout.module';
import { NameDialogModule } from 'src/app/modules/name-dialog/name-dialog.module';
import { SidenavModule } from 'src/app/modules/sidenav/sidenav.module';
import { TopViewModule } from 'src/app/modules/top-view/top-view.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';
import { OriginPipeModule } from 'src/app/pipes/origin-pipe/origin-pipe.module';
import { RedirectPipeModule } from 'src/app/pipes/redirect-pipe/redirect-pipe.module';

import { AdditionalOriginsComponent } from './additional-origins/additional-origins.component';
import { AppCreateComponent } from './app-create/app-create.component';
import { AppDetailComponent } from './app-detail/app-detail.component';
import { AuthMethodDialogComponent } from './app-detail/auth-method-dialog/auth-method-dialog.component';
import { AppSecretDialogComponent } from './app-secret-dialog/app-secret-dialog.component';
import { AppsRoutingModule } from './apps-routing.module';
import { RedirectUrisComponent } from './redirect-uris/redirect-uris.component';

@NgModule({
  declarations: [
    AppCreateComponent,
    AppDetailComponent,
    AppSecretDialogComponent,
    RedirectUrisComponent,
    AdditionalOriginsComponent,
    AuthMethodDialogComponent,
  ],
  imports: [
    CommonModule,
    A11yModule,
    RedirectPipeModule,
    NameDialogModule,
    AppRadioModule,
    AppsRoutingModule,
    FormsModule,
    InfoRowModule,
    TranslateModule,
    OriginPipeModule,
    ReactiveFormsModule,
    HasRoleModule,
    SidenavModule,
    MatChipsModule,
    CreateLayoutModule,
    ClientKeysModule,
    HasRolePipeModule,
    MatIconModule,
    MatSelectModule,
    MatButtonModule,
    MatProgressSpinnerModule,
    MatProgressBarModule,
    MatDialogModule,
    MatCheckboxModule,
    CardModule,
    TopViewModule,
    MatMenuModule,
    MatTooltipModule,
    TranslateModule,
    MatStepperModule,
    MatRadioModule,
    CopyToClipboardModule,
    MatSlideToggleModule,
    InputModule,
    MetaLayoutModule,
    MatSliderModule,
    CodemirrorModule,
    ChangesModule,
    InfoSectionModule,
  ],
  exports: [TranslateModule],
})
export default class AppsModule {}
