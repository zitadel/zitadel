import { Interceptor } from "@connectrpc/connect";
import { context, propagation } from "@opentelemetry/api";

export const tracingInterceptor: Interceptor = (next) => async (req) => {
    const fields: Record<string, string> = {};
    const activeContext = context.active();
    propagation.inject(activeContext, fields);

    console.log("Tracing Interceptor - Active Context:", activeContext);
    console.log("Tracing Interceptor - Injected Fields:", fields);

    for (const [key, value] of Object.entries(fields)) {
        req.header.set(key, value);
    }

    return next(req);
};
