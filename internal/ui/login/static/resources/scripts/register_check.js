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
disableSubmit(CheckRegisterPwPolicy, button);
