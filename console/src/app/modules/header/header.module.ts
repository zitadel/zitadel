import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatRippleModule } from '@angular/material/core';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatToolbarModule } from '@angular/material/toolbar';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { OutsideClickModule } from 'src/app/directives/outside-click/outside-click.module';
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
    MatButtonModule,
    HasRoleModule,
    MatProgressSpinnerModule,
    TranslateModule,
    OrgContextModule,
    OutsideClickModule,
    AvatarModule,
    AccountsCardModule,
    HasRolePipeModule,
  ],
  exports: [HeaderComponent],
})
export class HeaderModule {}
