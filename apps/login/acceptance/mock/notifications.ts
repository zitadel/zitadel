import { Router } from "express";
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
    let notifications: { [key: string]: Object } = {}
    router.post(EMAIL_PATH, (req, res) => {
        // Receive email webhook
        const { contextInfo: { recipientEmailAddress } } = req.body
        console.log("saving email for", recipientEmailAddress);
        notifications[recipientEmailAddress] = req.body
        res.send('Email!\n')
    })
    router.post(SMS_PATH, (req, res) => {
        // Receive SMS webhook
        const { contextInfo: { recipientPhoneNumber } } = req.body
        console.log("saving SMS for", recipientPhoneNumber);
        notifications[recipientPhoneNumber] = req.body
        res.send('SMS!\n')
    })
    router.get('/notifications/:recipient', (req, res) => {
        // Return notification for recipient to test case and remove it from memory
        const { recipient } = req.params
        const notification = notifications[recipient]
        if (!notification) {
            console.log("no notification found for", recipient, "in", Object.keys(notifications));
            return res.status(404).send('No message found\n')
        }
        console.log("returning and removing notification for", recipient);
        delete notifications[recipient]
        res.contentType('application/json').json(notification)
    })
    router.get('/notifications', (req, res) => {
        // for debugging purposes
        res.contentType('application/json').json(notifications)
    })
}
