import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatRippleModule } from '@angular/material/core';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { SharedModule } from 'src/app/modules/shared/shared.module';
import { ShortcutsModule } from 'src/app/modules/shortcuts/shortcuts.module';

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
    MatTooltipModule,
    SharedModule,
    MatProgressSpinnerModule,
    ShortcutsModule,
    MatRippleModule,
  ],
})
export class HomeModule {}
