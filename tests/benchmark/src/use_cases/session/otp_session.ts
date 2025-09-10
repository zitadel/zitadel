import { loginByUsernamePassword } from '../../login_ui';
import { createOrg, removeOrg } from '../../org';
import { createHuman, setEmailOTPOnHuman, User } from '../../user';
import { Trend } from 'k6/metrics';
import { Config, MaxVUs } from '../../config';
import { createSession, setSession } from '../../session';

export async function setup() {
  const tokens = loginByUsernamePassword(Config.admin as User);
  console.log('setup: admin signed in');

  const org = await createOrg(tokens.accessToken!);
  console.log(`setup: org (${org.organizationId}) created`);

  const humanPromises = Array.from({ length: MaxVUs() }, (_, i) => {
    return createHuman(`zitizen-${i}`, org, tokens.accessToken!);
  });

  const humans = (await Promise.all(humanPromises)).map((user) => {
    setEmailOTPOnHuman(user, org, tokens.accessToken!);
    return { userId: user.userId, loginName: user.loginNames[0], password: 'Password1!' };
  });
  console.log(`setup: ${humans.length} users created`);

  return { tokens, org, users: humans };
}

// implements the flow described in
// https://zitadel.com/docs/guides/integrate/login-ui/oidc-standard
const otpSessionTrend = new Trend('otp_session_duration', true);
export default async function (data: any) {
  const start = new Date();
  let session = await createSession(data.org, data.tokens.accessToken, {
    user: {
      loginName: data.users[__VU - 1].loginName,
    },
  });
  const sessionId = (session as any).sessionId;

  session = await setSession(sessionId, session, data.tokens.accessToken, {
    otpEmail: {
      return_code: {},
    },
  });

  session = await setSession(sessionId, session, data.tokens.accessToken, null, {
    otpEmail: {
      code: session.challenges.otpEmail,
    },
  });

  otpSessionTrend.add(new Date().getTime() - start.getTime());
}

export function teardown(data: any) {
  removeOrg(data.org, data.tokens.accessToken);
}
