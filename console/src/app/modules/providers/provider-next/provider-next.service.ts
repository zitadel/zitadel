import { Injectable, Injector, Type } from '@angular/core';
import { BehaviorSubject, combineLatestWith, defer, from, Observable, of, shareReplay, switchMap, take } from 'rxjs';
import { catchError, filter, map, timeout } from 'rxjs/operators';
import { EnvironmentService } from 'src/app/services/environment.service';
import { CopyUrl } from './provider-next.component';
import { ManagementService } from 'src/app/services/mgmt.service';
import { AdminService } from 'src/app/services/admin.service';
import { IDPOwnerType } from 'src/app/proto/generated/zitadel/idp_pb';
import { ToastService } from 'src/app/services/toast.service';
import { Data, ParamMap } from '@angular/router';
import { LoginPolicyService } from 'src/app/services/login-policy.service';
import { PolicyComponentServiceType } from '../../policies/policy-component-types.enum';
import { NewFeatureService } from 'src/app/services/new-feature.service';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';

@Injectable({
  providedIn: 'root',
})
export class ProviderNextService {
  private readonly loginV2BaseUri$: Observable<string | undefined>;
  constructor(
    private readonly env: EnvironmentService,
    private readonly toast: ToastService,
    private readonly loginPolicySvc: LoginPolicyService,
    private readonly injector: Injector,
    private readonly featureService: NewFeatureService,
  ) {
    this.loginV2BaseUri$ = this.getLoginV2BaseUri();
  }

  private getLoginV2BaseUri(): Observable<string | undefined> {
    return defer(() => this.featureService.getInstanceFeatures()).pipe(
      timeout(1000),
      // we try to load the features if this fails or takes too long we just assume loginV2 is not available
      catchError(() => of({ loginV2: undefined })),
      map((features) => features?.loginV2?.baseUri),
      // we only try this once as the backup plan is not too bad
      // and in most cases this will work
      shareReplay({ refCount: false, bufferSize: 1 }),
      takeUntilDestroyed(),
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
      combineLatestWith(this.loginV2BaseUri$),
      map(([env, loginV2BaseUri]) => [
        {
          label: 'Login V1 Callback URL',
          url: `${env.issuer}/ui/login/login/externalidp/callback`,
        },
        {
          label: 'Login V2 Callback URL',
          // if we don't have a loginV2BaseUri we provide a placeholder url so the user knows what to fill in
          // this is not ideal but better than nothing
          url: loginV2BaseUri ? `${loginV2BaseUri}idps/callback` : '{LOGIN V2 Hostname}/idps/callback',
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
