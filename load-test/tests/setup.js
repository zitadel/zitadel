import url from './url.js';
import http from 'k6/http';
import { check } from 'k6';
import { Trend } from 'k6/metrics';

import { Config } from './config.js';

export default async function Setup(accessToken) {
    const org = await createOrg(accessToken);
    const user = await createUser('gigi', org, accessToken);
    return {org, user};
}

const createOrgTrend = new Trend('setup_create_org_duration', true);
export function createOrg(accessToken) {
    return new Promise((resolve, reject) => {
        let response = http.asyncRequest(
            'POST', 
            url('/v2beta/organizations'), 
            JSON.stringify({
                name: 'load-test-' + new Date(Date.now()).toISOString() + `-${__VU}`
            }), 
            {
                headers: {
                    authorization: `Bearer ${accessToken}`,
                    'Content-Type': 'application/json',
                    'x-zitadel-orgid': Config.orgId
                }
            }
        )

        response.then((res) => {
            check(res, {
                'org created': (r) => r.status === 201  || reject(`unable to create org status: ${res.status} || body: ${res.body}`)
            });
        
            createOrgTrend.add(res.timings.duration);
        
            resolve(res.json());
        });
    })
}

const createUserTrend = new Trend('setup_create_user_duration', true);
export function createUser(username, org, accessToken){
    return new Promise((resolve, reject) => {
        let response = http.asyncRequest(
            'POST',
            url('/v2beta/users/human'), 
            JSON.stringify({
                username: username,
                organization: {
                    orgId: org.organizationId
                },
                profile: {
                    givenName: 'Gigi',
                    familyName: 'Zitizen',
                },
                email: {
                    email: `zitizen-@caos.ch`,
                    isVerified: true,
                },
                password: {
                    password: 'Password1!',
                    changeRequired: false
                }
            }), 
            {
                headers: {
                    authorization: `Bearer ${accessToken}`,
                    'Content-Type': 'application/json',
                    'x-zitadel-orgid': org.organizationId
                }
            }
        );
    
        response.then((res) => {
            check(res, {
                "create user is status ok": (r) => r.status === 201
            }) || reject(`unable to create user(username: ${username}) status: ${res.status} body: ${res.body}`);
        
            createUserTrend.add(res.timings.duration);
        
            resolve(http.get(
                url(`/v2beta/users/${res.json().userId}`), 
                {
                    headers: {
                        authorization: `Bearer ${accessToken}`,
                        'Content-Type': 'application/json',
                        'x-zitadel-orgid': org.organizationId
                    }
                }
            ).json().user);
        })
    })
}