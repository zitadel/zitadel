import { Injectable } from '@angular/core';
import { forkJoin, Observable, switchMap } from 'rxjs';
import { map } from 'rxjs/operators';
import { EnvironmentService, Environment } from '../../../services/environment.service';
import { TranslateService } from '@ngx-translate/core';
import { CopyUrl, Next } from './provider-next.component';

@Injectable({
  providedIn: 'root',
})
export class ProviderNextService {
  constructor(
    private env: EnvironmentService,
    private translateSvc: TranslateService,
  ) {}

  next(
    providerName: string,
    configureTitleI18nKey: string,
    configureDescriptionI18nKey: string,
    configureLink: string,
    autofillLink$: Observable<string>,
    copyUrls: (env: Environment) => CopyUrl[],
  ): Observable<Next> {
    return forkJoin([
      this.env.env,
      this.translateSvc.get(configureTitleI18nKey, { provider: providerName }),
      this.translateSvc.get(configureDescriptionI18nKey, { provider: providerName }),
    ]).pipe(
      switchMap(([environment, title, description]) =>
        autofillLink$.pipe(
          map((autofillLink) => ({
            copyUrls: copyUrls(environment),
            configureTitle: title as string,
            configureDescription: description as string,
            configureLink: configureLink,
            autofillLink: autofillLink,
          })),
        ),
      ),
    );
  }

  callbackUrls(env: Environment): CopyUrl[] {
    return [
      {
        label: 'ZITADEL Callback URL',
        url: `${env.issuer}/ui/login/login/externalidp/callback`,
      },
    ];
  }
}
