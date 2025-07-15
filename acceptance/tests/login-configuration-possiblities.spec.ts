import test from "@playwright/test";

test("login with mfa setup, mfa setup prompt", async ({ page }) => {
  test.skip();
  // Given the organization has enabled at least one mfa types
  // Given the user has a password but no mfa registered
  // User authenticates with login name and password
  // User is prompted to setup a mfa, mfa providers are listed, the user can choose the provider
});

test("login with mfa setup, no mfa setup prompt", async ({ page }) => {
  test.skip();
  // Given the organization has set "multifactor init check time" to 0
  // Given the organization has enabled mfa types
  // Given the user has a password but no mfa registered
  // User authenticates with loginname and password
  // user is directly loged in and not prompted to setup mfa
});

test("login with mfa setup, force mfa for local authenticated users", async ({ page }) => {
  test.skip();
  // Given the organization has enabled force mfa for local authentiacted users
  // Given the organization has enabled all possible mfa types
  // Given the user has a password but no mfa registered
  // User authenticates with loginname and password
  // User is prompted to setup a mfa, all possible mfa providers are listed, the user can choose the provider
});

test("login with mfa setup, force mfa - local user", async ({ page }) => {
  test.skip();
  // Given the organization has enabled force mfa for local authentiacted users
  // Given the organization has enabled all possible mfa types
  // Given the user has a password but no mfa registered
  // User authenticates with loginname and password
  // User is prompted to setup a mfa, all possible mfa providers are listed, the user can choose the provider
});

test("login with mfa setup, force mfa - external user", async ({ page }) => {
  test.skip();
  // Given the organization has enabled force mfa
  // Given the organization has enabled all possible mfa types
  // Given the user has an idp but no mfa registered
  // enter login name
  // redirect to configured external idp
  // User is prompted to setup a mfa, all possible mfa providers are listed, the user can choose the provider
});

test("login with mfa setup, force mfa - local user, wrong password", async ({ page }) => {
  test.skip();
  // Given the organization has a password lockout policy set to 1 on the max password attempts
  // Given the user has only a password as auth methos
  // enter login name
  // enter wrong password
  // User will get an error "Wrong password"
  // enter password
  // User will get an error "Max password attempts reached - user is locked. Please reach out to your administrator"
});
