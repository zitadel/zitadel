import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatTooltipModule } from '@angular/material/tooltip';

import { CardComponent } from './card.component';

@NgModule({
  declarations: [CardComponent],
  imports: [CommonModule, MatIconModule, MatButtonModule, MatTooltipModule],
  exports: [CardComponent],
})
export class CardModule {}
