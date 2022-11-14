const pattern =
  /^[a-zA-Z0-9.!#$%&'*+\/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*$/;

const htmlElements = document.querySelectorAll('input[type="email"]');

function validateEmail(element) {
  const valid = pattern.test(element.value);
  if (!valid) {
    element.setAttribute("color", "warn");
  } else {
    element.setAttribute("color", "primary");
  }
}

htmlElements.forEach((element) => {
  element.addEventListener("input", () => validateEmail(element));
});
