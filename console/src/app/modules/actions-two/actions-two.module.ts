import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActionsTwoActionsComponent } from './actions-two-actions/actions-two-actions.component';
import { ActionsTwoTargetsComponent } from './actions-two-targets/actions-two-targets.component';
import { ActionsTwoRoutingModule } from './actions-two-routing.module';
import { TranslateModule } from '@ngx-translate/core';
import { MatButtonModule } from '@angular/material/button';

@NgModule({
  declarations: [ActionsTwoActionsComponent, ActionsTwoTargetsComponent],
  imports: [CommonModule, MatButtonModule, ActionsTwoRoutingModule, TranslateModule],
  exports: [ActionsTwoActionsComponent, ActionsTwoTargetsComponent],
})
export default class ActionsTwoModule {}
