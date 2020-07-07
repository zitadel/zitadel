import { CommonModule } from '@angular/common';
import { NgModule, NO_ERRORS_SCHEMA } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatDialogModule } from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { CardModule } from 'src/app/modules/card/card.module';
import { ChangesModule } from 'src/app/modules/changes/changes.module';
import { MetaLayoutModule } from 'src/app/modules/meta-layout/meta-layout.module';
import { UserGrantsModule } from 'src/app/modules/user-grants/user-grants.module';
import { PipesModule } from 'src/app/pipes/pipes.module';

import { DetailFormModule } from '../detail-form/detail-form.module';
import { UserDetailRoutingModule } from './user-detail-routing.module';
import { UserDetailComponent } from './user-detail.component';
import { UserMfaComponent } from './user-mfa/user-mfa.component';


@NgModule({
    declarations: [
        UserDetailComponent,
        UserMfaComponent,
    ],
    imports: [
        CommonModule,
        UserDetailRoutingModule,
        ChangesModule,
        FormsModule,
        ReactiveFormsModule,
        DetailFormModule,
        MatDialogModule,
        MetaLayoutModule,
        PipesModule,
        MatFormFieldModule,
        UserGrantsModule,
        MatInputModule,
        MatButtonModule,
        MatIconModule,
        CardModule,
        MatProgressBarModule,
        MatTooltipModule,
        HasRoleModule,
        TranslateModule,
    ],
    schemas: [
        NO_ERRORS_SCHEMA, // used for metainfo
    ],
})
export class UserDetailModule { }
