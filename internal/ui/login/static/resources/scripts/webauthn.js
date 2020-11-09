function checkWebauthnSupported(button, func) {
    let support = document.getElementsByClassName("wa-support");
    let noSupport = document.getElementsByClassName("wa-no-support");
    if (typeof (PublicKeyCredential) === undefined) {
        for (let item of noSupport) {
            item.classList.remove('hidden');
        }
        for (let item of support) {
            item.classList.add('hidden');
        }
        return
    }
    document.getElementById(button).addEventListener('click', func);
}

function webauthnError(error) {
    let err = document.getElementById('wa-error');
    err.getElementsByClassName('cause')[0].innerText = error.message;
    err.classList.remove('hidden');
}

function bufferDecode(value) {
    return base64js.decode(value);
    // return Uint8Array.from(atob(value), c => c.charCodeAt(0));
}

function bufferEncode(value) {
    return base64js.encode(value)
        .replace(/\+/g, "-")
        .replace(/\//g, "_")
        .replace(/=/g, "");
}
