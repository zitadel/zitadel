import { Injectable } from '@angular/core';
import { forkJoin, Observable, switchMap } from 'rxjs';
import { map } from 'rxjs/operators';
import { EnvironmentService, Environment } from '../../../services/environment.service';
import { TranslateService } from '@ngx-translate/core';
import { CopyUrl } from './provider-next.component';

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
    activateLink$: Observable<string>,
    instance: boolean,
    configureTitleI18nKey: string,
    configureDescriptionI18nKey: string,
    configureLink: string,
    autofillLink$: Observable<string>,
    copyUrls: (env: Environment) => CopyUrl[],
  ): Observable<any> {
    return forkJoin([
      this.env.env,
      this.translateSvc.get(configureTitleI18nKey, { provider: providerName }),
      this.translateSvc.get(configureDescriptionI18nKey, { provider: providerName }),
    ]).pipe(
      switchMap(([environment, title, description]) =>
        autofillLink$.pipe(
          switchMap((autofillLink) => activateLink$.pipe(
            map((activateLink) => ({
              copyUrls: copyUrls(environment),
              configureTitle: title as string,
              configureDescription: description as string,
              configureLink: configureLink,
              autofillLink: autofillLink,
              activateLink: activateLink,
              instance: instance,
            })),
          )),
        ),
      ))
  }

  callbackUrls(): Observable<CopyUrl[]> {
    return this.env.env.pipe(
      map((env) => [
        {
          label: 'ZITADEL Callback URL',
          url: `${env.issuer}/ui/login/login/externalidp/callback`,
        },
      ]),
    );
  }
}
