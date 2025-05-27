function checkWebauthnSupported(func, optionalClickId) {
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

  // if id is provided add click event only, otherwise call the function directly
  if (optionalClickId) {
    document.getElementById(optionalClickId).addEventListener("click", func);
  } else {
    func();
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
