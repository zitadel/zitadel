function CheckChangePwPolicy() {
  const pwNew = document.getElementById("change-new-password");
  const pwNewValue = pwNew.value;
  const pwNewConfirmation = document.getElementById(
    "change-password-confirmation"
  );
  const pwNewConfirmationValue = pwNewConfirmation.value;

  ComplexityPolicyCheck(pwNew, pwNewConfirmation);

  return pwNewValue == pwNewConfirmationValue;
}

let button = document.getElementById("change-password-button");
disableSubmit(CheckChangePwPolicy, button);
