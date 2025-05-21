document.addEventListener(
  "DOMContentLoaded",
  () => {
    const form = document.getElementsByTagName("form")[0];
    if (form) {
      form.addEventListener("submit", (event) => {
        event.preventDefault(); // Prevent the default form submission
        checkWebauthnSupported(registerCredential);
      });
    }
  }
);

async function registerCredential() {
  document.getElementById("wa-error").classList.add("hidden");

  let opt;
  try {
    opt = JSON.parse(window.atob(document.getElementsByName("credentialCreationData")[0].value));
  } catch (e) {
    webauthnError({ message: "Failed to parse credential creation data." });
    return;
  }

  try {
    opt.publicKey.challenge = bufferDecode(opt.publicKey.challenge, "publicKey.challenge");
    opt.publicKey.user.id = bufferDecode(opt.publicKey.user.id, "publicKey.user.id");
    if (opt.publicKey.excludeCredentials) {
      for (let i = 0; i < opt.publicKey.excludeCredentials.length; i++) {
        if (opt.publicKey.excludeCredentials[i].id !== null) {
          opt.publicKey.excludeCredentials[i].id = bufferDecode(opt.publicKey.excludeCredentials[i].id, "publicKey.excludeCredentials");
        }
      }
    }
  } catch (e) {
    webauthnError({ message: "Failed to decode buffer data." });
    return;
  }

  try {
    const credential = await navigator.credentials.create({
      publicKey: opt.publicKey,
    });

    createCredential(credential);
  } catch (err) {
    webauthnError(err);
  }
}

function createCredential(newCredential) {
  let attestationObject = new Uint8Array(
    newCredential.response.attestationObject
  );
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

  document.getElementsByName("credentialData")[0].value = window.btoa(data);
  document.getElementsByTagName("form")[0].submit();
}
