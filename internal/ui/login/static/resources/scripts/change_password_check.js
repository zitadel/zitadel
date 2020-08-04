function CheckChangePwPolicy() {
    let policyElement = document.getElementById("change-new-password")
    let pwNew = policyElement.value;
    let pwNewConfirmation = document.getElementById("change-password-confirmation").value;

    if (ComplexityPolicyCheck(policyElement, pwNew) === false) {
        return false;
    }

    return pwNew == pwNewConfirmation;
}

let button = document.getElementById("change-password-button");
disableSubmit(CheckChangePwPolicy, button);