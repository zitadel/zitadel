// if an autofilled input is deleted we remove the attribute
function detectDelete(event) {
  const key = event.key;
  if (key === "Backspace" || key === "Delete") {
    event.target.isAutofilled = false;
  }
}

// if the autofill associated animation is detected we add a property
// and check if submit button should be disabled or not
function autofill(target, checks, form, inputs, button) {
  if (!target.isAutofilled) {
    target.isAutofilled = true;
    target.dispatchEvent(new CustomEvent("autofill", { bubbles: true }));
    toggleButton(checks, form, inputs, button);
  }
}

function disableSubmit(checks, button) {
  let form = document.getElementsByTagName("form")[0];
  let inputs = form.getElementsByTagName("input");
  if (button) {
    button.disabled = true;
  }
  addRequiredEventListener(inputs, checks, form, button);
  disableDoubleSubmit(form, button);

  toggleButton(checks, form, inputs, button);
}

function addRequiredEventListener(inputs, checks, form, button) {
  let eventType = "input";
  for (i = 0; i < inputs.length; i++) {
    if (inputs[i].required) {
      eventType = "input";
      if (inputs[i].type === "checkbox") {
        eventType = "click";
      }

      inputs[i].addEventListener(eventType, function () {
        toggleButton(checks, form, inputs, button);
      });

      if (inputs[i].type !== "checkbox") {
        // hack for Chrome, add an animationstart event listener
        // if input is autofilled: https://gist.github.com/jonathantneal/d462fc2bf761a10c9fca60eb634f6977?permalink_comment_id=2901919
        inputs[i].addEventListener("animationstart", (event) =>
          autofill(event.target, checks, form, inputs, button)
        );

        inputs[i].addEventListener("keydown", detectDelete);
      }
    }
  }
}

function disableDoubleSubmit(form, button) {
  form.addEventListener("submit", function () {
    document.body.classList.add("waiting");
    button.disabled = true;
  });
}

function toggleButton(checks, form, inputs, button) {
  if (checks !== undefined) {
    if (checks() === false) {
      button.disabled = true;
      return;
    }
  }
  const targetValue = !allRequiredDone(form, inputs);
  button.disabled = targetValue;
}

function allRequiredDone(form, inputs) {
  for (i = 0; i < inputs.length; i++) {
    if (inputs[i].required) {
      if (inputs[i].type === "checkbox" && !inputs[i].checked) {
        return false;
      }
      if (inputs[i].value === "" && !inputs[i].isAutofilled) {
        return false;
      }
    }
  }
  return true;
}
