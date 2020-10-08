import { Injectable } from '@angular/core';
import { Request, UnaryInterceptor, UnaryResponse } from 'grpc-web';


const i18nHeader = 'Accept-Language';
@Injectable({ providedIn: 'root' })
/**
 * Set the navigator language as header to all grpc requests
 */
export class I18nInterceptor<TReq = unknown, TResp = unknown> implements UnaryInterceptor<TReq, TResp> {
    constructor() { }

    public async intercept(request: Request<TReq, TResp>, invoker: any): Promise<UnaryResponse<TReq, TResp>> {
        const metadata = request.getMetadata();

        const navLang = navigator.language;
        metadata[i18nHeader] = navLang;

        return invoker(request).then((response: any) => {
            return response;
        }).catch((error: any) => {
            return Promise.reject(error);
        });
    }
}
