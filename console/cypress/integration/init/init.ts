//import { apiAuth } from "../../support/api/apiauth";
//import { ensureMachineUserExists, ensureUserDoesntExist } from "../../support/api/users";
import { login, User } from "../../support/login/users";

describe('initialize organisation', () => {
    it('initializes', () => {
        login(User.IAMAdminUser)
    })
})