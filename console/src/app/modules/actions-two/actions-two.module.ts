import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActionsTwoActionsComponent } from './actions-two-actions/actions-two-actions.component';
import { ActionsTwoTargetsComponent } from './actions-two-targets/actions-two-targets.component';

@NgModule({
  declarations: [ActionsTwoActionsComponent, ActionsTwoTargetsComponent],
  imports: [CommonModule],
  exports: [ActionsTwoActionsComponent, ActionsTwoTargetsComponent],
})
export class ActionsTwoModule {}
