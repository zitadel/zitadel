validate();

document.querySelectorAll('input[name="provider"]').forEach((input) => {
  input.addEventListener("change", validate);
});

function validate() {
  const checkedMfaMethod = document.querySelector(
    'input[name="provider"]:checked'
  );
  const submitButton = document.querySelector(
    'button.lgn-raised-button[type="submit"]'
  );

  if (checkedMfaMethod && submitButton) {
    submitButton.disabled = false;
  } else if (submitButton) {
    submitButton.disabled = true;
  }
}
