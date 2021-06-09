import { HttpEvent, HttpHandler, HttpInterceptor, HttpRequest } from '@angular/common/http';
import { OAuthModuleConfig } from 'angular-oauth2-oidc';
import { Observable } from 'rxjs';

import { Org } from '../../proto/generated/zitadel/org_pb';
import { StorageService } from '../storage.service';

const orgKey = 'x-zitadel-orgid';
const ORG_STORAGE_KEY = 'organization';
export abstract class HttpOrgInterceptor implements HttpInterceptor {
  private org!: Org.AsObject;

  protected get validUrls(): string[] {
    return this.oauthModuleConfig.resourceServer.allowedUrls || [];
  }

  constructor(
    private storageService: StorageService,
    protected oauthModuleConfig: OAuthModuleConfig,
  ) {
    const org: Org.AsObject | null = (this.storageService.getItem(ORG_STORAGE_KEY));

    if (org) {
      this.org = org;
    }
  }

  public intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    if (!this.urlValidation(req.url)) {
      return next.handle(req);
    }

    return next.handle(req.clone({
      setHeaders: {
        [orgKey]: this.org.id
      },
    }));
  }

  private urlValidation(toIntercept: string): boolean {
    const URLS = ['https://api.zitadel.dev/assets', 'https://api.zitadel.ch/assets'];

    return URLS.findIndex(url => toIntercept.startsWith(url)) > -1;
  }
}
