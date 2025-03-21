import { Injectable } from '@angular/core';
import { Request, RpcError, StatusCode, UnaryInterceptor, UnaryResponse } from 'grpc-web';
import { Org } from 'src/app/proto/generated/zitadel/org_pb';

import { StorageKey, StorageLocation, StorageService } from '../storage.service';
import { ConnectError, Interceptor } from '@connectrpc/connect';
import { firstValueFrom, identity, Observable, Subject } from 'rxjs';
import { debounceTime, filter, map } from 'rxjs/operators';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';

const ORG_HEADER_KEY = 'x-zitadel-orgid';
@Injectable({ providedIn: 'root' })
export class OrgInterceptor<TReq = unknown, TResp = unknown> implements UnaryInterceptor<TReq, TResp> {
  constructor(private readonly orgInterceptorProvider: OrgInterceptorProvider) {}

  public async intercept(request: Request<TReq, TResp>, invoker: any): Promise<UnaryResponse<TReq, TResp>> {
    const metadata = request.getMetadata();

    const orgId = this.orgInterceptorProvider.getOrgId();
    if (orgId) {
      metadata[ORG_HEADER_KEY] = orgId;
    }

    return invoker(request).catch(this.orgInterceptorProvider.handleError);
  }
}

export function NewConnectWebOrgInterceptor(orgInterceptorProvider: OrgInterceptorProvider): Interceptor {
  return (next) => async (req) => {
    if (!req.header.get(ORG_HEADER_KEY)) {
      const orgId = orgInterceptorProvider.getOrgId();
      if (orgId) {
        req.header.set(ORG_HEADER_KEY, orgId);
      }
    }

    return next(req).catch(orgInterceptorProvider.handleError);
  };
}

@Injectable({ providedIn: 'root' })
export class OrgInterceptorProvider {
  constructor(private storageService: StorageService) {}

  getOrgId() {
    const org: Org.AsObject | null = this.storageService.getItem(StorageKey.organization, StorageLocation.session);
    return org?.id;
  }

  handleError = (error: any): never => {
    if (!(error instanceof RpcError) && !(error instanceof ConnectError)) {
      throw error;
    }

    if (
      error instanceof RpcError &&
      error.code === StatusCode.PERMISSION_DENIED &&
      error.message.startsWith("Organisation doesn't exist")
    ) {
      this.storageService.removeItem(StorageKey.organization, StorageLocation.session);
    }

    throw error;
  };
}
