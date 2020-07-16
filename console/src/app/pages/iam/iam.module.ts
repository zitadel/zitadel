import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatAutocompleteModule } from '@angular/material/autocomplete';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatChipsModule } from '@angular/material/chips';
import { MatDialogModule } from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatPaginatorModule } from '@angular/material/paginator';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSortModule } from '@angular/material/sort';
import { MatTableModule } from '@angular/material/table';
import { MatTabsModule } from '@angular/material/tabs';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { CardModule } from 'src/app/modules/card/card.module';
import { ChangesModule } from 'src/app/modules/changes/changes.module';
import { MetaLayoutModule } from 'src/app/modules/meta-layout/meta-layout.module';

import { IamContributorsModule } from './iam-contributors/iam-contributors.module';
import { IamRoutingModule } from './iam-routing.module';
import { IamViewsComponent } from './iam-views/iam-views.component';
import { IamComponent } from './iam.component';



@NgModule({
    declarations: [IamComponent, IamViewsComponent],
    imports: [
        CommonModule,
        IamRoutingModule,
        ChangesModule,
        CardModule,
        MatAutocompleteModule,
        MatChipsModule,
        MatButtonModule,
        HasRoleModule,
        MatCheckboxModule,
        MetaLayoutModule,
        MatIconModule,
        MatTabsModule,
        MatTableModule,
        MatPaginatorModule,
        MatFormFieldModule,
        MatSortModule,
        MatTooltipModule,
        ReactiveFormsModule,
        MatProgressSpinnerModule,
        FormsModule,
        TranslateModule,
        MatDialogModule,
        IamContributorsModule,
    ],
})
export class IamModule { }
