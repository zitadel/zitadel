import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';

import { HasFeatureDirective } from './has-feature.directive';



@NgModule({
  declarations: [
    HasFeatureDirective,
  ],
  imports: [
    CommonModule,
  ],
  exports: [
    HasFeatureDirective,
  ],
})
export class HasFeatureModule { }
