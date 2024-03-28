import { loginByUsernamePassword } from '../login_ui.js';
import Setup from '../setup.js';
import { userinfo } from '../oidc.js';
import { Config } from '../config.js';

export async function setup() {
  const tokens = loginByUsernamePassword(Config.admin);
  const setup = await Setup(tokens.accessToken);
  return {tokens: loginByUsernamePassword({loginName: setup.user.loginNames[0], password: 'Password1!'})};
}

export default function (data) {
  userinfo(data.tokens.accessToken);
}

