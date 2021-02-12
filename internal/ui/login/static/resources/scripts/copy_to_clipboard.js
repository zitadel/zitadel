const copyToClipboard = str => {
    navigator.clipboard.writeText(str);
};

let copyButton = document.getElementById("copy");
copyButton.addEventListener("click", copyToClipboard(copyButton.getAttribute("data-copy")));
