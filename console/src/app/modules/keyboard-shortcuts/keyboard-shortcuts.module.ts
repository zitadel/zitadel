import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';

import { KeyboardShortcutsComponent } from './keyboard-shortcuts.component';

@NgModule({
  declarations: [KeyboardShortcutsComponent],
  imports: [CommonModule, FormsModule, TranslateModule, HasRoleModule, MatButtonModule, MatIconModule],
})
export class KeyboardShortcutsModule {}
