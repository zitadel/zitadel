import { platformBrowserDynamic } from '@angular/platform-browser-dynamic';
import { AccountModule } from './account.module';
platformBrowserDynamic()
  .bootstrapModule(AccountModule)
  .catch(() => {
    document.body.textContent = 'Account unavailable / Conta indisponível';
  });
