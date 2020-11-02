function bufferDecode(value) {
    return Uint8Array.from(atob(value), c => c.charCodeAt(0));
}

function bufferEncode(value) {
    return base64js.fromByteArray(value)
        .replace(/\+/g, "-")
        .replace(/\//g, "_")
        .replace(/=/g, "");
}

function registerCredential() {
    let opt = JSON.parse(atob(document.getElementsByName('credentialCreationData')[0].value));
    opt.publicKey.challenge = bufferDecode(opt.publicKey.challenge);
    opt.publicKey.user.id = bufferDecode(opt.publicKey.user.id);
    if (opt.publicKey.excludeCredentials) {
        for (let i = 0; i < opt.publicKey.excludeCredentials.length; i++) {
            opt.publicKey.excludeCredentials[i].id = bufferDecode(opt.publicKey.excludeCredentials[i].id);
        }
    }
    console.log(opt);
    navigator.credentials.create({
        publicKey: opt.publicKey
    }).then(function (credential) {
        console.log(credential);
        createCredential(credential);
    }).catch(function (err) {
        alert(err.name);
        alert(err.message);
    });
}

function createCredential(newCredential) {
    let attestationObject = new Uint8Array(newCredential.response.attestationObject);
    let clientDataJSON = new Uint8Array(newCredential.response.clientDataJSON);
    let rawId = new Uint8Array(newCredential.rawId);

    let data = JSON.stringify({
        id: newCredential.id,
        rawId: bufferEncode(rawId),
        type: newCredential.type,
        response: {
            attestationObject: bufferEncode(attestationObject),
            clientDataJSON: bufferEncode(clientDataJSON),
        },
    });

    document.getElementsByName('credentialData')[0].value = btoa(data);
    document.getElementsByTagName('form')[0].submit();
}

function login() {
    let makeAssertionOptions = JSON.parse(atob(document.getElementsByName('credentialAssertionData')[0].value));
    console.log("Assertion Options:");
    console.log(makeAssertionOptions);
    makeAssertionOptions.publicKey.challenge = bufferDecode(makeAssertionOptions.publicKey.challenge);
    makeAssertionOptions.publicKey.allowCredentials.forEach(function (listItem) {
        listItem.id = bufferDecode(listItem.id)
    });
    console.log(makeAssertionOptions);
    navigator.credentials.get({
        publicKey: makeAssertionOptions.publicKey
    })
        .then(function (credential) {
            console.log(credential);
            verifyAssertion(credential);
        }).catch(function (err) {
        alert(err.name);
        alert(err.message);
    });
}


function verifyAssertion(assertedCredential) {
    // Move data into Arrays incase it is super long
    console.log('calling verify')
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
