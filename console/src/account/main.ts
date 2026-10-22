import { provideZoneChangeDetection } from '@angular/core';
import { platformBrowserDynamic } from '@angular/platform-browser-dynamic';
import { AccountModule } from './account.module';
platformBrowserDynamic()
  // Native components update ordinary fields after RPCs and rely on Zone change detection.
  .bootstrapModule(AccountModule, { applicationProviders: [provideZoneChangeDetection()] })
  .catch(() => {
    document.body.textContent = 'Account unavailable / Conta indisponível';
  });
