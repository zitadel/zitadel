import { loginByUsernamePassword } from '../login_ui.js';
import Setup from '../setup.js';
import { userinfo } from '../oidc.js';

// not using SharedArray here will mean that the code in the function call (that is what loads and
// parses the json) will be executed per each VU which also means that there will be a complete copy
// per each VU
const admin = JSON.parse(open('../data/admin.json'));

export async function setup() {
  const tokens = loginByUsernamePassword(admin);
  const setup = await Setup(tokens.accessToken);
  return {tokens: loginByUsernamePassword({loginName: setup.user.loginNames[0], password: 'Password1!'})};
}

export default function (data) {
  userinfo(data.tokens.accessToken);
}

