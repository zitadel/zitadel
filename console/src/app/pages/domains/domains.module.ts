import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyProgressSpinnerModule as MatProgressSpinnerModule } from '@angular/material/legacy-progress-spinner';
import { MatLegacyTooltipModule as MatTooltipModule } from '@angular/material/legacy-tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { CopyToClipboardModule } from 'src/app/directives/copy-to-clipboard/copy-to-clipboard.module';
import { ActionKeysModule } from 'src/app/modules/action-keys/action-keys.module';
import { CardModule } from 'src/app/modules/card/card.module';
import { InfoSectionModule } from 'src/app/modules/info-section/info-section.module';
import { InputModule } from 'src/app/modules/input/input.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';

import { AddDomainDialogModule } from './add-domain-dialog/add-domain-dialog.module';
import { DomainVerificationComponent } from './domain-verification/domain-verification.component';
import { DomainsRoutingModule } from './domains-routing.module';
import { DomainsComponent } from './domains.component';

@NgModule({
  declarations: [DomainsComponent, DomainVerificationComponent],
  imports: [
    DomainsRoutingModule,
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
    TranslateModule,
    InfoSectionModule,
    MatProgressSpinnerModule,
  ],
})
export default class DomainsModule {}
