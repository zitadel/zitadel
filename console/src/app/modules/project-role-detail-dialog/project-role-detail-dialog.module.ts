import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyRadioModule as MatRadioModule } from '@angular/material/legacy-radio';
import { TranslateModule } from '@ngx-translate/core';

import { InfoSectionModule } from '../info-section/info-section.module';
import { InputModule } from '../input/input.module';
import { ProjectRoleDetailDialogComponent } from './project-role-detail-dialog.component';

@NgModule({
  declarations: [ProjectRoleDetailDialogComponent],
  imports: [
    CommonModule,
    TranslateModule,
    InputModule,
    MatButtonModule,
    FormsModule,
    ReactiveFormsModule,
    MatRadioModule,
    InfoSectionModule,
  ],
  exports: [ProjectRoleDetailDialogComponent],
})
export class ProjectRoleDetailDialogModule {}
