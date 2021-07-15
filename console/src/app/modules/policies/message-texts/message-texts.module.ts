import { TextFieldModule } from '@angular/cdk/text-field';
import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { DetailLayoutModule } from 'src/app/modules/detail-layout/detail-layout.module';
import { InputModule } from 'src/app/modules/input/input.module';
import { HasFeaturePipeModule } from 'src/app/pipes/has-feature-pipe/has-feature-pipe.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';

import { EditTextModule } from '../../edit-text/edit-text.module';
import { FormFieldModule } from '../../form-field/form-field.module';
import { InfoSectionModule } from '../../info-section/info-section.module';
import { LinksModule } from '../../links/links.module';
import { MessageTextsRoutingModule } from './message-texts-routing.module';
import { MessageTextsComponent } from './message-texts.component';

@NgModule({
  declarations: [MessageTextsComponent],
  imports: [
    MessageTextsRoutingModule,
    CommonModule,
    InfoSectionModule,
    ReactiveFormsModule,
    FormsModule,
    InputModule,
    FormFieldModule,
    EditTextModule,
    MatButtonModule,
    HasFeaturePipeModule,
    MatIconModule,
    HasRoleModule,
    HasRolePipeModule,
    MatTooltipModule,
    TranslateModule,
    MatTooltipModule,
    DetailLayoutModule,
    MatProgressSpinnerModule,
    LinksModule,
    TextFieldModule,
  ],
})
export class MessageTextsPolicyModule { }
