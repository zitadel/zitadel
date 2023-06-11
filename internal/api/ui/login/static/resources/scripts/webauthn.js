function checkWebauthnSupported(button, func) {
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
  document.getElementById(button).addEventListener("click", func);
}

function webauthnError(error) {
  let err = document.getElementById("wa-error");
  err.getElementsByClassName("cause")[0].innerText = error.message;
  err.classList.remove("hidden");
}

function bufferDecode(value, name) {
  return coerceToArrayBuffer(value, name);
}

function bufferEncode(value, name) {
  return coerceToBase64Url(value, name);
}
