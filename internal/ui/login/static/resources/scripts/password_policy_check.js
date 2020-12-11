function ComplexityPolicyCheck(policyElement, pwNew) {
    let minLength = policyElement.dataset.minlength;
    let upperRegex = policyElement.dataset.hasUppercase;
    let lowerRegex = policyElement.dataset.hasLowercase;
    let numberRegex = policyElement.dataset.hasNumber;
    let symbolRegex = policyElement.dataset.hasSymbol;

    let valid = true;

    let minlengthelem = document.getElementById('minlength');
    if (pwNew.length >= minLength) {
        ValidPolicy(minlengthelem);
        valid = true;
    } else {
        InvalidPolicy(minlengthelem);
        valid = false;
    }
    let upper = document.getElementById('uppercase');
    if (upperRegex !== "") {
        if (RegExp(upperRegex).test(pwNew)) {
            ValidPolicy(upper);
            valid = true;
        } else {
            InvalidPolicy(upper);
            valid = false;
        }
    }
    let lower = document.getElementById('lowercase');
    if (lowerRegex !== "") {
        if (RegExp(lowerRegex).test(pwNew)) {
            ValidPolicy(lower);
            valid = true;
        } else {
            InvalidPolicy(lower);
            valid = false;
        }
    }
    let number = document.getElementById('number');
    if (numberRegex != "") {
        if (RegExp(numberRegex).test(pwNew)) {
            ValidPolicy(number);
            valid = true;
        } else {
            InvalidPolicy(number);
            valid = false;
        }
    }
    let symbol = document.getElementById('symbol');
    if (symbolRegex != "") {
        if (RegExp(symbolRegex).test(pwNew)) {
            ValidPolicy(symbol);
            valid = true;
        } else {
            InvalidPolicy(symbol);
            valid = false;
        }
    }
    return valid;
}
function ValidPolicy(element) {
    element.classList.remove('invalid');
    document.getElementById('change-new-password').setAttribute("color", "primary");
    element.getElementsByTagName('i')[0].classList.remove('la-times');
    element.getElementsByTagName('i')[0].classList.remove('lgn-warn');
    element.getElementsByTagName('i')[0].classList.add('la-check');
    element.getElementsByTagName('i')[0].classList.add('lgn-valid');

    // element.getElementsByTagName('i')[0].innerText = 'check';
}

function InvalidPolicy(element) {
    element.classList.add('invalid');
    document.getElementById('change-new-password').setAttribute("color", "warn");
    element.getElementsByTagName('i')[0].classList.remove('lgn-valid');
    element.getElementsByTagName('i')[0].classList.remove('la-check');
    element.getElementsByTagName('i')[0].classList.add('lgn-warn');
    element.getElementsByTagName('i')[0].classList.add('la-times');
    // element.getElementsByTagName('i')[0].innerText = 'clear';
}
