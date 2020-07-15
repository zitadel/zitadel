function CheckRegisterPwPolicy() {
    let pwNew = document.getElementById("register-password").value;
    let pwNewConfirmation = document.getElementById("register-password-confirmation").value;
    let button = document.getElementById("register-button");

    if (pwNew == "" || pwNewConfirmation == "" || pwNew != pwNewConfirmation) {
        button.classList.add("disabled");
    } else {
        button.classList.remove("disabled");
    }

    ComplexityPolicyCheck(button, pwNew);
}

let newPWRegister = document.getElementById("register-password");
newPWRegister.addEventListener('input', CheckRegisterPwPolicy);

let newPWConfirmationRegister = document.getElementById("register-password-confirmation");
newPWConfirmationRegister.addEventListener('input', CheckRegisterPwPolicy);
