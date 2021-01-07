function CheckRegisterPwPolicy() {
    let policyElement = document.getElementById("register-password");
    let pwNew = policyElement.value;
    let pwNewConfirmation = document.getElementById("register-password-confirmation").value;

    if (ComplexityPolicyCheck(policyElement, pwNew) === false) {
        return false;
    }

    return pwNew == pwNewConfirmation;
}

let button = document.getElementById("register-button");

// set formfield color primary if all required field are correct
const pwdInput = document.getElementById('register-password');
if (pwdInput) {
    pwdInput.addEventListener('input', () => {
        if (ComplexityPolicyCheck(pwdInput, pwdInput.value)) {
            pwdInput.setAttribute("color", "primary");
        } else {
            pwdInput.setAttribute("color", "warn");
        }
    });
}

disableSubmit(CheckRegisterPwPolicy, button);
