import { OverlayModule } from '@angular/cdk/overlay';
import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatRippleModule } from '@angular/material/core';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyProgressSpinnerModule as MatProgressSpinnerModule } from '@angular/material/legacy-progress-spinner';
import { MatToolbarModule } from '@angular/material/toolbar';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';

import { AccountsCardModule } from '../accounts-card/accounts-card.module';
import { ActionKeysModule } from '../action-keys/action-keys.module';
import { AvatarModule } from '../avatar/avatar.module';
import { OrgContextModule } from '../org-context/org-context.module';
import { HeaderComponent } from './header.component';

@NgModule({
  declarations: [HeaderComponent],
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    RouterModule,
    MatToolbarModule,
    ActionKeysModule,
    MatRippleModule,
    MatIconModule,
    OverlayModule,
    MatButtonModule,
    HasRoleModule,
    MatProgressSpinnerModule,
    TranslateModule,
    OrgContextModule,
    AvatarModule,
    AccountsCardModule,
    HasRolePipeModule,
  ],
  exports: [HeaderComponent],
})
export class HeaderModule {}
