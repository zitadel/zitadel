function CheckRegisterPwPolicy() {
  const pwNew = document.getElementById("register-password");
  const pwNewConfirmation = document.getElementById(
    "register-password-confirmation"
  );
  const pwNewValue = pwNew.value;
  const pwNewConfirmationValue = pwNewConfirmation.value;

  if (ComplexityPolicyCheck(pwNew) === false) {
    pwNew.setAttribute("color", "warn");
    return false;
  } else {
    pwNew.setAttribute("color", "primary");
  }

  if (pwNewValue !== pwNewConfirmationValue && pwNewConfirmationValue !== "") {
    pwNewConfirmation.setAttribute("color", "warn");
  } else {
    pwNewConfirmation.setAttribute("color", "primary");
  }

  return pwNewValue == pwNewConfirmationValue;
}

let button = document.getElementById("register-button");
disableSubmit(CheckRegisterPwPolicy, button);
