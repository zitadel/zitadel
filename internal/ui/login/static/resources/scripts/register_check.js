function CheckRegisterPwPolicy() {
    let policyElement = document.getElementById("register-password");
    let pwNew = policyElement.value;
    let pwNewConfirmation = document.getElementById("register-password-confirmation").value;

    if (ComplexityPolicyCheck(policyElement, pwNew, pwNewConfirmation) === false) {
        policyElement.setAttribute("color", "warn");
        return false;
    } else {
        policyElement.setAttribute("color", "primary");
    }

    return pwNew == pwNewConfirmation;
}

let button = document.getElementById("register-button");
disableSubmit(CheckRegisterPwPolicy, button);
