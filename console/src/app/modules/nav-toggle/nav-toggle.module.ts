import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';

import { NavToggleComponent } from './nav-toggle.component';

@NgModule({
  declarations: [NavToggleComponent],
  imports: [CommonModule, RouterModule],
  exports: [NavToggleComponent],
})
export class NavToggleModule {}
