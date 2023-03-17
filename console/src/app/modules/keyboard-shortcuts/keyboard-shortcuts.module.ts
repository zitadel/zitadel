import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';

import { KeyboardShortcutsComponent } from './keyboard-shortcuts.component';

@NgModule({
  declarations: [KeyboardShortcutsComponent],
  imports: [CommonModule, FormsModule, TranslateModule, HasRoleModule, MatButtonModule, MatIconModule],
})
export class KeyboardShortcutsModule {}
