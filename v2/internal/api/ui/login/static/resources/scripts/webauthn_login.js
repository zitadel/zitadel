document.addEventListener('DOMContentLoaded', checkWebauthnSupported('btn-login', login));

function login() {
    document.getElementById('wa-error').classList.add('hidden');

    let makeAssertionOptions = JSON.parse(atob(document.getElementsByName('credentialAssertionData')[0].value));
    makeAssertionOptions.publicKey.challenge = bufferDecode(makeAssertionOptions.publicKey.challenge);
    makeAssertionOptions.publicKey.allowCredentials.forEach(function (listItem) {
        listItem.id = bufferDecode(listItem.id)
    });
    navigator.credentials.get({
        publicKey: makeAssertionOptions.publicKey
    }).then(function (credential) {
            verifyAssertion(credential);
        }).catch(function (err) {
            webauthnError(err);
    });
}

function verifyAssertion(assertedCredential) {
    let authData = new Uint8Array(assertedCredential.response.authenticatorData);
    let clientDataJSON = new Uint8Array(assertedCredential.response.clientDataJSON);
    let rawId = new Uint8Array(assertedCredential.rawId);
    let sig = new Uint8Array(assertedCredential.response.signature);
    let userHandle = new Uint8Array(assertedCredential.response.userHandle);

    let data = JSON.stringify({
        id: assertedCredential.id,
        rawId: bufferEncode(rawId),
        type: assertedCredential.type,
        response: {
            authenticatorData: bufferEncode(authData),
            clientDataJSON: bufferEncode(clientDataJSON),
            signature: bufferEncode(sig),
            userHandle: bufferEncode(userHandle),
        },
    })

    document.getElementsByName('credentialData')[0].value = btoa(data);
    document.getElementsByTagName('form')[0].submit();
}