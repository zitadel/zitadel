import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActionsTwoActionsComponent } from './actions-two-actions/actions-two-actions.component';
import { ActionsTwoTargetsComponent } from './actions-two-targets/actions-two-targets.component';
import { ActionsTwoRoutingModule } from './actions-two-routing.module';

@NgModule({
  declarations: [ActionsTwoActionsComponent, ActionsTwoTargetsComponent],
  imports: [CommonModule, ActionsTwoRoutingModule],
  exports: [ActionsTwoActionsComponent, ActionsTwoTargetsComponent],
})
export default class ActionsTwoModule {}
