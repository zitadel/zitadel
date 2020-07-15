function ComplexityPolicyCheck(button, pwNew) {
    console.log("ComplexityPolicyCheck");
    let minLength = document.getElementById("min-length").value;
    let upperRegex = document.getElementById("uppercase-regex").value;
    let lowerRegex = document.getElementById("lowercase-regex").value;
    let numberRegex = document.getElementById("hasnumebr-regex").value;
    let symbolRegex = document.getElementById("hassymbol-regex").value;

    let minlengthelem = document.getElementById('minlength')
    if (pwNew.length >= minLength) {
        ValidPolicy(minlengthelem);
    } else {
        InvalidPolicy(minlengthelem, button);
    }
    let upper = document.getElementById('uppercase')
    if (upperRegex !== "") {
        if (RegExp(upperRegex).test(pwNew)) {
            ValidPolicy(upper);
        } else {
            button.classList.add("disabled");
            InvalidPolicy(upper, button);
        }
    }
    let lower = document.getElementById('lowercase')
    if (lowerRegex !== "") {
        if (RegExp(lowerRegex).test(pwNew)) {
            ValidPolicy(lower);
        } else {
            InvalidPolicy(lower, button);
        }
    }
    let number = document.getElementById('number')
    if (numberRegex != "") {
       if (RegExp(numberRegex).test(pwNew)) {
           ValidPolicy(number);
        } else {
           button.classList.add("disabled");
           InvalidPolicy(number, button);
        }
    }
    let symbol = document.getElementById('symbol')
    if (symbolRegex != "") {
        if (RegExp(symbolRegex).test(pwNew)) {
            ValidPolicy(symbol);
        } else {
            button.classList.add("disabled");
            InvalidPolicy(symbol, button);
        }
    }
}
function ValidPolicy(element) {
    element.classList.remove('invalid')
    element.getElementsByTagName('i')[0].innerText = 'check';
}

function InvalidPolicy(element, button) {
    button.classList.add("disabled");
    element.classList.add('invalid')
    element.getElementsByTagName('i')[0].innerText = 'clear';
}
