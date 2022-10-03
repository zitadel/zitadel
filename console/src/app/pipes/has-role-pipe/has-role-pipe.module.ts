import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';

import { HasRolePipe } from './has-role.pipe';

@NgModule({
  declarations: [HasRolePipe],
  imports: [CommonModule],
  exports: [HasRolePipe],
})
export class HasRolePipeModule {}
