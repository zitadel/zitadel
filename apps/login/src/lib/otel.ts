import { Interceptor } from "@connectrpc/connect";
import { context, propagation } from "@opentelemetry/api";

export const tracingInterceptor: Interceptor = (next) => async (req) => {
    const fields: Record<string, string> = { "traceparent": "", "tracestate": "" };
    const activeContext = context.active();
    console.log("Tracing Interceptor - Active Context before injection:", activeContext);
    propagation.inject(activeContext, fields);

    console.log("Tracing Interceptor - Active Context:", activeContext);
    console.log("Tracing Interceptor - Injected Fields:", fields);

    // This makes the request to e.g. 'http://localhost:8080/zitadel.org.v2.OrganizationService/ListOrganizations'
    for (const [key, value] of Object.entries(fields)) {
        req.header.set(key, value);
    }

    return next(req);
};
