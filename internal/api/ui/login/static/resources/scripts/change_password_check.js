function CheckChangePwPolicy() {
  let pwNew = document.getElementById("change-new-password");
  let pwNewValue = pwNew.value;
  let pwNewConfirmation = document.getElementById(
    "change-password-confirmation"
  );
  let pwNewConfirmationValue = pwNewConfirmation.value;

  if (
    ComplexityPolicyCheck(pwNew, pwNewValue, pwNewConfirmationValue) === false
  ) {
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

  return pwNewValue === pwNewConfirmationValue;
}

let button = document.getElementById("change-password-button");
disableSubmit(CheckChangePwPolicy, button);
