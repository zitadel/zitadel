function bufferDecode(value) {
    return Uint8Array.from(atob(value), c => c.charCodeAt(0));
}

function bufferEncode(value) {
    return base64js.fromByteArray(value)
        .replace(/\+/g, "-")
        .replace(/\//g, "_")
        .replace(/=/g, "");
}

function register() {
    let opt = JSON.parse(atob(document.getElementsByName('credentialCreationData')[0].value));
    opt.publicKey.challenge = bufferDecode(opt.publicKey.challenge);
    opt.publicKey.user.id = bufferDecode(opt.publicKey.user.id);
    console.log(opt);
    navigator.credentials.create({
        publicKey: opt.publicKey
    }).then(function (credential) {
        console.log(credential);
        createCredential(credential);
    }).catch(function (err) {
        console.log(err.name);
    });
}

function createCredential(newCredential) {
    let attestationObject = new Uint8Array(newCredential.response.attestationObject);
    let clientDataJSON = new Uint8Array(newCredential.response.clientDataJSON);
    let rawId = new Uint8Array(newCredential.rawId);
    //
    // let a = btoa(rawId);
    // let aa = bufferEncode(rawId);
    // console.log(a, aa);
    //
    let data = JSON.stringify({
        id: newCredential.id,
        rawId: bufferEncode(rawId),
        type: newCredential.type,
        response: {
            attestationObject: bufferEncode(attestationObject),
            clientDataJSON: bufferEncode(clientDataJSON),
        },
    })
    // let j = JSON.stringify(newCredential);
    let b = btoa(data)
    document.getElementsByName('credentialData')[0].value = b;
    document.getElementsByTagName('form')[0].submit();
}
