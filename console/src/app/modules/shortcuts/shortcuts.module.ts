import { DragDropModule } from '@angular/cdk/drag-drop';
import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';

import { ShortcutsComponent } from './shortcuts.component';

@NgModule({
  declarations: [ShortcutsComponent],
  imports: [CommonModule, DragDropModule],
  exports: [ShortcutsComponent],
})
export class ShortcutsModule {}
