import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { RouterModule } from '@angular/router';

import { DetailLayoutComponent } from './detail-layout.component';

@NgModule({
  declarations: [DetailLayoutComponent],
  imports: [CommonModule, MatIconModule, MatButtonModule, RouterModule],
  exports: [DetailLayoutComponent],
})
export class DetailLayoutModule {}
