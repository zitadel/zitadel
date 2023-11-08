import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatRadioModule } from '@angular/material/radio';
import { TranslateModule } from '@ngx-translate/core';

import { MatDialogModule } from '@angular/material/dialog';
import { InfoSectionModule } from '../info-section/info-section.module';
import { InputModule } from '../input/input.module';
import { ProjectRoleDetailDialogComponent } from './project-role-detail-dialog.component';

@NgModule({
  declarations: [ProjectRoleDetailDialogComponent],
  imports: [
    CommonModule,
    TranslateModule,
    InputModule,
    MatDialogModule,
    MatButtonModule,
    FormsModule,
    ReactiveFormsModule,
    MatRadioModule,
    InfoSectionModule,
  ],
  exports: [ProjectRoleDetailDialogComponent],
})
export class ProjectRoleDetailDialogModule {}
