import { HTTPRequest, launch, Page } from 'puppeteer'

module.exports = (on, config) => {

  config.env.newEmail = "demo@caos.ch"
  config.env.newUserName = "demo"
  config.env.newFirstName = "demofirstname"
  config.env.newLastName = "demolastname"
  config.env.newPhonenumber = "+41 123456789"

  config.env.newMachineUserName = "machineusername"
  config.env.newMachineName = "name"
  config.env.newMachineDesription = "description"

  on('task', {
    login({username, password }) {
      return login(config.env.consoleUrl, `${config.env.consoleUrl}/auth/callback`, username, password)
    },
  })  

  // IMPORTANT return the updated config object
  return config

}


async function login(loginUrl: string, callbackUrl: string, username: string, password: string) {

  const browser = await launch({
    headless: true,
    args: ['--no-sandbox', '--disable-setuid-sandbox']
  });

  const page = await browser.newPage();
  try {

    await page.setRequestInterception(true);
    page.on('request', interceptRequest(callbackUrl));

    await page.setViewport({ width: 1280, height: 800 });
    await page.goto(loginUrl);

    // Enter credentials.
    await submitInput(page, '#loginName', username);
    const response = await submitInput(page, '#password', password);

    // The login failed.
    if (response.status() >= 400) {
      throw new Error(`'Login with user ${username} failed, error ${response.status()}`);
    }

    // Redirected to MFA/consent/... which is not implemented yet.
    const url = response.url();
    if (url.indexOf(callbackUrl) !== 0) {
      throw new Error(`User was redirected to unexpected location: ${url}`);
    }

    // Now let's fetch all cookies.
    const { cookies } = await page._client.send('Network.getAllCookies', {});
    return {
      callbackUrl: url,
      cookies
    };
  } finally {
    await page.close();
    await browser.close();
  }
  return null
}

async function submitInput(page: Page, elementSelector: string, value: string) {

  const input = await page.waitForSelector(elementSelector, { timeout: 10000, visible: true })
  await input.type(value);
  const [response] = await Promise.all([page.waitForNavigation({ waitUntil: 'networkidle0' }), page.click('#submit-button')]);
  return response;
}


function interceptRequest(callbackUrl: string) {
  return (request: HTTPRequest) => {
    const url = request.url();
    if (request.isNavigationRequest() && url.indexOf(callbackUrl) === 0) {
      request.respond({ body: url, status: 200 });
    } else {
      request.continue();
    }
  };
}
