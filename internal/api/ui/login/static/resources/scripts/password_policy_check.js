function ComplexityPolicyCheck(policyElement, pwNew, pwNewConfirmation) {
    let minLength = policyElement.dataset.minlength;
    let upperRegex = policyElement.dataset.hasUppercase;
    let lowerRegex = policyElement.dataset.hasLowercase;
    let numberRegex = policyElement.dataset.hasNumber;
    let symbolRegex = policyElement.dataset.hasSymbol;

    let invalid = 0;

    let minlengthelem = document.getElementById('minlength');
    if (pwNew.length >= minLength) {
        ValidPolicy(minlengthelem);
    } else {
        InvalidPolicy(minlengthelem);
        invalid++;
    }
    let upper = document.getElementById('uppercase');
    if (upperRegex !== "") {
        if (RegExp(upperRegex).test(pwNew)) {
            ValidPolicy(upper);
        } else {
            InvalidPolicy(upper);
            invalid++;
        }
    }
    let lower = document.getElementById('lowercase');
    if (lowerRegex !== "") {
        if (RegExp(lowerRegex).test(pwNew)) {
            ValidPolicy(lower);
        } else {
            InvalidPolicy(lower);
            invalid++;
        }
    }
    let number = document.getElementById('number');
    if (numberRegex !== "") {
        if (RegExp(numberRegex).test(pwNew)) {
            ValidPolicy(number);
        } else {
            InvalidPolicy(number);
            invalid++;
        }
    }
    let symbol = document.getElementById('symbol');
    if (symbolRegex !== "") {
        if (RegExp(symbolRegex).test(pwNew)) {
            ValidPolicy(symbol);
        } else {
            InvalidPolicy(symbol);
            invalid++;
        }
    }
    let confirmation = document.getElementById('confirmation');
    if (pwNew === pwNewConfirmation && pwNewConfirmation !== "" ) {
        ValidPolicy(confirmation);
    } else {
        InvalidPolicy(confirmation);
        invalid++;
    }
    return invalid===0;
}

function ValidPolicy(element) {
    element.classList.remove('invalid');
    element.getElementsByTagName('i')[0].classList.remove('lgn-icon-times-solid');
    element.getElementsByTagName('i')[0].classList.remove('lgn-warn');
    element.getElementsByTagName('i')[0].classList.add('lgn-icon-check-solid');
    element.getElementsByTagName('i')[0].classList.add('lgn-valid');
}

function InvalidPolicy(element) {
    element.classList.add('invalid');
    element.getElementsByTagName('i')[0].classList.remove('lgn-valid');
    element.getElementsByTagName('i')[0].classList.remove('lgn-icon-check-solid');
    element.getElementsByTagName('i')[0].classList.add('lgn-warn');
    element.getElementsByTagName('i')[0].classList.add('lgn-icon-times-solid');
}
