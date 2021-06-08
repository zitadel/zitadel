import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatSelectModule } from '@angular/material/select';
import { TranslateModule } from '@ngx-translate/core';
import { AvatarModule } from 'src/app/modules/avatar/avatar.module';
import { InputModule } from 'src/app/modules/input/input.module';

import { DetailFormComponent } from './detail-form.component';

@NgModule({
  declarations: [
    DetailFormComponent,
  ],
  imports: [
    AvatarModule,
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    TranslateModule,
    MatSelectModule,
    MatButtonModule,
    MatIconModule,
    TranslateModule,
    InputModule,
  ],
  exports: [
    DetailFormComponent,
  ],
})
export class DetailFormModule { }
