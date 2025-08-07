function ComplexityPolicyCheck(passwordElement, passwordConfirmationElement) {
  const minLength = passwordElement.dataset.minlength;
  const upperRegex = passwordElement.dataset.hasUppercase;
  const lowerRegex = passwordElement.dataset.hasLowercase;
  const numberRegex = passwordElement.dataset.hasNumber;
  const symbolRegex = passwordElement.dataset.hasSymbol;

  let invalid = 0;

  const minLengthElem = document.getElementById("minlength");
  if (passwordElement.value.length >= minLength) {
    ValidPolicy(minLengthElem);
  } else {
    InvalidPolicy(minLengthElem);
    invalid++;
  }

  const maxLengthElem = document.getElementById("maxlength");
  if (passwordElement.value.length <= 70) {
    ValidPolicyFlipped(maxLengthElem);
  } else {
    InvalidPolicyFlipped(maxLengthElem);
    invalid++;
  }

  const upper = document.getElementById("uppercase");
  if (upperRegex !== "") {
    if (RegExp(upperRegex).test(passwordElement.value)) {
      ValidPolicy(upper);
    } else {
      InvalidPolicy(upper);
      invalid++;
    }
  }

  const lower = document.getElementById("lowercase");
  if (lowerRegex !== "") {
    if (RegExp(lowerRegex).test(passwordElement.value)) {
      ValidPolicy(lower);
    } else {
      InvalidPolicy(lower);
      invalid++;
    }
  }

  const number = document.getElementById("number");
  if (numberRegex !== "") {
    if (RegExp(numberRegex).test(passwordElement.value)) {
      ValidPolicy(number);
    } else {
      InvalidPolicy(number);
      invalid++;
    }
  }

  const symbol = document.getElementById("symbol");
  if (symbolRegex !== "") {
    if (RegExp(symbolRegex).test(passwordElement.value)) {
      ValidPolicy(symbol);
    } else {
      InvalidPolicy(symbol);
      invalid++;
    }
  }

  const confirmation = document.getElementById("confirmation");
  if (
    passwordElement.value === passwordConfirmationElement.value &&
    passwordConfirmationElement.value !== ""
  ) {
    ValidPolicy(confirmation);
    passwordConfirmationElement.setAttribute("color", "primary");
  } else {
    InvalidPolicy(confirmation);
    passwordConfirmationElement.setAttribute("color", "warn");
  }

  if (invalid > 0) {
    passwordElement.setAttribute("color", "warn");
    return false;
  } else {
    passwordElement.setAttribute("color", "primary");
    return true;
  }
}

function ValidPolicy(element) {
  element.classList.remove("invalid");
  element.getElementsByTagName("i")[0].classList.remove("lgn-icon-times-solid");
  element.getElementsByTagName("i")[0].classList.remove("lgn-warn");
  element.getElementsByTagName("i")[0].classList.add("lgn-icon-check-solid");
  element.getElementsByTagName("i")[0].classList.add("lgn-valid");
}

function ValidPolicyFlipped(element) {
  element.classList.add("valid");
  element.getElementsByTagName("i")[0].classList.remove("lgn-warn");
  element.getElementsByTagName("i")[0].classList.remove("lgn-icon-times-solid");
  element.getElementsByTagName("i")[0].classList.add("lgn-valid");
  element.getElementsByTagName("i")[0].classList.add("lgn-icon-check-solid");
}

function InvalidPolicy(element) {
  element.classList.add("invalid");
  element.getElementsByTagName("i")[0].classList.remove("lgn-valid");
  element.getElementsByTagName("i")[0].classList.remove("lgn-icon-check-solid");
  element.getElementsByTagName("i")[0].classList.add("lgn-warn");
  element.getElementsByTagName("i")[0].classList.add("lgn-icon-times-solid");
}

function InvalidPolicyFlipped(element) {
  element.classList.remove("valid");
  element.getElementsByTagName("i")[0].classList.remove("lgn-icon-check-solid");
  element.getElementsByTagName("i")[0].classList.remove("lgn-valid");
  element.getElementsByTagName("i")[0].classList.add("lgn-icon-times-solid");
  element.getElementsByTagName("i")[0].classList.add("lgn-warn");
}
