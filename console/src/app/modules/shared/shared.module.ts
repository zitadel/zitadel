import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { QuicklinkModule } from 'ngx-quicklink';

@NgModule({
  declarations: [],
  imports: [CommonModule, QuicklinkModule],
  exports: [QuicklinkModule],
})
export class SharedModule {}
