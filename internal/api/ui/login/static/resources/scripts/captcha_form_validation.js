function onRecaptchaLoad() {
  let button = document.getElementById("submit-button");
  let recaptchaContainer = document.querySelector(".g-recaptcha");

  if (!recaptchaContainer) {
    console.error("No reCAPTCHA container found");
    return;
  }

  // Render ReCaptcha in element with id `captcha-container`
  // Add custom event disptch using grecaptcha callbacks to trigger toggleButton()
  let recaptchaId = grecaptcha.enterprise.render(recaptchaContainer, {
    sitekey: recaptchaContainer.getAttribute("data-sitekey"),
    callback: function () {
      window.dispatchEvent(new Event("captchaChanged"));
    },
    "expired-callback": function () {
      window.dispatchEvent(new Event("captchaChanged"));
    },
    "error-callback": function () {
      window.dispatchEvent(new Event("captchaChanged"));
    },
  });

  function captchaCheck() {
    return grecaptcha.enterprise.getResponse(recaptchaId) !== "";
  }

  disableSubmit(captchaCheck, button);
}

// Check if ReCaptcha is already loaded
if (typeof grecaptcha !== "undefined") {
  onRecaptchaLoad();
} else {
  window.addEventListener("load", function () {
    if (typeof grecaptcha !== "undefined") {
      onRecaptchaLoad();
    } else {
      console.error("grecaptcha is still undefined after window load");
    }
  });
}
