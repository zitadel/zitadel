import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatDialogModule } from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatPaginatorModule } from '@angular/material/paginator';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTableModule } from '@angular/material/table';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { QRCodeModule } from 'angularx-qrcode';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { CardModule } from 'src/app/modules/card/card.module';
import { ChangesModule } from 'src/app/modules/changes/changes.module';
import { MetaLayoutModule } from 'src/app/modules/meta-layout/meta-layout.module';
import { SharedModule } from 'src/app/modules/shared/shared.module';
import { UserGrantsModule } from 'src/app/modules/user-grants/user-grants.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe.module';

import { AuthUserDetailComponent } from './auth-user-detail/auth-user-detail.component';
import { AuthUserMfaComponent } from './auth-user-detail/auth-user-mfa/auth-user-mfa.component';
import { CodeDialogComponent } from './auth-user-detail/code-dialog/code-dialog.component';
import { DialogOtpComponent } from './auth-user-detail/dialog-otp/dialog-otp.component';
import { ThemeSettingComponent } from './auth-user-detail/theme-setting/theme-setting.component';
import { DetailFormModule } from './detail-form/detail-form.module';
import { PasswordComponent } from './password/password.component';
import { UserDetailRoutingModule } from './user-detail-routing.module';
import { UserDetailComponent } from './user-detail/user-detail.component';
import { UserMfaComponent } from './user-detail/user-mfa/user-mfa.component';

@NgModule({
    declarations: [
        AuthUserDetailComponent,
        UserDetailComponent,
        DialogOtpComponent,
        AuthUserMfaComponent,
        UserMfaComponent,
        ThemeSettingComponent,
        PasswordComponent,
        CodeDialogComponent,
    ],
    imports: [
        UserDetailRoutingModule,
        ChangesModule,
        CommonModule,
        FormsModule,
        ReactiveFormsModule,
        DetailFormModule,
        MatDialogModule,
        QRCodeModule,
        MetaLayoutModule,
        HasRolePipeModule,
        MatFormFieldModule,
        UserGrantsModule,
        MatInputModule,
        MatButtonModule,
        MatIconModule,
        CardModule,
        MatProgressSpinnerModule,
        MatProgressBarModule,
        MatTooltipModule,
        HasRoleModule,
        TranslateModule,
        MatTableModule,
        MatPaginatorModule,
        SharedModule,
    ],
})
export class UserDetailModule { }
