function CheckChangePwPolicy() {
    let policyElement = document.getElementById("change-new-password");
    let pwNew = policyElement.value;
    let pwNewConfirmation = document.getElementById("change-password-confirmation").value;

    if (ComplexityPolicyCheck(policyElement, pwNew) === false) {
        policyElement.setAttribute("color", "warn");
        return false;
    } else {
        policyElement.setAttribute("color", "primary");
    }

    return pwNew == pwNewConfirmation;
}

let button = document.getElementById("change-password-button");
disableSubmit(CheckChangePwPolicy, button);