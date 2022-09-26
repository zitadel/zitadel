import { LayoutModule } from '@angular/cdk/layout';
import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';

import { MetaLayoutComponent } from './meta-layout.component';

@NgModule({
  declarations: [MetaLayoutComponent],
  imports: [CommonModule, MatButtonModule, LayoutModule],
  exports: [MetaLayoutComponent],
})
export class MetaLayoutModule {}
