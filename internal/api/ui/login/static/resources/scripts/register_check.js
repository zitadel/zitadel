function CheckRegisterPwPolicy() {
  const pwNew = document.getElementById("register-password");
  const pwNewConfirmation = document.getElementById(
    "register-password-confirmation"
  );
  const pwNewValue = pwNew.value;
  const pwNewConfirmationValue = pwNewConfirmation.value;

  ComplexityPolicyCheck(pwNew, pwNewConfirmation);

  return pwNewValue == pwNewConfirmationValue;
}

function onRecaptchaLoad() {
  let button = document.getElementById("register-button");
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

  disableSubmit([CheckRegisterPwPolicy, captchaCheck], button);
}

if (document.getElementById("captcha") !== null) {
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
} else {
  let button = document.getElementById("register-button");
  disableSubmit(CheckRegisterPwPolicy, button);
}
