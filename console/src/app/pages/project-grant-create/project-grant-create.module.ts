import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatAutocompleteModule } from '@angular/material/autocomplete';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatChipsModule } from '@angular/material/chips';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { CardModule } from 'src/app/modules/card/card.module';

import { ProjectRolesModule } from '../../modules/project-roles/project-roles.module';
import { ProjectGrantCreateRoutingModule } from './project-grant-create-routing.module';
import { ProjectGrantCreateComponent } from './project-grant-create.component';

@NgModule({
    declarations: [ProjectGrantCreateComponent],
    imports: [
        ProjectGrantCreateRoutingModule,
        CommonModule,
        MatAutocompleteModule,
        MatChipsModule,
        MatButtonModule,
        MatFormFieldModule,
        MatInputModule,
        CardModule,
        MatCheckboxModule,
        ProjectRolesModule,
        MatIconModule,
        MatTooltipModule,
        ReactiveFormsModule,
        MatProgressSpinnerModule,
        FormsModule,
        TranslateModule,
    ],
})
export class ProjectGrantCreateModule { }
