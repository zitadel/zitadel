import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { SharedModule } from 'src/app/modules/shared/shared.module';

import { SignedoutRoutingModule } from './signedout-routing.module';
import { SignedoutComponent } from './signedout.component';

@NgModule({
  declarations: [SignedoutComponent],
  imports: [CommonModule, SignedoutRoutingModule, MatButtonModule, MatTooltipModule, TranslateModule, SharedModule],
})
export class SignedoutModule {}
