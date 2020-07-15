
function ComplexityPolicyCheck() {
    console.log("GUGUS");
    let pwNew = document.getElementById("new-password").value;
    let pwNewConfirmation = document.getElementById("password-confirmation").value;
    let oldPW = document.getElementById("old-password").value;
    let minLength = document.getElementById("min-length").value;
    let upperRegex = document.getElementById("uppercase-regex").value;
    let lowerRegex = document.getElementById("lowercase-regex").value;
    let numberRegex = document.getElementById("hasnumebr-regex").value;
    let symbolRegex = document.getElementById("hassymbol-regex").value;

    if (oldPW == "" || pwNew == "" || pwNewConfirmation == "" || pwNew != pwNewConfirmation) {
        document.getElementById("change-button").classList.add("disabled");
    } else {
        document.getElementById("change-button").classList.remove("disabled");
    }


    let minlengthelem = document.getElementById('minlength')
    if (pwNew.length >= minLength) {
        ValidPolicy(minlengthelem);
    } else {
        InvalidPolicy(minlengthelem);
    }
    let upper = document.getElementById('uppercase')
    if (upperRegex !== "") {
        if (RegExp(upperRegex).test(pwNew)) {
            ValidPolicy(upper);
        } else {
            document.getElementById("change-button").classList.add("disabled");
            InvalidPolicy(upper);
        }
    }
    let lower = document.getElementById('lowercase')
    if (lowerRegex !== "") {
        if (RegExp(lowerRegex).test(pwNew)) {
            ValidPolicy(lower);
        } else {
            InvalidPolicy(lower);
        }
    }
    let number = document.getElementById('number')
    if (numberRegex != "") {
       if (RegExp(numberRegex).test(pwNew)) {
           ValidPolicy(number);
        } else {
           document.getElementById("change-button").classList.add("disabled");
           InvalidPolicy(number);
        }
    }
    let symbol = document.getElementById('symbol')
    if (symbolRegex != "") {
        if (RegExp(symbolRegex).test(pwNew)) {
            ValidPolicy(symbol);
        } else {
            document.getElementById("change-button").classList.add("disabled");
            InvalidPolicy(symbol);
        }
    }
}

function InvalidPolicy(element) {
    document.getElementById("change-button").classList.add("disabled");
    element.classList.add('invalid')
    element.getElementsByTagName('i')[0].innerText = 'clear';
}

function ValidPolicy(element) {
    element.classList.remove('invalid')
    element.getElementsByTagName('i')[0].innerText = 'check';
}

let newPW = document.getElementById("new-password");
newPW.addEventListener('input', ComplexityPolicyCheck);

let newPWConfirmation = document.getElementById("password-confirmation");
newPWConfirmation.addEventListener('input', ComplexityPolicyCheck);

