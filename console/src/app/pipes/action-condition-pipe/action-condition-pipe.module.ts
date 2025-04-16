import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';

import { ActionConditionPipe } from './action-condition-pipe.pipe';

@NgModule({
  declarations: [ActionConditionPipe],
  imports: [CommonModule],
  exports: [ActionConditionPipe],
})
export class ActionConditionPipeModule {}
