import { Trend } from "k6/metrics";
import http from 'k6/http';
import url from './url.js';
import { check } from "k6";


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
        
            const user = http.get(
                url(`/v2beta/users/${res.json().userId}`), 
                {
                    headers: {
                        authorization: `Bearer ${accessToken}`,
                        'Content-Type': 'application/json',
                        'x-zitadel-orgid': org.organizationId
                    }
                }
            );
            resolve(user.json().user);
        }).catch((reason) => {
            reject(reason);
        });
    })
}

const updateHumanTrend = new Trend('update_human_duration', true);
export function updateHuman(payload = {}, userId, org, accessToken) {
    return new Promise((resolve, reject) => {
        let response = http.asyncRequest(
            'PUT',
            url(`/v2beta/users/${userId}`),
            JSON.stringify(payload),
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
                "update user is status ok": (r) => r.status === 201
            });
            updateHumanTrend.add(res.timings.duration);
            resolve(res);
        }).catch((reason) => {
            reject(reason);
        });
    });
}

const lockUserTrend = new Trend('lock_user_duration', true);
export function lockUser(userId, org, accessToken) {
    return new Promise((resolve, reject) => {
        let response = http.asyncRequest(
            'POST',
            url(`/v2beta/users/${userId}/lock`),
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
                "update user is status ok": (r) => r.status === 201
            });
            lockUserTrend.add(res.timings.duration);
            resolve(res);
        }).catch((reason) => {
            reject(reason);
        });
    });
}

const deleteUserTrend = new Trend('delete_user_duration', true);
export function deleteUser(userId, org, accessToken) {
    return new Promise((resolve, reject) => {
        let response = http.asyncRequest(
            'DELETE',
            url(`/v2beta/users/${userId}`),
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
                "update user is status ok": (r) => r.status === 201
            });
            deleteUserTrend.add(res.timings.duration);
            resolve(res);
        }).catch((reason) => {
            reject(reason);
        });
    });
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
        
            createMachineTrend.add(res.timings.duration);
        
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
            
            addMachinePatTrend.add(res.timings.duration);
            resolve(res.json());
        });
    });
}