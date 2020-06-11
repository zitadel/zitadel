const copyToClipboard = str => {
    navigator.clipboard.writeText(str);
}

let copyButton = document.getElementsByClassName("copy")[0];
copyButton.addEventListener("click", copyToClipboard(copyButton.getAttribute("data-copy").value));
