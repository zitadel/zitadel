import { apiAuth } from "../../support/api/apiauth";
import { ensureHumanUserExists, ensureUserDoesntExist } from "../../support/api/users";
import { login, User, username } from "../../support/login/users";

describe('humans', () => {

    const humansPath = `${Cypress.env('consoleUrl')}/users/list/humans`
    const testHumanUserName = 'e2ehumanusername'

    ;[User.OrgOwner].forEach(user => {
 
        describe(`as user "${user}"`, () => {

            beforeEach(()=> {
                login(user)
                cy.visit(humansPath)
                cy.get('[data-cy=timestamp]')
            })

            describe('add', () => {
                before(`ensure it doesn't exist already`, () => {
                    apiAuth().then(apiCallProperties => {
                        ensureUserDoesntExist(apiCallProperties, testHumanUserName)
                    })
                })

                it('should add a user', () => {
                    cy.contains('a', 'New').click()
                    cy.url().should('contain', 'users/create')
                    cy.get('[formcontrolname^=email]').type(username('e2ehuman'))
                    //force needed due to the prefilled username prefix
                    cy.get('[formcontrolname^=userName]').type(testHumanUserName, {force: true})
                    cy.get('[formcontrolname^=firstName]').type('e2ehumanfirstname')
                    cy.get('[formcontrolname^=lastName]').type('e2ehumanlastname')
                    cy.get('[formcontrolname^=phone]').type('+41 123456789')
                    cy.get('button').filter(':contains("Create")').should('be.visible').click()
                    cy.contains('User created successfully')
                    cy.visit(humansPath);
                    cy.contains("tr", testHumanUserName)
                })        
            })
            
            describe('remove', () => {
                before('ensure it exists', () => {
                    apiAuth().then(api => {
                        ensureHumanUserExists(api, testHumanUserName)
                    })                    
                })

                it('should delete a human user', () => {
                    cy.get('h1')
                        .contains('Users')
                        .parent()
                        .contains("tr", testHumanUserName, { timeout: 1000 })
                        .find('button')
                        //force due to angular hidden buttons
                        .click({force: true})
                    cy.get('span.title')
                        .contains('Delete User')
                        .parent()
                        .find('button')
                        .contains('Delete')
                        .click()
                    cy.contains('User deleted successfully')
                    cy.get(`[text*=${testHumanUserName}]`).should('not.exist');                    
                })
            })
        })
    })
})
/*
describe("users", ()=> {

    before(()=> {
        cy.consolelogin(Cypress.env('username'), Cypress.env('password'), Cypress.env('consoleUrl'))
    })

    it('should show personal information', () => {
        cy.log(`USER: show personal information`);
        //click on user information 
        cy.get('a[href*="users/me"').eq(0).click()
        cy.url().should('contain', '/users/me')
    })

    it('should show users', () => {
        cy.visit(Cypress.env('consoleUrl') + '/users/list/humans')
        cy.url().should('contain', 'users/list/humans')
    })
})

*/