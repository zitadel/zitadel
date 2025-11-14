import express from 'express'
import * as notifications from './notifications.js'
import { readFileSync } from 'fs'

console.log("Starting mock server from", process.cwd());
if (!process.env.ZITADEL_ADMIN_TOKEN_FILE) {
    throw new Error("ZITADEL_ADMIN_TOKEN_FILE environment variable is not set");
}
if (!process.env.ZITADEL_API_URL) {
    throw new Error("ZITADEL_API_URL environment variable is not set");
}

console.log("Using the following configuration:");
console.log(`ZITADEL_ADMIN_TOKEN_FILE: ${process.env.ZITADEL_ADMIN_TOKEN_FILE}`);
console.log(`ZITADEL_API_URL: ${process.env.ZITADEL_API_URL}`);

const selfPort = 3333
const apiToken = readFileSync(process.env.ZITADEL_ADMIN_TOKEN_FILE!).toString().trim();
const apiUrl = process.env.ZITADEL_API_URL!;

const router = express()
router.use(express.json())

await Promise.all([
    notifications.setup(selfPort, apiUrl, apiToken)
])

notifications.serve(router)

router.get('/ready', (_, res) => {
    res.send('OK\n')
})

router.listen(selfPort, () => {
    console.log(`Server ready at: http://localhost:${selfPort}`)
})

