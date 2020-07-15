function CheckChangePwPolicy() {
    let policyElement = document.getElementById("change-new-password")
    let pwNew = policyElement.value;
    let pwNewConfirmation = document.getElementById("change-password-confirmation").value;
    let oldPW = document.getElementById("change-old-password").value;
    let button = document.getElementById("change-password-button");

    if (oldPW == "" || pwNew == "" || pwNewConfirmation == "" || pwNew != pwNewConfirmation) {
        button.classList.add("disabled");
    } else {
        button.classList.remove("disabled");
    }

    ComplexityPolicyCheck(policyElement, button, pwNew);
}

let newPWChange = document.getElementById("change-new-password");
newPWChange.addEventListener('input', CheckChangePwPolicy);

let newPWConfirmationChange = document.getElementById("change-password-confirmation");
newPWConfirmationChange.addEventListener('input', CheckChangePwPolicy);
