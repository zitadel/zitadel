import { Injectable } from '@angular/core';
import { BehaviorSubject, combineLatest, from, map, Observable, of, switchMap } from 'rxjs';

import { ManagementService } from './mgmt.service';
import { AddCustomLoginPolicyRequest, AddCustomLoginPolicyResponse } from '../proto/generated/zitadel/management_pb';
import { Duration } from 'google-protobuf/google/protobuf/duration_pb';
import { IDPOwnerType, Provider } from '../proto/generated/zitadel/idp_pb';
import { AdminService } from './admin.service';
import { catchError } from 'rxjs/operators';
import { LoginPolicy } from '../proto/generated/zitadel/policy_pb';

@Injectable({
  providedIn: 'root',
})
export class ActivateIdpService {
  constructor() {}

  public activateIdp(
    service: AdminService | ManagementService,
    id: string,
    owner?: IDPOwnerType,
    policy?: LoginPolicy.AsObject,
  ): Observable<any> {
    const isAdmin = service instanceof AdminService;
    if (isAdmin) {
      service.addIDPToLoginPolicy(id);
    }

    if (!isAdmin && owner !== IDPOwnerType.IDP_OWNER_TYPE_SYSTEM && owner !== IDPOwnerType.IDP_OWNER_TYPE_ORG) {
      throw new Error('Must specify owner for management service');
    }

    return from(service.addIDPToLoginPolicy(id!, owner!)).pipe(
      catchError((error) => {
        if (isAdmin || error.code != 5) {
          throw error;
        }
        // No org policy was found, so we create a new one
        return from(policy ? of(policy) : from(service.getLoginPolicy()).pipe(map((policy) => policy.policy))).pipe(
          switchMap((policy) => {
            if (!policy?.isDefault) {
              // There is already an org policy
              throw error;
            }
            return from(this.addLoginPolicy(service, policy)).pipe(
              switchMap(() => from(service.addIDPToLoginPolicy(id!, owner!))),
            );
          }),
        );
      }),
    );
  }

  public addLoginPolicy(
    service: ManagementService,
    policy: LoginPolicy.AsObject,
  ): Promise<AddCustomLoginPolicyResponse.AsObject> {
    const mgmtreq = new AddCustomLoginPolicyRequest();
    mgmtreq.setAllowExternalIdp(policy.allowExternalIdp);
    mgmtreq.setAllowRegister(policy.allowRegister);
    mgmtreq.setAllowUsernamePassword(policy.allowUsernamePassword);
    mgmtreq.setForceMfa(policy.forceMfa);
    mgmtreq.setPasswordlessType(policy.passwordlessType);
    mgmtreq.setHidePasswordReset(policy.hidePasswordReset);
    mgmtreq.setMultiFactorsList(policy.multiFactorsList);
    mgmtreq.setSecondFactorsList(policy.secondFactorsList);

    const pcl = new Duration()
      .setSeconds(policy.passwordCheckLifetime?.seconds ?? 0)
      .setNanos(policy.passwordCheckLifetime?.nanos ?? 0);
    mgmtreq.setPasswordCheckLifetime(pcl);

    const elcl = new Duration()
      .setSeconds(policy.externalLoginCheckLifetime?.seconds ?? 0)
      .setNanos(policy.externalLoginCheckLifetime?.nanos ?? 0);
    mgmtreq.setExternalLoginCheckLifetime(elcl);

    const misl = new Duration()
      .setSeconds(policy.mfaInitSkipLifetime?.seconds ?? 0)
      .setNanos(policy.mfaInitSkipLifetime?.nanos ?? 0);
    mgmtreq.setMfaInitSkipLifetime(misl);

    const sfcl = new Duration()
      .setSeconds(policy.secondFactorCheckLifetime?.seconds ?? 0)
      .setNanos(policy.secondFactorCheckLifetime?.nanos ?? 0);
    mgmtreq.setSecondFactorCheckLifetime(sfcl);

    const mficl = new Duration()
      .setSeconds(policy.multiFactorCheckLifetime?.seconds ?? 0)
      .setNanos(policy.multiFactorCheckLifetime?.nanos ?? 0);
    mgmtreq.setMultiFactorCheckLifetime(mficl);

    mgmtreq.setAllowDomainDiscovery(policy.allowDomainDiscovery);
    mgmtreq.setIgnoreUnknownUsernames(policy.ignoreUnknownUsernames);
    mgmtreq.setDefaultRedirectUri(policy.defaultRedirectUri);

    mgmtreq.setDisableLoginWithEmail(policy.disableLoginWithEmail);
    mgmtreq.setDisableLoginWithPhone(policy.disableLoginWithPhone);
    mgmtreq.setForceMfaLocalOnly(policy.forceMfaLocalOnly);

    return service.addCustomLoginPolicy(mgmtreq);
  }
}
