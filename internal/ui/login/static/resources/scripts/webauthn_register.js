document.getElementById('btn-register').addEventListener('click', function () {
    registerCredential();
});

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