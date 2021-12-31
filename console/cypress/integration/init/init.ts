//import { apiAuth } from "../../support/api/apiauth";
//import { ensureMachineUserExists, ensureUserDoesntExist } from "../../support/api/users";
import { login, User } from "../../support/login/users";


describe('initialize organisation', () => {
    it('initializes', () => {
        const adminPw = 'Password1!'

        login(User.IAMAdminUser, false, adminPw, null, null, () => {
// TODO: Not always
//            cy.contains('button', 'skip').click()
/*            cy.get('#change-old-password').type(adminPw)
            cy.get('#change-new-password').type(adminPw)
            cy.get('#change-password-confirmation').type(adminPw)
            cy.pause()
            cy.contains('button', 'next').click()
            cy.contains('button', 'next').click()*/
        })
    })
})