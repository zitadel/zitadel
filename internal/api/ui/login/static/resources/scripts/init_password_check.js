function CheckInitPwPolicy() {
  let pwNew = document.getElementById("password");
  let pwNewValue = pwNew.value;
  const pwNewConfirmation = document.getElementById("passwordconfirm");
  let pwNewConfirmationValue = pwNewConfirmation.value;

  if (ComplexityPolicyCheck(pwNew) === false) {
    return false;
  }

  return pwNewValue == pwNewConfirmationValue;
}

let button = document.getElementById("init-button");
disableSubmit(CheckInitPwPolicy, button);
