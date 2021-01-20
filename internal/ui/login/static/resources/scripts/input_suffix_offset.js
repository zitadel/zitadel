const suffix = document.getElementById('default-login-suffix');
const suffixInput = document.getElementsByClassName('lgn-suffix-input')[0];

calculateOffset();
suffix.addEventListener("DOMCharacterDataModified", calculateOffset);

function calculateOffset() {
    // add suffix width to inner right padding of the input field
    if (suffix && suffixInput) {
        suffixInput.style.paddingRight = `${(suffix.offsetWidth ?? 0) + 10}px`;
    }
}