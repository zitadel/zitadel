// Run from the repository root: bun console/scripts/self-service-email.test.cjs
// Console has no configured unit-test target; load the real methods with framework dependencies stubbed.
const assert = require('node:assert/strict');
const fs = require('node:fs');
const vm = require('node:vm');
const root = require('node:path').resolve(__dirname, '../..') + '/';
const ts = require(root + 'console/node_modules/typescript');
function load(relative) {
  const source = fs.readFileSync(root + relative, 'utf8');
  const js = ts.transpileModule(source, {
    compilerOptions: { module: ts.ModuleKind.CommonJS, target: ts.ScriptTarget.ES2022, experimentalDecorators: true },
  }).outputText;
  const exports = {};
  vm.runInNewContext(js, { exports, require: () => ({ Component: () => () => {}, Injectable: () => () => {} }) });
  return exports;
}
const { NewAuthService } = load('console/src/app/services/new-auth.service.ts');
const { AuthUserDetailComponent } = load(
  'console/src/app/pages/users/user-detail/auth-user-detail/auth-user-detail.component.ts',
);
(async () => {
  let resolve,
    reject,
    calls = 0,
    refreshed = 0,
    notices = 0,
    failures = 0;
  const transport = {
    authNew: {
      resendMyEmailVerification: (request) => {
        assert.equal(Object.keys(request).length, 0);
        calls++;
        return new Promise((ok, fail) => {
          resolve = ok;
          reject = fail;
        });
      },
    },
  };
  const component = Object.create(AuthUserDetailComponent.prototype);
  let pending = false;
  component.emailVerificationPending = () => pending;
  component.emailVerificationPending.set = (value) => {
    pending = value;
  };
  component.newAuthService = new NewAuthService(transport, {});
  component.newMgmtService = new Proxy(
    {},
    {
      get() {
        throw Error('Management API must not be used for self-service');
      },
    },
  );
  component.toast = {
    showInfo: () => {
      notices++;
    },
    showError: () => {
      failures++;
    },
  };
  component.invalidateUser = async () => {
    refreshed++;
  };
  const first = component.resendEmailVerification();
  assert.equal(calls, 1, 'resend must use the self-service Auth API');
  assert.equal(pending, true);
  await component.resendEmailVerification();
  assert.equal(calls, 1);
  resolve({});
  await first;
  assert.equal(pending, false);
  assert.equal(notices, 1);
  assert.equal(refreshed, 1);
  const second = component.resendEmailVerification();
  reject(Error('Network failure'));
  await second;
  assert.equal(pending, false);
  assert.equal(failures, 1);
  assert.equal(notices, 1);
  const retry = component.resendEmailVerification();
  resolve({});
  await retry;
  assert.equal(calls, 3);
  console.log('PASS: self-service transport, duplicate prevention, success refresh, error handling, retry');
})().catch((e) => {
  console.error(e);
  process.exit(1);
});
