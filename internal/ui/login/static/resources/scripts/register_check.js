function CheckRegisterPwPolicy() {
    let policyElement = document.getElementById("register-password");
    let pwNew = policyElement.value;
    let pwNewConfirmation = document.getElementById("register-password-confirmation").value;
    let button = document.getElementById("register-button");

    if (pwNew == "" || pwNewConfirmation == "" || pwNew != pwNewConfirmation) {
        button.disabled = true;
    } else {
        button.disabled = false;
    }

    ComplexityPolicyCheck(policyElement, button, pwNew);
}

let newPWRegister = document.getElementById("register-password");
newPWRegister.addEventListener('input', CheckRegisterPwPolicy);

let newPWConfirmationRegister = document.getElementById("register-password-confirmation");
newPWConfirmationRegister.addEventListener('input', CheckRegisterPwPolicy);
