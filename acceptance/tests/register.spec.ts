import {test} from "@playwright/test";
import {registerWithPassword} from './register';
import {loginWithPassword} from "./login";

test("register with password", async ({page}) => {
    const username = "register@example.com"
    const password = "Password1!"
    await registerWithPassword(page, "firstname", "lastname", username, password, password)
    await loginWithPassword(page, username, password)
});
