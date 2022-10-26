

describe('login', () => {
    describe('existing user with username password', () => {
        it('login policy has no 2fa registered');
        it('login policy has not allowed passwordless');
        it('login policy mfa prompt lifetime is set to 0 which results in no prompt');
        it('A user with password and verified email does exist');

        it('User has to enter loginname');
        it('User has to enter password');
        it('No mfa should be reuqests');
        it('No passwordless is requested');
        it('User is redirected to application (requested redirect url)');
    });

    describe('existing user with 2fa otp', () => {
        it('login policy has otp registered');
        it('login policy has not allowed passwordless');
        it('login policy mfa prompt lifetime is set to 0 which results in no prompt');
        it('A user with password, verified email and otp does exist');

        it('User has to enter loginname');
        it('User has to enter password');
        it('User has to enter code for otp');
        it('No passwordless is requested');
        it('User is redirected to application (requested redirect url)');
    });

    describe('existing user with 2fa webauthn', () => {
        it('login policy has webauthn registered');
        it('login policy has not allowed passwordless');
        it('login policy mfa prompt lifetime is set to 0 which results in no prompt');
        it('A user with password, verified email and webauthn does exist');

        it('User has to enter loginname');
        it('User has to enter password');
        it('User has to prove webauthn');
        it('No passwordless is requested');
        it('User is redirected to application (requested redirect url)');
    });

    describe('existing user with passwordless', () => {
        it('login policy has webauthn as mfa registered');
        it('login policy has allowed passwordless');
        it('login policy mfa prompt lifetime is set to 0 which results in no prompt');
        it('A user with password, verified email and passwordless does exist');

        it('User has to enter loginname');
        it('User has to enter verify webauthn');
        it('User is redirected to application (requested redirect url)');
    });

    describe('existing user with external idp (google)', () => {
        it('login policy has enabled login with external idp');
        it('login policy has disabled login with username password');
        it('google idp is added and activated')
        it('login policy mfa prompt lifetime is set to 0 which results in no prompt');
        it('A user with linked google idp is existing');

        it('Auto redirect to google idp (because only possibility to login)');
        it('user authenticates in google');
        it('User with linked idp is found in zitadel');
        it('User is redirected to application (requested redirect url)');
    });

    describe('login with external idp (google) with auto register user', () => {
        it('login policy has enabled login with external idp');
        it('login policy has disabled login with username password');
        it('google idp is added and activated');
        it('auto register is enabled on google');
        it('login policy has register allowed');
        it('login policy mfa prompt lifetime is set to 0 which results in no prompt');
        it('No users with this linked idp does exist');

        it('Auto redirect to google idp (because only possibility to login)');
        it('user authenticates in google');
        it('Screen to choose either link user or register is shown');
        it('User chooses register');
        it('No verify email or init user screen is shown');
        it('No mfa prompt is shown');
        it('User is redirected to application (requested redirect url)');
    });

    describe('login with external idp (google) with manully register user', () => {
        it('login policy has enabled login with external idp');
        it('login policy has disabled login with username password');
        it('google idp is added and activated');
        it('auto register is disabled on google');
        it('login policy has register allowed');
        it('login policy mfa prompt lifetime is set to 0 which results in no prompt');
        it('No users with this linked idp does exist');

        it('Auto redirect to google idp (because only possibility to login)');
        it('user authenticates in google');
        it('Screen to choose either link user or register is shown');
        it('User chooses register');
        it('Screen with user fields is shown, user clicks register');
        it('No verify email or init user screen is shown');
        it('No mfa prompt is shown');
        it('User is redirected to application (requested redirect url)');
    });

    describe('login with external idp (google) with link user', () => {
        it('login policy has enabled login with external idp');
        it('login policy has disabled login with username password');
        it('google idp is added and activated');
        it('login policy has register allowed');
        it('login policy mfa prompt lifetime is set to 0 which results in no prompt');
        it('No users with this linked idp does exist');
        it('A user with username, password and active state exists');

        it('Auto redirect to google idp (because only possibility to login)');
        it('user authenticates in google');
        it('Screen to choose either link user or register is shown');
        it('User chooses link');
        it('User has to enter loginname');
        it('User has to enter password');
        it('Sucessful linked screen is shown');
        it('No verify email or init user screen is shown');
        it('No mfa prompt is shown');
        it('User is redirected to application (requested redirect url)');
    });

    describe('login with external idp (azure ad) with auto register user', () => {
        it('login policy has enabled login with external idp');
        it('login policy has disabled login with username password');
        it('azure ad idp is added and activated');
        it('auto register is enabled on azure ad');
        it('login policy has register allowed');
        it('login policy mfa prompt lifetime is set to 0 which results in no prompt');
        it('No users with this linked idp does exist');

        it('Auto redirect to azure ad idp (because only possibility to login)');
        it('user authenticates in azure ad');
        it('Screen to choose either link user or register is shown');
        it('User chooses register');
        it('No verify email or init user screen is shown'); // At the moment this is not possible, because azure ad does not send email verified attribute, so we should get a email verify screen but not an init screen
        it('No mfa prompt is shown');
        it('User is redirected to application (requested redirect url)');
    });

});
