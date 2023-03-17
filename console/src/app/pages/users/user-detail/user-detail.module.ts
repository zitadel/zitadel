import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyCheckboxModule as MatCheckboxModule } from '@angular/material/legacy-checkbox';
import { MatLegacyDialogModule as MatDialogModule } from '@angular/material/legacy-dialog';
import { MatLegacyMenuModule as MatMenuModule } from '@angular/material/legacy-menu';
import { MatLegacyProgressSpinnerModule as MatProgressSpinnerModule } from '@angular/material/legacy-progress-spinner';
import { MatLegacyTableModule as MatTableModule } from '@angular/material/legacy-table';
import { MatLegacyTabsModule as MatTabsModule } from '@angular/material/legacy-tabs';
import { MatLegacyTooltipModule as MatTooltipModule } from '@angular/material/legacy-tooltip';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { CopyToClipboardModule } from 'src/app/directives/copy-to-clipboard/copy-to-clipboard.module';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { MemberCreateDialogModule } from 'src/app/modules/add-member-dialog/member-create-dialog.module';
import { CardModule } from 'src/app/modules/card/card.module';
import { ChangesModule } from 'src/app/modules/changes/changes.module';
import { DetailLayoutModule } from 'src/app/modules/detail-layout/detail-layout.module';
import { InfoSectionModule } from 'src/app/modules/info-section/info-section.module';
import { InputModule } from 'src/app/modules/input/input.module';
import { MachineKeysModule } from 'src/app/modules/machine-keys/machine-keys.module';
import { MembershipsTableModule } from 'src/app/modules/memberships-table/memberships-table.module';
import { MetaLayoutModule } from 'src/app/modules/meta-layout/meta-layout.module';
import { PaginatorModule } from 'src/app/modules/paginator/paginator.module';
import { PasswordComplexityViewModule } from 'src/app/modules/password-complexity-view/password-complexity-view.module';
import { PersonalAccessTokensModule } from 'src/app/modules/personal-access-tokens/personal-access-tokens.module';
import { RefreshTableModule } from 'src/app/modules/refresh-table/refresh-table.module';
import { ShowTokenDialogModule } from 'src/app/modules/show-token-dialog/show-token-dialog.module';
import { SidenavModule } from 'src/app/modules/sidenav/sidenav.module';
import { TableActionsModule } from 'src/app/modules/table-actions/table-actions.module';
import { TopViewModule } from 'src/app/modules/top-view/top-view.module';
import { UserGrantsModule } from 'src/app/modules/user-grants/user-grants.module';
import { WarnDialogModule } from 'src/app/modules/warn-dialog/warn-dialog.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe/has-role-pipe.module';
import { LocalizedDatePipeModule } from 'src/app/pipes/localized-date-pipe/localized-date-pipe.module';
import { TimestampToDatePipeModule } from 'src/app/pipes/timestamp-to-date-pipe/timestamp-to-date-pipe.module';

import { MatLegacySelectModule as MatSelectModule } from '@angular/material/legacy-select';
import { QRCodeModule } from 'angularx-qrcode';
import { MetadataModule } from 'src/app/modules/metadata/metadata.module';
import { CountryCallingCodesService } from 'src/app/services/country-calling-codes.service';
import { InfoRowModule } from '../../../modules/info-row/info-row.module';
import { AuthFactorDialogComponent } from './auth-user-detail/auth-factor-dialog/auth-factor-dialog.component';
import { AuthPasswordlessComponent } from './auth-user-detail/auth-passwordless/auth-passwordless.component';
import { DialogPasswordlessComponent } from './auth-user-detail/auth-passwordless/dialog-passwordless/dialog-passwordless.component';
import { AuthUserDetailComponent } from './auth-user-detail/auth-user-detail.component';
import { AuthUserMfaComponent } from './auth-user-detail/auth-user-mfa/auth-user-mfa.component';
import { CodeDialogComponent } from './auth-user-detail/code-dialog/code-dialog.component';
import { DialogU2FComponent } from './auth-user-detail/dialog-u2f/dialog-u2f.component';
import { EditDialogComponent } from './auth-user-detail/edit-dialog/edit-dialog.component';
import { ResendEmailDialogComponent } from './auth-user-detail/resend-email-dialog/resend-email-dialog.component';
import { ContactComponent } from './contact/contact.component';
import { DetailFormMachineModule } from './detail-form-machine/detail-form-machine.module';
import { DetailFormModule } from './detail-form/detail-form.module';
import { ExternalIdpsComponent } from './external-idps/external-idps.component';
import { PasswordComponent } from './password/password.component';
import { PhoneDetailComponent } from './phone-detail/phone-detail.component';
import { MachineSecretDialogComponent } from './user-detail/machine-secret-dialog/machine-secret-dialog.component';
import { PasswordlessComponent } from './user-detail/passwordless/passwordless.component';
import { UserDetailComponent } from './user-detail/user-detail.component';
import { UserMfaComponent } from './user-detail/user-mfa/user-mfa.component';

@NgModule({
  declarations: [
    AuthUserDetailComponent,
    UserDetailComponent,
    EditDialogComponent,
    AuthUserMfaComponent,
    AuthPasswordlessComponent,
    UserMfaComponent,
    PasswordlessComponent,
    PasswordComponent,
    CodeDialogComponent,
    ExternalIdpsComponent,
    ContactComponent,
    ResendEmailDialogComponent,
    DialogU2FComponent,
    DialogPasswordlessComponent,
    AuthFactorDialogComponent,
    PhoneDetailComponent,
    MachineSecretDialogComponent,
  ],
  providers: [CountryCallingCodesService],
  imports: [
    ChangesModule,
    CommonModule,
    SidenavModule,
    MatTabsModule,
    FormsModule,
    ReactiveFormsModule,
    MembershipsTableModule,
    DetailFormModule,
    DetailFormMachineModule,
    WarnDialogModule,
    MatDialogModule,
    QRCodeModule,
    ShowTokenDialogModule,
    MetaLayoutModule,
    MatCheckboxModule,
    MetadataModule,
    TopViewModule,
    HasRolePipeModule,
    UserGrantsModule,
    MatButtonModule,
    PersonalAccessTokensModule,
    MatIconModule,
    CardModule,
    MatProgressSpinnerModule,
    MatTooltipModule,
    HasRoleModule,
    TranslateModule,
    MatTableModule,
    InfoRowModule,
    PaginatorModule,
    MatMenuModule,
    RouterModule,
    RefreshTableModule,
    CopyToClipboardModule,
    DetailLayoutModule,
    TableActionsModule,
    PasswordComplexityViewModule,
    MemberCreateDialogModule,
    TimestampToDatePipeModule,
    LocalizedDatePipeModule,
    InputModule,
    MachineKeysModule,
    InfoSectionModule,
    MatSelectModule,
  ],
})
export class UserDetailModule {}
