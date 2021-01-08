const copyToClipboard = str => {
    navigator.clipboard.writeText(str);
    // const icon = copyButton.getElementsByClassName('las')[0];
    // if (icon && icon.classList.contains('la-clipboard')) {
    //     icon.remove('la-clipboard');
    //     icon.classList.add('la-clipboard-check');

    //     setTimeout(() => {
    //         if (icon.classList.contains('la-clipboard-check')) {
    //             icon.classList.remove('la-clipboard-check');
    //             icon.classList.add('la-clipboard');
    //         }
    //     }, 3000);
    // }
};

let copyButton = document.getElementById("copy");
copyButton.addEventListener("click", copyToClipboard(copyButton.getAttribute("data-copy")));
