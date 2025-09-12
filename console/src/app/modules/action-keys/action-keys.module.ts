import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';

import { ActionKeysComponent } from './action-keys.component';

@NgModule({
  declarations: [ActionKeysComponent],
  imports: [CommonModule],
  exports: [ActionKeysComponent],
})
export class ActionKeysModule {}
