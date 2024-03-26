import {ChangeDetectorRef, computed, Injectable, Injector, signal, Signal, Type, WritableSignal} from '@angular/core';
import {BehaviorSubject, combineLatestWith, defer, from, Observable, of, shareReplay, switchMap, take} from 'rxjs';
import { filter, map, tap } from 'rxjs/operators';
import { EnvironmentService } from '../../../services/environment.service';
import { CopyUrl } from './provider-next.component';
import { ManagementService } from '../../../services/mgmt.service';
import { AdminService } from '../../../services/admin.service';
import { IDPOwnerType } from '../../../proto/generated/zitadel/idp_pb';
import { ToastService } from '../../../services/toast.service';
import { Data, ParamMap } from '@angular/router';
import { ActivateIdpService } from '../../../services/activate-idp.service';
import { PolicyComponentServiceType } from '../../policies/policy-component-types.enum';
import {toSignal} from "@angular/core/rxjs-interop";


@Injectable({
  providedIn: 'root',
})
export class ProviderNextServiceV2 {
  constructor(
    private env: EnvironmentService,
    private toast: ToastService,
    private addIdpSvc: ActivateIdpService,
    private injector: Injector,
  ) {}

  service(routeData$: Observable<Data>): Signal<ManagementService | AdminService> {
    const routeDataSignal = toSignal(routeData$);
    return computed(() => {
      const data = routeDataSignal();
      if (!data) {
        throw new Error('No Data');
      }
      switch (data['serviceType']) {
        case PolicyComponentServiceType.MGMT:
          return this.injector.get(ManagementService as Type<ManagementService>);
        case PolicyComponentServiceType.ADMIN:
          return this.injector.get(AdminService as Type<AdminService>);
        default:
          throw new Error('Unknown Service Type');
      }
    });
  }

  id(paramMap$: Observable<ParamMap>, justCreated: Signal<string>): Signal<string> {
    const paramMapSignal = toSignal(paramMap$);
    return computed(() => {
      const params = paramMapSignal();
      if (!params) {
        throw new Error('No Params');
      }
      const created = justCreated()
      if (created) {
        return created;
      }
      return params.get('id') || ''
    });
  }

  exists(id: Signal<string>): Signal<boolean> {
    return computed(() => !!id());
  }
  autofillLink(id: Signal<string>, link: string): Signal<string> {
    return computed(() => {
      const isID = id();
      return isID ? link : '';
    });
  }

  activateLink(
    id: Signal<string>,
    justActivated: Signal<boolean>,
    link: string,
    service: Signal<ManagementService | AdminService>,
  ): Observable<string> {
    console.log("aufgerufen")
    return defer(() => {
      const idValue = id();
      const activated = justActivated();
      console.log("id", idValue, "activated", activated)
      if (!idValue || activated) {
        return of('');
      }
      console.log("fetching")
      return from(service().getLoginPolicy()).pipe(
        map((policy) => {
          return policy.policy?.idpsList.find((idp) => idp.idpId === idValue) ? '' : link;
        }),
        tap((policy) => console.log("fetched", policy)),
      );
    })
  }

  callbackUrls(): Signal<CopyUrl[]> {
    const envSignal = toSignal(this.env.env);
    return computed(() => [
      {
        label: 'ZITADEL Callback URL',
        url: `${envSignal()?.issuer}/ui/login/login/externalidp/callback`,
      },
    ]);
  }

  expandWhatNow(
    id: Signal<string | null>,
    activateLink$: Observable<string>,
    justCreated: Signal<string>,
  ): Signal<boolean> {
    const activateLinkSignal = toSignal(activateLink$);
    return computed(() => {
      const idValue = !!id();
      const link = !!activateLinkSignal();
      const created = !!justCreated();
      return !idValue || link || created;
    })
  }

  activate(
    id: Signal<string | null>,
    emitActivated: WritableSignal<boolean>,
    service: Signal<ManagementService | AdminService>,
  ): void {
      const idValue = id();
      if (!idValue) {
        throw new Error('No ID');
      }
      const serviceValue = service();
      if (!serviceValue) {
        throw new Error('No Service');
      }
      this.addIdpSvc.activateIdp(
        serviceValue,
        idValue,
        serviceValue instanceof AdminService ? IDPOwnerType.IDP_OWNER_TYPE_SYSTEM : IDPOwnerType.IDP_OWNER_TYPE_ORG,
      )      .subscribe({
        next: () => {
          this.toast.showInfo('POLICY.LOGIN_POLICY.PROVIDER_ADDED', true);
          emitActivated.set(true);
        },
        error: this.toast.showError,
      });
  }
}
