import { Injectable } from '@angular/core';
import { TranslateService } from '@ngx-translate/core';
import { Request, UnaryInterceptor, UnaryResponse } from 'grpc-web';

const i18nHeader = 'Accept-Language';
@Injectable({ providedIn: 'root' })
/**
 * Set the navigator language as header to all grpc requests
 */
export class I18nInterceptor<TReq = unknown, TResp = unknown> implements UnaryInterceptor<TReq, TResp> {
  constructor(private translate: TranslateService) {}

  public intercept(request: Request<TReq, TResp>, invoker: any): Promise<UnaryResponse<TReq, TResp>> {
    const metadata = request.getMetadata();

    const navLang = this.translate.currentLang ?? navigator.language;
    if (navLang) {
      metadata[i18nHeader] = navLang;
    }

    return invoker(request);
  }
}
