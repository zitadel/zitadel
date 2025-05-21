import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';

import { TypeSafeCellDefDirective } from './type-safe-cell-def.directive';

@NgModule({
  declarations: [TypeSafeCellDefDirective],
  imports: [CommonModule],
  exports: [TypeSafeCellDefDirective],
})
export class TypeSafeCellDefModule {}
