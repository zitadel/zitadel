import { Injectable } from '@angular/core';
import { from, map, Observable, of, switchMap } from 'rxjs';

import { ManagementService } from './mgmt.service';
import { AddCustomLoginPolicyRequest, AddCustomLoginPolicyResponse } from '../proto/generated/zitadel/management_pb';
import { Duration } from 'google-protobuf/google/protobuf/duration_pb';
import { IDPOwnerType } from '../proto/generated/zitadel/idp_pb';
import { AdminService } from './admin.service';
import { catchError } from 'rxjs/operators';
import { LoginPolicy } from '../proto/generated/zitadel/policy_pb';

@Injectable({
  providedIn: 'root',
})
export class LoginPolicyService {
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
            return from(this.createCustomLoginPolicy(service, policy, id));
          }),
        );
      }),
    );
  }

  public createCustomLoginPolicy(
    service: ManagementService,
    fromDefaultPolicy: LoginPolicy.AsObject,
    activateOrgIdp?: string,
  ): Promise<AddCustomLoginPolicyResponse.AsObject> {
    const mgmtreq = new AddCustomLoginPolicyRequest();
    mgmtreq.setAllowExternalIdp(fromDefaultPolicy.allowExternalIdp);
    mgmtreq.setAllowRegister(fromDefaultPolicy.allowRegister);
    mgmtreq.setAllowUsernamePassword(fromDefaultPolicy.allowUsernamePassword);
    mgmtreq.setForceMfa(fromDefaultPolicy.forceMfa);
    mgmtreq.setPasswordlessType(fromDefaultPolicy.passwordlessType);
    mgmtreq.setHidePasswordReset(fromDefaultPolicy.hidePasswordReset);
    mgmtreq.setMultiFactorsList(fromDefaultPolicy.multiFactorsList);
    mgmtreq.setSecondFactorsList(fromDefaultPolicy.secondFactorsList);

    const pcl = new Duration()
      .setSeconds(fromDefaultPolicy.passwordCheckLifetime?.seconds ?? 0)
      .setNanos(fromDefaultPolicy.passwordCheckLifetime?.nanos ?? 0);
    mgmtreq.setPasswordCheckLifetime(pcl);

    const elcl = new Duration()
      .setSeconds(fromDefaultPolicy.externalLoginCheckLifetime?.seconds ?? 0)
      .setNanos(fromDefaultPolicy.externalLoginCheckLifetime?.nanos ?? 0);
    mgmtreq.setExternalLoginCheckLifetime(elcl);

    const misl = new Duration()
      .setSeconds(fromDefaultPolicy.mfaInitSkipLifetime?.seconds ?? 0)
      .setNanos(fromDefaultPolicy.mfaInitSkipLifetime?.nanos ?? 0);
    mgmtreq.setMfaInitSkipLifetime(misl);

    const sfcl = new Duration()
      .setSeconds(fromDefaultPolicy.secondFactorCheckLifetime?.seconds ?? 0)
      .setNanos(fromDefaultPolicy.secondFactorCheckLifetime?.nanos ?? 0);
    mgmtreq.setSecondFactorCheckLifetime(sfcl);

    const mficl = new Duration()
      .setSeconds(fromDefaultPolicy.multiFactorCheckLifetime?.seconds ?? 0)
      .setNanos(fromDefaultPolicy.multiFactorCheckLifetime?.nanos ?? 0);
    mgmtreq.setMultiFactorCheckLifetime(mficl);

    mgmtreq.setAllowDomainDiscovery(fromDefaultPolicy.allowDomainDiscovery);
    mgmtreq.setIgnoreUnknownUsernames(fromDefaultPolicy.ignoreUnknownUsernames);
    mgmtreq.setDefaultRedirectUri(fromDefaultPolicy.defaultRedirectUri);

    mgmtreq.setDisableLoginWithEmail(fromDefaultPolicy.disableLoginWithEmail);
    mgmtreq.setDisableLoginWithPhone(fromDefaultPolicy.disableLoginWithPhone);
    mgmtreq.setForceMfaLocalOnly(fromDefaultPolicy.forceMfaLocalOnly);

    mgmtreq.setIdpsList(
      fromDefaultPolicy.idpsList.map((idp) => addIdpMessage(idp.idpId, IDPOwnerType.IDP_OWNER_TYPE_SYSTEM)),
    );
    if (activateOrgIdp) {
      mgmtreq.addIdps(addIdpMessage(activateOrgIdp, IDPOwnerType.IDP_OWNER_TYPE_ORG));
    }
    return service.addCustomLoginPolicy(mgmtreq);
  }
}

function addIdpMessage(id: string, owner: IDPOwnerType): AddCustomLoginPolicyRequest.IDP {
  const addIdp = new AddCustomLoginPolicyRequest.IDP();
  addIdp.setIdpId(id);
  addIdp.setOwnertype(owner);
  return addIdp;
}
