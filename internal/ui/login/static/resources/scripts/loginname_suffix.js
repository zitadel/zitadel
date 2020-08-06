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
}

window.addEventListener('DOMContentLoaded', (event) => {
    RenderDefaultLoginnameSuffix();
});

document.getElementById("orgname").addEventListener('input', function () {
    RenderDefaultLoginnameSuffix();
});
