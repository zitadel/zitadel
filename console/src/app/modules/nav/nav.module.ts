import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';

import { NavComponent } from './nav.component';

@NgModule({
  declarations: [
    NavComponent
  ],
  imports: [
    CommonModule,
    FormsModule,
    TranslateModule,
    RouterModule,
  ],
  exports: [
    NavComponent,
  ]
})
export class NavModule { }
