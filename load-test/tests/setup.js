import url from './url.js';
import http from 'k6/http';
import { check } from 'k6';
import { Trend } from 'k6/metrics';

import { Config } from './config.js';

export default async function Setup(accessToken) {
    const org = await createOrg(accessToken);
    const user = await createHuman('gigi', org, accessToken);
    return {org, user};
}

const createOrgTrend = new Trend('setup_create_org_duration', true);
export function createOrg(accessToken) {
    return new Promise((resolve, reject) => {
        let response = http.asyncRequest(
            'POST', 
            url('/v2beta/organizations'), 
            JSON.stringify({
                name: `load-test-${new Date(Date.now()).toISOString()}`
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

const createHumanTrend = new Trend('setup_create_human_duration', true);
export function createHuman(username, org, accessToken){
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
        
            createHumanTrend.add(res.timings.duration);
        
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

const createMachineTrend = new Trend('setup_create_machine_duration', true);
export function createMachine(username, org, accessToken){
    return new Promise((resolve, reject) => {
        let response = http.asyncRequest(
            'POST',
            url('/management/v1/users/machine'), 
            JSON.stringify({
                userName: username,
                name: username,
                // bearer
                access_token_type: 0
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
                "create user is status ok": (r) => r.status === 200
            }) || reject(`unable to create user(username: ${username}) status: ${res.status} body: ${res.body}`);
        
            createHumanTrend.add(res.timings.duration);
        
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
        });
    })
}

const addMachinePatTrend = new Trend('setup_add_machine_pat_duration', true);
export function addMachinePat(userId, org, accessToken){
    return new Promise((resolve, reject) => {
        let response = http.asyncRequest(
            'POST',
            url(`/management/v1/users/${userId}/pats`), 
            null, 
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
                "add pat status ok": (r) => r.status === 200
            }) || reject(`unable to add pat (user id: ${userId}) status: ${res.status} body: ${res.body}`);
            resolve(res.json());
        });
    });
}