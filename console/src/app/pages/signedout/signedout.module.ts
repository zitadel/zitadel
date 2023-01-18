import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyTooltipModule as MatTooltipModule } from '@angular/material/legacy-tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { SignedoutRoutingModule } from './signedout-routing.module';
import { SignedoutComponent } from './signedout.component';

@NgModule({
  declarations: [SignedoutComponent],
  imports: [CommonModule, SignedoutRoutingModule, MatButtonModule, MatTooltipModule, TranslateModule],
})
export default class SignedoutModule {}
