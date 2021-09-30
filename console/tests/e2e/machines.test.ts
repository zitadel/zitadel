import { test, expect, Page } from '@playwright/test';
import { login } from './commands/login'
import { ORG_OWNER } from "./models/users"
import fetch from 'node-fetch'
import { APICallProperties, prepareAPICalls } from './commands/apiauth';
import { checkStatus } from './fetch/status';

const TEST_MACHINE_USERNAME="testmachineusername"
const MACHINES_URL=`/users/list/machines`

test.describe('machine', () => {

    ;[ORG_OWNER].forEach(user => {

        test.describe(`as user ${user.username}`, () => {

            var page: Page
            var api: APICallProperties

            test.beforeAll(async ({browser}) => {
                const context = await browser.newContext({
                    recordVideo: {
                        dir: './tests/e2e/results/machine'
                    },
                    recordHar: {
                        path: './tests/e2e/results/machine'
                    }
                })
                page = await context.newPage()
                const res = await Promise.all([login(page, user), prepareAPICalls()])
                api = res[1]
            })

            test.beforeEach(async () => {
                // Navigate to machines list
                await page.click('a:has-text("Service Users")');
                await expect(page).toHaveURL(MACHINES_URL);                
            })

            test.describe('add', () => {

                test.beforeAll(async () => {

                    const usersResp = await fetch(`${api.baseURL}/management/v1/users/_search`, {
                        method: 'POST',
                        headers: { Authorization: api.authHeader },
                    })

                    checkStatus(usersResp)

                    const users = await usersResp.json() as {result: [{userName: string, id: string}]}

                    var machineUser = users.result.find(user => user.userName === TEST_MACHINE_USERNAME)
                    if (machineUser) {
                        const delResp = await fetch(`${api.baseURL}/management/v1/users/${machineUser.id}`, {
                            method: 'DELETE',
                            headers: { Authorization: api.authHeader },
                        })
                        checkStatus(delResp)
                    }
                })

                test('should add a machine', async () => {
                    
                    // Click new button
                    await Promise.all([
                        page.waitForNavigation(),
                        page.click('a:has-text("New")')
                    ]);
                    
                    // Fill username
                    await page.fill('input', TEST_MACHINE_USERNAME);
                    
                    // Fill name
                    await page.fill('text=Name* The input field is empty. >> input', 'name');
                    
                    // Fill description
                    await page.fill('#cnsl-input-3', 'description');
                    
                    // Submit
                    await Promise.all([
                        page.waitForNavigation(),
                        page.click('button:has-text("Create")')
                    ]);
                    
                    // Navigate to machines list
                    await page.click('a:has-text("Service Users")');
                    await expect(page).toHaveURL(MACHINES_URL);
                    
                    await page.waitForSelector(`table:has-text("${TEST_MACHINE_USERNAME}")`, { strict: true, state: 'attached' })                      

    /*                        cy.contains('a', 'New').click()
                    cy.url().should('contain', 'users/create-machine')
                    //force needed due to the prefilled username prefix
                    cy.get('[formcontrolname^=userName]').type(Cypress.env('newMachineUserName'),{force: true})
                    cy.get('[formcontrolname^=name]').type(Cypress.env('newMachineName'))
                    cy.get('[formcontrolname^=description]').type(Cypress.env('newMachineDesription'))
                    cy.get('button').filter(':contains("Create")').should('be.visible').click()
                    cy.contains('User created successfully')
                    cy.visit(Cypress.env('consoleUrl') + '/users/list/machines');
                    cy.contains("tr", Cypress.env('newMachineUserName'))*/
                })
            })

            test.describe('remove', () => {

                test.beforeAll(async () => {
                    const resp = await fetch(`${api.baseURL}/management/v1/users/machine`, {
                        method: 'POST',
                        headers: {
                            Authorization: api.authHeader,
                            'Content-Type': 'application/json',
                        },
                        body: JSON.stringify({
                            'user_name': TEST_MACHINE_USERNAME,
                            name: 'test',
                            description: 'e2e delete user test',
                        }),                            
                    })

                    // 409 = does not exist
                    if (resp.status != 409) {
                        checkStatus(resp)
                    }
                })

                test('should delete a machine', async () => {

                    // Hover over the service account so the delete button appears
                    await page.hover(`text=${TEST_MACHINE_USERNAME}`)

                    // Click the delete button
                    await page.click(`tr:has-text("${TEST_MACHINE_USERNAME}") >> button`);

                    // Confirm deletion
                    await page.click('button:has-text("Delete")');

                    // User message appears
                    // await expect(page.locator('div:has-text=User deleted successfully')).toBeVisible()

                    await page.waitForSelector(`table:has-text("${TEST_MACHINE_USERNAME}")`, { strict: true, state: 'detached' })
                })
            })
        })
    })
})
