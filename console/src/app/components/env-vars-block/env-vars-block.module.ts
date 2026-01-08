import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';

import { EnvVarsBlockComponent } from './env-vars-block.component';

@NgModule({
  declarations: [EnvVarsBlockComponent],
  imports: [CommonModule, MatButtonModule, MatIconModule, MatTooltipModule, TranslateModule],
  exports: [EnvVarsBlockComponent],
})
export class EnvVarsBlockModule {}
