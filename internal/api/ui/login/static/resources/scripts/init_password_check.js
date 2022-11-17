function CheckInitPwPolicy() {
  const pwNew = document.getElementById("password");
  const pwNewValue = pwNew.value;
  const pwNewConfirmation = document.getElementById("passwordconfirm");
  const pwNewConfirmationValue = pwNewConfirmation.value;

  ComplexityPolicyCheck(pwNew, pwNewConfirmation);

  return pwNewValue == pwNewConfirmationValue;
}

let button = document.getElementById("init-button");
disableSubmit(CheckInitPwPolicy, button);
