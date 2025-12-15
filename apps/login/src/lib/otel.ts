import { Interceptor } from "@connectrpc/connect";
import { context, propagation } from "@opentelemetry/api";

export const tracingInterceptor: Interceptor = (next) => async (req) => {
    const fields: Record<string, string> = {};
    propagation.inject(context.active(), fields);

    for (const [key, value] of Object.entries(fields)) {
        req.header.set(key, value);
    }

    return next(req);
};
