import { DragDropModule } from '@angular/cdk/drag-drop';
import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { MatTooltipModule } from '@angular/material/tooltip';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';

import { ShortcutsComponent } from './shortcuts.component';

@NgModule({
  declarations: [ShortcutsComponent],
  imports: [
    CommonModule,
    MatButtonModule,
    MatTooltipModule,
    RouterModule,
    DragDropModule,
    HasRoleModule,
    TranslateModule,
    MatIconModule,
  ],
  exports: [ShortcutsComponent],
})
export class ShortcutsModule {}
