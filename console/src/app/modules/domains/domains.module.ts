import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { CopyToClipboardModule } from 'src/app/directives/copy-to-clipboard/copy-to-clipboard.module';
import { ActionKeysModule } from 'src/app/modules/action-keys/action-keys.module';
import { CardModule } from 'src/app/modules/card/card.module';
import { InfoSectionModule } from 'src/app/modules/info-section/info-section.module';
import { InputModule } from 'src/app/modules/input/input.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';

import { AddDomainDialogModule } from './add-domain-dialog/add-domain-dialog.module';
import { DomainVerificationComponent } from './domain-verification/domain-verification.component';
import { DomainsComponent } from './domains.component';
import { MatDialogModule } from '@angular/material/dialog';

@NgModule({
  declarations: [DomainsComponent, DomainVerificationComponent],
  imports: [
    AddDomainDialogModule,
    CommonModule,
    MatIconModule,
    CardModule,
    HasRolePipeModule,
    ActionKeysModule,
    InfoSectionModule,
    MatButtonModule,
    MatTooltipModule,
    CopyToClipboardModule,
    InputModule,
    MatDialogModule,
    TranslateModule,
    InfoSectionModule,
    MatProgressSpinnerModule,
  ],
  exports: [DomainsComponent],
})
export default class DomainsModule {}
