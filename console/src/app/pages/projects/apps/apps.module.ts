import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatButtonToggleModule } from '@angular/material/button-toggle';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatChipsModule } from '@angular/material/chips';
import { MatDialogModule } from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatMenuModule } from '@angular/material/menu';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSelectModule } from '@angular/material/select';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { CardModule } from 'src/app/modules/card/card.module';

import { AppCreateComponent } from './app-create/app-create.component';
import { AppDetailComponent } from './app-detail/app-detail.component';
import { AppSecretDialogComponent } from './app-secret-dialog/app-secret-dialog.component';
import { AppsRoutingModule } from './apps-routing.module';

@NgModule({
    declarations: [
        AppCreateComponent,
        AppDetailComponent,
        AppSecretDialogComponent,
    ],
    imports: [
        CommonModule,
        AppsRoutingModule,
        FormsModule,
        TranslateModule,
        ReactiveFormsModule,
        HasRoleModule,
        MatFormFieldModule,
        MatInputModule,
        MatMenuModule,
        MatChipsModule,
        MatIconModule,
        MatSelectModule,
        MatButtonToggleModule,
        MatButtonModule,
        MatProgressSpinnerModule,
        MatProgressBarModule,
        MatDialogModule,
        MatCheckboxModule,
        CardModule,
        MatTooltipModule,
        TranslateModule,
    ],
    exports: [TranslateModule],
})
export class AppsModule { }
