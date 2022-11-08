function ComplexityPolicyCheck(passwordElement) {
  let minLength = passwordElement.dataset.minlength;
  let upperRegex = passwordElement.dataset.hasUppercase;
  let lowerRegex = passwordElement.dataset.hasLowercase;
  let numberRegex = passwordElement.dataset.hasNumber;
  let symbolRegex = passwordElement.dataset.hasSymbol;

  let invalid = 0;

  let minLengthElem = document.getElementById("minlength");
  if (passwordElement.value.length >= minLength) {
    ValidPolicy(minlengthelem);
  } else {
    InvalidPolicy(minlengthelem);
    invalid++;
  }
  let upper = document.getElementById("uppercase");
  if (upperRegex !== "") {
    if (RegExp(upperRegex).test(passwordElement.value)) {
      ValidPolicy(upper);
    } else {
      InvalidPolicy(upper);
      invalid++;
    }
  }
  let lower = document.getElementById("lowercase");
  if (lowerRegex !== "") {
    if (RegExp(lowerRegex).test(passwordElement.value)) {
      ValidPolicy(lower);
    } else {
      InvalidPolicy(lower);
      invalid++;
    }
  }
  let number = document.getElementById("number");
  if (numberRegex !== "") {
    if (RegExp(numberRegex).test(passwordElement.value)) {
      ValidPolicy(number);
    } else {
      InvalidPolicy(number);
      invalid++;
    }
  }
  let symbol = document.getElementById("symbol");
  if (symbolRegex !== "") {
    if (RegExp(symbolRegex).test(passwordElement.value)) {
      ValidPolicy(symbol);
    } else {
      InvalidPolicy(symbol);
      invalid++;
    }
  }

  return invalid === 0;
}

function ValidPolicy(element) {
  element.classList.remove("invalid");
  element.getElementsByTagName("i")[0].classList.remove("lgn-icon-times-solid");
  element.getElementsByTagName("i")[0].classList.remove("lgn-warn");
  element.getElementsByTagName("i")[0].classList.add("lgn-icon-check-solid");
  element.getElementsByTagName("i")[0].classList.add("lgn-valid");
}

function InvalidPolicy(element) {
  element.classList.add("invalid");
  element.getElementsByTagName("i")[0].classList.remove("lgn-valid");
  element.getElementsByTagName("i")[0].classList.remove("lgn-icon-check-solid");
  element.getElementsByTagName("i")[0].classList.add("lgn-warn");
  element.getElementsByTagName("i")[0].classList.add("lgn-icon-times-solid");
}
