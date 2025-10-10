import { Router, Request } from "express";
import { createServerTransport } from '@zitadel/client/node'
import { createClientFor } from '@zitadel/client'
import { AdminService } from "@zitadel/proto/zitadel/admin_pb";

const EMAIL_PATH = '/email'
const SMS_PATH = '/sms'

export async function setup(selfPort: number, apiUrl: string, apiToken: string) {
    const transport = createServerTransport(apiToken, {baseUrl: apiUrl})
    const client = createClientFor(AdminService)(transport)
    const selfBaseUrl = `http://localhost:${selfPort}`
    const { id: emailProviderId } = await client.addEmailProviderHTTP({ endpoint: selfBaseUrl + EMAIL_PATH, "description": "Email provider for login acceptance testing" })
    await client.activateEmailProvider({ id: emailProviderId })
    const { id: smsProviderId } = await client.addSMSProviderHTTP({ endpoint: selfBaseUrl + SMS_PATH, "description": "SMS provider for login acceptance testing" })
    await client.activateSMSProvider({ id: smsProviderId })
}

export function serve(router: Router) {
    channel(router, EMAIL_PATH, (req) => req.body.contextInfo.recipientEmailAddress)
    channel(router, SMS_PATH, (req) => req.body.contextInfo.recipientPhoneNumber)
}

function channel(router: Router, path: string, extractAddress: (req: Request) => string) {
    let notifications: { [key: string]: Object } = {}
    router.post(path, (req, res) => {
        const address = extractAddress(req)
        console.log("saving message for", address);
        notifications[address] = req.body
        res.send('OK\n')
    })
    router.get(`/notifications${path}/:address`, (req, res) => {
        // Return notification for recipient to test case and remove it from memory
        const { address } = req.params
        const notification = notifications[address]
        if (!notification) {
            console.log("no notification found for", address, "in", Object.keys(notifications));
            return res.status(404).send('No message found\n')
        }
        console.log("returning and removing notification for", address);
        delete notifications[address]
        res.contentType('application/json').json(notification)
    })
    router.get(`/notifications${path}`, (req, res) => {
        // for debugging purposes
        res.contentType('application/json').json(notifications)
    })
}
