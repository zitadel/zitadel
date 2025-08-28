import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';

import { HasRoleDirective } from './has-role.directive';

@NgModule({
  declarations: [HasRoleDirective],
  imports: [CommonModule],
  exports: [HasRoleDirective],
})
export class HasRoleModule {}
