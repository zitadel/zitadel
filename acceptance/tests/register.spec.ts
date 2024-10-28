import {test} from "@playwright/test";
import {registerWithPassword} from './register';
import {loginWithPassword} from "./login";

test("register with password", async ({page}) => {
    const firstname = "firstname"
    const lastname = "lastname"
    const username = "register@example.com"
    const password = "Password1!"
    await registerWithPassword(page, firstname, lastname, username, password, password)
    await page.getByRole("heading", {name: "Welcome " + lastname + " " + lastname + "!"}).click();
    await loginWithPassword(page, username, password)
});
