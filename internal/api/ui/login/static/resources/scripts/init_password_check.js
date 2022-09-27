function CheckInitPwPolicy() {
    let policyElement = document.getElementById("password");
    let pwNew = policyElement.value;
    let pwNewConfirmation = document.getElementById("passwordconfirm").value;

    if (ComplexityPolicyCheck(policyElement, pwNew, pwNewConfirmation) === false) {
        return false;
    }

    return pwNew == pwNewConfirmation;
}

let button = document.getElementById("init-button");
disableSubmit(CheckInitPwPolicy, button);