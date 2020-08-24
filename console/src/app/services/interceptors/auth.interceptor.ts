import { Injectable } from '@angular/core';

@Injectable({ providedIn: 'root' })
export class AuthInterceptor {
    constructor() { }

    public intercept(request: any, invoker: any): any {
        // Update the request message before the RPC.
        console.log(request);
        const reqMsg = request.getRequestMessage();
        reqMsg.setMessage('[Intercept request]' + reqMsg.getMessage());

        // After the RPC returns successfully, update the response.
        return invoker(request).then((response: any) => {
            // You can also do something with response metadata here.
            console.log(response.getMetadata());

            // Update the response message.
            const responseMsg = response.getResponseMessage();
            responseMsg.setMessage('[Intercept response]' + responseMsg.getMessage());

            return response;
        });
    };
}
