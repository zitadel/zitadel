import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatRadioModule } from '@angular/material/radio';
import { TranslateModule } from '@ngx-translate/core';

import { InfoSectionModule } from '../info-section/info-section.module';
import { ProjectPrivateLabelingDialogComponent } from './project-private-labeling-dialog.component';



@NgModule({
  declarations: [ProjectPrivateLabelingDialogComponent],
  imports: [
    CommonModule,
    TranslateModule,
    MatButtonModule,
    FormsModule,
    MatRadioModule,
    InfoSectionModule
  ],
})
export class ProjectPrivateLabelingDialogModule { }
