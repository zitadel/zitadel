import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';
import { OidcWebKeysComponent } from './oidc-webkeys.component';
import { RefreshTableModule } from 'src/app/modules/refresh-table/refresh-table.module';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatTableModule } from '@angular/material/table';
import { MatMenuModule } from '@angular/material/menu';
import { TableActionsModule } from 'src/app/modules/table-actions/table-actions.module';
import { MatButtonModule } from '@angular/material/button';
import { ActionKeysModule } from 'src/app/modules/action-keys/action-keys.module';
import { MatIconModule } from '@angular/material/icon';
import { FormFieldModule } from 'src/app/modules/form-field/form-field.module';
import { MatSelectModule } from '@angular/material/select';
import { ReactiveFormsModule } from '@angular/forms';
import { CardModule } from 'src/app/modules/card/card.module';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { OidcWebKeysCreateComponent } from './oidc-webkeys-create/oidc-webkeys-create.component';
import { OidcWebKeysTableComponent } from './oidc-webkeys-table/oidc-webkeys-table.component';
import { TimestampToDatePipeModule } from 'src/app/pipes/timestamp-to-date-pipe/timestamp-to-date-pipe.module';
import { LocalizedDatePipeModule } from 'src/app/pipes/localized-date-pipe/localized-date-pipe.module';
import { OidcWebKeysInactiveTableComponent } from './oidc-webkeys-inactive-table/oidc-webkeys-inactive-table.component';
import { TypeSafeCellDefDirective } from './type-safe-cell-def.directive';
import { TimestampToDatePipe } from '../../../pipes/timestamp-to-date-pipe/timestamp-to-date.pipe';
import { MatTooltipModule } from '@angular/material/tooltip';

@NgModule({
  declarations: [
    OidcWebKeysComponent,
    OidcWebKeysCreateComponent,
    OidcWebKeysTableComponent,
    OidcWebKeysInactiveTableComponent,
    TypeSafeCellDefDirective,
  ],
  providers: [TimestampToDatePipe],
  imports: [
    CommonModule,
    TranslateModule,
    RefreshTableModule,
    MatCheckboxModule,
    MatTableModule,
    MatMenuModule,
    TableActionsModule,
    MatButtonModule,
    ActionKeysModule,
    MatIconModule,
    FormFieldModule,
    MatSelectModule,
    ReactiveFormsModule,
    CardModule,
    MatProgressSpinnerModule,
    TimestampToDatePipeModule,
    LocalizedDatePipeModule,
    MatTooltipModule,
  ],
  exports: [OidcWebKeysComponent],
})
export class OidcWebkeysModule {}
