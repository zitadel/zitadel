import { OverlayModule } from '@angular/cdk/overlay';
import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatTooltipModule } from '@angular/material/tooltip';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';

import { AccountsCardModule } from '../accounts-card/accounts-card.module';
import { AvatarModule } from '../avatar/avatar.module';
import { SharedModule } from '../shared/shared.module';
import { NavComponent } from './nav.component';

@NgModule({
  declarations: [
    NavComponent,
  ],
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    TranslateModule,
    MatIconModule,
    RouterModule,
    MatTooltipModule,
    AvatarModule,
    AccountsCardModule,
    HasRolePipeModule,
    HasRoleModule,
    MatToolbarModule,
    MatMenuModule,
    MatButtonModule,
    MatProgressSpinnerModule,
    SharedModule,
    OverlayModule,
  ],
  exports: [
    NavComponent,
  ],
})
export class NavModule { }
