import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatDialogModule } from '@angular/material/dialog';
import { MatIconModule } from '@angular/material/icon';
import { MatPaginatorModule } from '@angular/material/paginator';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTableModule } from '@angular/material/table';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { QRCodeModule } from 'angularx-qrcode';
import { CopyToClipboardModule } from 'src/app/directives/copy-to-clipboard/copy-to-clipboard.module';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { MemberCreateDialogModule } from 'src/app/modules/add-member-dialog/member-create-dialog.module';
import { CardModule } from 'src/app/modules/card/card.module';
import { ChangesModule } from 'src/app/modules/changes/changes.module';
import { DetailLayoutModule } from 'src/app/modules/detail-layout/detail-layout.module';
import { InputModule } from 'src/app/modules/input/input.module';
import { MetaLayoutModule } from 'src/app/modules/meta-layout/meta-layout.module';
import { PasswordComplexityViewModule } from 'src/app/modules/password-complexity-view/password-complexity-view.module';
import { RefreshTableModule } from 'src/app/modules/refresh-table/refresh-table.module';
import { SharedModule } from 'src/app/modules/shared/shared.module';
import { UserGrantsModule } from 'src/app/modules/user-grants/user-grants.module';
import { WarnDialogModule } from 'src/app/modules/warn-dialog/warn-dialog.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';
import { LocalizedDatePipeModule } from 'src/app/pipes/localized-date-pipe/localized-date-pipe.module';
import { TimestampToDatePipeModule } from 'src/app/pipes/timestamp-to-date-pipe/timestamp-to-date-pipe.module';

import { AuthUserDetailComponent } from './auth-user-detail/auth-user-detail.component';
import { AuthUserMfaComponent } from './auth-user-detail/auth-user-mfa/auth-user-mfa.component';
import { CodeDialogComponent } from './auth-user-detail/code-dialog/code-dialog.component';
import { DialogOtpComponent } from './auth-user-detail/dialog-otp/dialog-otp.component';
import { EditDialogComponent } from './auth-user-detail/edit-dialog/edit-dialog.component';
import { ResendEmailDialogComponent } from './auth-user-detail/resend-email-dialog/resend-email-dialog.component';
import { ThemeSettingComponent } from './auth-user-detail/theme-setting/theme-setting.component';
import { ContactComponent } from './contact/contact.component';
import { DetailFormMachineModule } from './detail-form-machine/detail-form-machine.module';
import { DetailFormModule } from './detail-form/detail-form.module';
import { ExternalIdpsComponent } from './external-idps/external-idps.component';
import { AddKeyDialogModule } from './machine-keys/add-key-dialog/add-key-dialog.module';
import { MachineKeysComponent } from './machine-keys/machine-keys.component';
import { ShowKeyDialogModule } from './machine-keys/show-key-dialog/show-key-dialog.module';
import { MembershipsComponent } from './memberships/memberships.component';
import { PasswordComponent } from './password/password.component';
import { UserDetailRoutingModule } from './user-detail-routing.module';
import { UserDetailComponent } from './user-detail/user-detail.component';
import { UserMfaComponent } from './user-detail/user-mfa/user-mfa.component';

@NgModule({
    declarations: [
        AuthUserDetailComponent,
        UserDetailComponent,
        DialogOtpComponent,
        EditDialogComponent,
        AuthUserMfaComponent,
        UserMfaComponent,
        ThemeSettingComponent,
        PasswordComponent,
        CodeDialogComponent,
        MembershipsComponent,
        MachineKeysComponent,
        ExternalIdpsComponent,
        ContactComponent,
        ResendEmailDialogComponent,
    ],
    imports: [
        UserDetailRoutingModule,
        ChangesModule,
        CommonModule,
        FormsModule,
        ReactiveFormsModule,
        DetailFormModule,
        DetailFormMachineModule,
        WarnDialogModule,
        MatDialogModule,
        QRCodeModule,
        MetaLayoutModule,
        AddKeyDialogModule,
        ShowKeyDialogModule,
        MatCheckboxModule,
        HasRolePipeModule,
        UserGrantsModule,
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
        RefreshTableModule,
        CopyToClipboardModule,
        DetailLayoutModule,
        PasswordComplexityViewModule,
        MemberCreateDialogModule,
        TimestampToDatePipeModule,
        LocalizedDatePipeModule,
        InputModule,
    ],
})
export class UserDetailModule { }
