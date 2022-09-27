function RenderDefaultLoginnameSuffix() {
    let orgNameText = document.getElementById("orgname").value;
    let userName = document.getElementById("username");
    let defaultLoginNameSuffix = document.getElementById("default-login-suffix");

    let iamDomain = userName.dataset.iamDomain;
    let orgDomain = orgNameText.replace(" ", "-");
    if (orgDomain !== "") {
        defaultLoginNameSuffix.innerText = "@" + orgDomain.toLowerCase() + "." + iamDomain;
    } else {
        defaultLoginNameSuffix.innerText = "";
    }

    offsetLabel();
}

window.addEventListener('DOMContentLoaded', (event) => {
    RenderDefaultLoginnameSuffix();
});

document.getElementById("orgname").addEventListener('input', function () {
    RenderDefaultLoginnameSuffix();
});

function offsetLabel() {
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
}