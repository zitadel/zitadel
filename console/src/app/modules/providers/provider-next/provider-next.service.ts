import { Injectable, Injector, Type } from '@angular/core';
import { BehaviorSubject, combineLatestWith, from, Observable, of, shareReplay, switchMap, take } from 'rxjs';
import { filter, map, tap } from 'rxjs/operators';
import { EnvironmentService } from '../../../services/environment.service';
import { CopyUrl } from './provider-next.component';
import { ManagementService } from '../../../services/mgmt.service';
import { AdminService } from '../../../services/admin.service';
import { IDPOwnerType } from '../../../proto/generated/zitadel/idp_pb';
import { ToastService } from '../../../services/toast.service';
import { Data, ParamMap } from '@angular/router';
import { LoginPolicyService } from '../../../services/login-policy.service';
import { PolicyComponentServiceType } from '../../policies/policy-component-types.enum';

@Injectable({
  providedIn: 'root',
})
export class ProviderNextService {
  constructor(
    private env: EnvironmentService,
    private toast: ToastService,
    private loginPolicySvc: LoginPolicyService,
    private injector: Injector,
  ) {}

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
    );
  }

  exists(id$: Observable<string | null>): Observable<boolean> {
    return id$.pipe(
      map((id) => !!id),
      shareReplay(1),
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
          return this.loginPolicySvc.activateIdp(
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
