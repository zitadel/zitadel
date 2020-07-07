import { CommonModule } from '@angular/common';
import { NgModule, NO_ERRORS_SCHEMA } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatDialogModule } from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { QRCodeModule } from 'angularx-qrcode';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { CardModule } from 'src/app/modules/card/card.module';
import { ChangesModule } from 'src/app/modules/changes/changes.module';
import { MetaLayoutModule } from 'src/app/modules/meta-layout/meta-layout.module';
import { PipesModule } from 'src/app/pipes/pipes.module';

import { DetailFormModule } from '../../modules/detail-form/detail-form.module';
import { PasswordModule } from '../password/password.module';
import { AuthUserDetailRoutingModule } from './auth-user-detail-routing.module';
import { AuthUserDetailComponent } from './auth-user-detail.component';
import { AuthUserMfaComponent } from './auth-user-mfa/auth-user-mfa.component';
import { CodeDialogComponent } from './code-dialog/code-dialog.component';
import { DialogOtpComponent } from './dialog-otp/dialog-otp.component';
import { ThemeSettingComponent } from './theme-setting/theme-setting.component';

@NgModule({
    declarations: [
        AuthUserDetailComponent,
        DialogOtpComponent,
        AuthUserMfaComponent,
        ThemeSettingComponent,
        CodeDialogComponent,
    ],
    imports: [
        CommonModule,
        AuthUserDetailRoutingModule,
        ChangesModule,
        PasswordModule,
        FormsModule,
        ReactiveFormsModule,
        DetailFormModule,
        MatDialogModule,
        QRCodeModule,
        MetaLayoutModule,
        PipesModule,
        MatProgressSpinnerModule,
        MatFormFieldModule,
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
export class AuthUserDetailModule { }
