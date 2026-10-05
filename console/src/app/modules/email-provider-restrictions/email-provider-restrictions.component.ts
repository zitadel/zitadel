import { ChangeDetectionStrategy, Component, computed, inject, input } from '@angular/core';
import { TranslatePipe } from '@ngx-translate/core';
import { injectQuery } from '@tanstack/angular-query-experimental';
import { RouterLink } from '@angular/router';
import { MatButtonModule } from '@angular/material/button';

import { InfoSectionModule } from '../info-section/info-section.module';
import { NewAdminService } from 'src/app/services/new-admin.service';

// States the restrictions the operator of the system defined for the active email provider,
// e.g. a default provider shared by multiple instances.
// Nothing is rendered if the active provider is not restricted or cannot be read.
@Component({
  selector: 'cnsl-email-provider-restrictions',
  templateUrl: './email-provider-restrictions.component.html',
  styleUrls: ['./email-provider-restrictions.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [InfoSectionModule, TranslatePipe, RouterLink, MatButtonModule],
})
export class EmailProviderRestrictionsComponent {
  // 'provider' states all restrictions, 'texts' only the one concerning custom message texts
  public readonly mode = input<'provider' | 'texts'>('provider');

  private readonly adminService = inject(NewAdminService);
  private readonly providerQuery = injectQuery(() => this.adminService.getEmailProviderQueryOptions());

  protected readonly restrictions = computed(() => this.providerQuery.data()?.config?.restrictions);

  protected readonly limit = computed(() => {
    const limit = this.restrictions()?.sendingLimit;
    if (!limit) {
      return undefined;
    }
    const minutes = Math.round(Number(limit.window?.seconds ?? 0) / 60);
    // whole hours read better than minutes
    return minutes % 60 === 0
      ? { key: 'SMTP.RESTRICTIONS.LIMIT_HOURS', params: { count: limit.count, duration: minutes / 60 } }
      : { key: 'SMTP.RESTRICTIONS.LIMIT_MINUTES', params: { count: limit.count, duration: minutes } };
  });

  protected readonly show = computed(() => {
    const restrictions = this.restrictions();
    if (!restrictions) {
      return false;
    }
    return this.mode() === 'provider' || restrictions.customHtmlRestricted;
  });
}
