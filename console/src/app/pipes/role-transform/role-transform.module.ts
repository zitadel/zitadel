import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';

import { RoleTransformPipe } from './role-transform.pipe';

@NgModule({
  declarations: [RoleTransformPipe],
  imports: [CommonModule],
  exports: [RoleTransformPipe],
})
export class RoleTransformPipeModule {}
