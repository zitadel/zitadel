import { Component } from '@angular/core';

import { PolicyComponentServiceType } from '../policy-component-types.enum';
import * as SMTPKnownProviders from '../../smtp-provider/known-smtp-providers-settings';
import { TranslatePipe } from '@ngx-translate/core';
import { SMTPTableModule } from '../../smtp-table/smtp-table.module';
import { RouterLink } from '@angular/router';
import { KeyValuePipe, TitleCasePipe } from '@angular/common';
import { MatIcon } from '@angular/material/icon';

@Component({
  selector: 'cnsl-notification-smtp-provider',
  templateUrl: './notification-smtp-provider.component.html',
  styleUrls: ['./notification-smtp-provider.component.scss'],
  imports: [TranslatePipe, SMTPTableModule, RouterLink, TitleCasePipe, MatIcon, KeyValuePipe],
})
export class NotificationSMTPProviderComponent {
  protected readonly PolicyComponentServiceType = PolicyComponentServiceType;
  protected readonly providers = { ...SMTPKnownProviders, generic: { description: 'generic' } } as const;
}
