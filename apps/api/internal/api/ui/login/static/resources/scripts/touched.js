let inputs = document.getElementsByTagName("input");

if (inputs && inputs.length) {
  for (let input of inputs) {
    input.addEventListener("focus", () => {
      input.classList.add("lgn-focused");
    });
    input.addEventListener("blur", () => {
      input.classList.add("lgn-touched");
      input.classList.remove("lgn-focused");
    });
  }
}
