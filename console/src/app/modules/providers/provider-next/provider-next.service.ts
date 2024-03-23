import { Injectable, Injector, Type } from '@angular/core';
import { BehaviorSubject, combineLatestWith, forkJoin, from, Observable, of, shareReplay, switchMap, take } from 'rxjs';
import { filter, map, tap } from 'rxjs/operators';
import { Environment, EnvironmentService } from '../../../services/environment.service';
import { TranslateService } from '@ngx-translate/core';
import { CopyUrl } from './provider-next.component';
import { ManagementService } from '../../../services/mgmt.service';
import { AdminService } from '../../../services/admin.service';
import { IDPOwnerType } from '../../../proto/generated/zitadel/idp_pb';
import { ToastService } from '../../../services/toast.service';
import { Data, ParamMap } from '@angular/router';
import { ActivateIdpService } from '../../../services/activate-idp.service';
import { PolicyComponentServiceType } from '../../policies/policy-component-types.enum';

@Injectable({
  providedIn: 'root',
})
export class ProviderNextService {
  constructor(
    private env: EnvironmentService,
    private translateSvc: TranslateService,
    private toast: ToastService,
    private addIdpSvc: ActivateIdpService,
    private injector: Injector,
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
          switchMap((autofillLink) =>
            activateLink$.pipe(
              map((activateLink) => ({
                copyUrls: copyUrls(environment),
                configureTitle: title as string,
                configureDescription: description as string,
                configureLink: configureLink,
                autofillLink: autofillLink,
                activateLink: activateLink,
                instance: instance,
              })),
            ),
          ),
        ),
      ),
    );
  }

  service(routeData: Observable<Data>): Observable<ManagementService | AdminService> {
    return routeData.pipe(
      map((data) => {
        switch (data['serviceType']) {
          case PolicyComponentServiceType.MGMT:
            return this.injector.get(ManagementService as Type<ManagementService>);
          case PolicyComponentServiceType.ADMIN:
            return this.injector.get(AdminService as Type<AdminService>);
          default:
            throw new Error('Unknown Service Type');
        }
      }),
      shareReplay(1),
    );
  }

  id(paramMap: Observable<ParamMap>, justCreated$: Observable<string>): Observable<string | null> {
    return paramMap.pipe(
      // The ID observable should also emit when the IDP was just created
      combineLatestWith(justCreated$),
      map(([params, created]) => (created ? created : params.get('id'))),
      shareReplay(1),
      tap(console.info),
    );
  }

  exists(id$: Observable<string | null>): Observable<boolean> {
    return id$.pipe(
      map((id) => !!id),
      shareReplay(1),
      tap(console.info),
    );
  }
  autofillLink(id$: Observable<string | null>, link: string): Observable<string> {
    return id$.pipe(
      filter((id) => !!id),
      map(() => link),
      shareReplay(1),
    );
  }

  activateLink(
    id$: Observable<string | null>,
    justActivated$: Observable<boolean>,
    link: string,
    service$: Observable<ManagementService | AdminService>,
  ): Observable<string> {
    return id$.pipe(
      combineLatestWith(justActivated$, service$),
      // Because we also want to emit when the IDP is not active, we return an empty string if the IDP does not exist
      switchMap(([id, activated, service]) =>
        (!id || activated
          ? of(false)
          : from(service.getLoginPolicy()).pipe(map((policy) => !policy.policy?.idpsList.find((idp) => idp.idpId === id)))
        ).pipe(map((show) => (!show ? '' : link))),
      ),
      shareReplay(1),
    );
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

  expandWhatNow(
    id$: Observable<string | null>,
    activateLink$: Observable<string>,
    justCreated$: Observable<string>,
  ): Observable<boolean> {
    return id$.pipe(
      combineLatestWith(activateLink$, justCreated$),
      map(([id, activateLink, created]) => !id || activateLink || created),
      map((expand) => !!expand),
      shareReplay(1),
    );
  }

  activate(
    id$: Observable<string | null>,
    emitActivated$: BehaviorSubject<boolean>,
    service$: Observable<ManagementService | AdminService>,
  ): void {
    id$
      .pipe(
        combineLatestWith(service$),
        take(1),
        switchMap(([id, service]) => {
          if (!id) {
            throw new Error('No ID');
          }
          return this.addIdpSvc.activateIdp(
            service,
            id,
            service instanceof AdminService ? IDPOwnerType.IDP_OWNER_TYPE_SYSTEM : IDPOwnerType.IDP_OWNER_TYPE_ORG,
          );
        }),
      )
      .subscribe({
        next: () => {
          this.toast.showInfo('POLICY.LOGIN_POLICY.PROVIDER_ADDED', true);
          emitActivated$.next(true);
        },
        error: (error) => this.toast.showError(error),
      });
  }
}
