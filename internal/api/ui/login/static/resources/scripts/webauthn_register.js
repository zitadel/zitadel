document.addEventListener(
  "DOMContentLoaded",
  checkWebauthnSupported("btn-register", registerCredential)
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
  try {
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

    console.log("Encoded data:", data);

    let credentialDataElement = document.getElementsByName("credentialData")[0];
    if (!credentialDataElement) {
      console.error("Element with name 'credentialData' not found.");
      webauthnError({ message: "Element with name 'credentialData' not found." });
      return;
    }

    credentialDataElement.value = window.btoa(data);
    console.log("Credential data set:", credentialDataElement.value);

    let form = document.getElementsByTagName("form")[0];
    if (!form) {
      console.error("Form element not found.");
      webauthnError({ message: "Form element not found." });
      return;
    }

    console.log("Submitting form...");
    form.submit();
  } catch (err) {
    webauthnError(err);
  }
}
