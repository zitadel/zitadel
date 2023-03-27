import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacySelectModule as MatSelectModule } from '@angular/material/legacy-select';
import { MatLegacyTooltipModule as MatTooltipModule } from '@angular/material/legacy-tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { DropzoneModule } from 'src/app/directives/dropzone/dropzone.module';
import { AvatarModule } from 'src/app/modules/avatar/avatar.module';
import { InputModule } from 'src/app/modules/input/input.module';

import { DetailFormComponent } from './detail-form.component';
import { ProfilePictureComponent } from './profile-picture/profile-picture.component';

@NgModule({
  declarations: [DetailFormComponent, ProfilePictureComponent],
  imports: [
    DropzoneModule,
    AvatarModule,
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    TranslateModule,
    MatSelectModule,
    MatButtonModule,
    MatTooltipModule,
    MatIconModule,
    TranslateModule,
    InputModule,
  ],
  exports: [DetailFormComponent],
})
export class DetailFormModule {}
