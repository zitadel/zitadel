import { ChangeDetectionStrategy, Component, computed, inject, input } from '@angular/core';
import { TranslatePipe } from '@ngx-translate/core';
import { injectQuery } from '@tanstack/angular-query-experimental';
import { RouterLink } from '@angular/router';
import { MatButtonModule } from '@angular/material/button';

import { InfoSectionModule } from '../info-section/info-section.module';
import { NewSettingsService } from 'src/app/services/new-settings.service';

// States the restrictions the ZITADEL operator defined for the active email provider,
// e.g. a default provider shared by multiple instances.
// Nothing is rendered if the active provider is not restricted or the settings cannot be read.
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

  private readonly settingsService = inject(NewSettingsService);
  // the general settings are readable by organization admins, which customize message texts as well
  private readonly settingsQuery = injectQuery(() => this.settingsService.getGeneralSettingsQueryOptions());

  // a failed reload keeps the previous data, which must not keep a stale banner
  protected readonly restrictions = computed(() =>
    this.settingsQuery.isError() ? undefined : this.settingsQuery.data()?.emailProviderRestrictions,
  );

  protected readonly limit = computed(() => {
    const limit = this.restrictions()?.sendingLimit;
    if (!limit) {
      return undefined;
    }
    return { count: limit.count, duration: formatWindow(Number(limit.window?.seconds ?? 0)) };
  });

  protected readonly show = computed(() => {
    const restrictions = this.restrictions();
    if (!restrictions) {
      return false;
    }
    return this.mode() === 'provider' || restrictions.customHtmlRestricted;
  });
}

// formatWindow states the window in the largest unit representing it exactly, e.g. 24 h, 5 min or 90 s.
function formatWindow(seconds: number): string {
  if (seconds % 3600 === 0) {
    return `${seconds / 3600} h`;
  }
  if (seconds % 60 === 0) {
    return `${seconds / 60} min`;
  }
  return `${seconds} s`;
}
