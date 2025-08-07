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

let button = document.getElementById("register-button");
disableSubmit(CheckRegisterPwPolicy, button);
