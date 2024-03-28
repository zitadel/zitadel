import http from "k6/http";
import url from "./url.js";
import { check } from "k6";

export function removeOrg(org, accessToken) {
    const response = http.del(
        url('/management/v1/orgs/me'),
        null,
        {
            headers: {
                authorization: `Bearer ${accessToken}`,
                'x-zitadel-orgid': org.organizationId
            }
        }
    );

    check(response, {
        'org removed': (r) => r.status === 200
    });

    return response.json()
}