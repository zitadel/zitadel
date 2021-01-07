function CheckChangePwPolicy() {
    let policyElement = document.getElementById("change-new-password");
    let pwNew = policyElement.value;
    let pwNewConfirmation = document.getElementById("change-password-confirmation").value;

    if (ComplexityPolicyCheck(policyElement, pwNew) === false) {
        return false;
    }

    return pwNew == pwNewConfirmation;
}

let button = document.getElementById("change-password-button");

// set formfield color primary if all required field are correct
const pwdInput = document.getElementById('change-new-password');
if (pwdInput) {
    pwdInput.addEventListener('input', () => {
        if (ComplexityPolicyCheck(pwdInput, pwdInput.value)) {
            pwdInput.setAttribute("color", "primary");
        } else {
            pwdInput.setAttribute("color", "warn");
        }
    });
}

disableSubmit(CheckChangePwPolicy, button);