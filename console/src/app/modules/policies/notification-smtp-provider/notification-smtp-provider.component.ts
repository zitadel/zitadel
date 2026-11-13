import { Component, ViewChild } from '@angular/core';

import { PolicyComponentServiceType } from '../policy-component-types.enum';
import * as SMTPKnownProviders from '../../smtp-provider/known-smtp-providers-settings';
import { TranslatePipe } from '@ngx-translate/core';
import { SMTPTableModule } from '../../smtp-table/smtp-table.module';
import { SMTPTableComponent } from '../../smtp-table/smtp-table.component';
import { RouterLink } from '@angular/router';
import { KeyValuePipe, TitleCasePipe } from '@angular/common';
import { MatIcon } from '@angular/material/icon';
import { MatDialog } from '@angular/material/dialog';
import {
  HttpProviderDialogComponent,
  HttpProviderDialogData,
} from '../../http-provider-dialog/http-provider-dialog.component';

@Component({
  selector: 'cnsl-notification-smtp-provider',
  templateUrl: './notification-smtp-provider.component.html',
  styleUrls: ['./notification-smtp-provider.component.scss'],
  imports: [TranslatePipe, SMTPTableModule, RouterLink, TitleCasePipe, MatIcon, KeyValuePipe],
})
export class NotificationSMTPProviderComponent {
  @ViewChild(SMTPTableComponent) private table?: SMTPTableComponent;

  protected readonly PolicyComponentServiceType = PolicyComponentServiceType;
  protected readonly providers = { ...SMTPKnownProviders, generic: { description: 'generic' } } as const;

  constructor(private readonly dialog: MatDialog) {}

  protected addHttpProvider(): void {
    this.dialog
      .open<HttpProviderDialogComponent, HttpProviderDialogData, boolean>(HttpProviderDialogComponent, {
        data: { type: 'email' },
        width: '500px',
      })
      .afterClosed()
      .subscribe((saved) => {
        if (saved) {
          this.table?.refreshPage();
        }
      });
  }
}
