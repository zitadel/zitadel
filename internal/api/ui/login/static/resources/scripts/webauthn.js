function checkWebauthnSupported(buttonId, func) {
  let support = document.getElementsByClassName("wa-support");
  let noSupport = document.getElementsByClassName("wa-no-support");
  if (!window.PublicKeyCredential) {
    for (let item of noSupport) {
      item.classList.remove("hidden");
    }
    for (let item of support) {
      item.classList.add("hidden");
    }
    return;
  }

  // Check if the browser supports WebAuthn, then add the event listener
  const button = document.getElementById(buttonId);
  if (button) {
    button.addEventListener("click", func);
  } else {
    console.warn(`Button with id '${buttonId}' not found. No event listener added.`);
  }
}

function webauthnError(error) {
  let err = document.getElementById("wa-error");
  let causeElement = err.getElementsByClassName("cause")[0];

  if (error.message) {
    causeElement.innerText = error.message;
  } else if (error.value) {
    causeElement.innerText = error.value;
  } else {
    console.error("Unknown error:", error);
    causeElement.innerText = "An unknown error occurred.";
  }

  err.classList.remove("hidden");
}

function bufferDecode(value, name) {
  return coerceToArrayBuffer(value, name);
}

function bufferEncode(value, name) {
  return coerceToBase64Url(value, name);
}
