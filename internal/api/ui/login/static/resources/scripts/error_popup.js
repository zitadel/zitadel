setTimeout(removeOverlay, 3000);
window.addEventListener('DOMContentLoaded', (event) => {
    document.addEventListener("mousemove", (event) => {removeOverlay()});
    document.addEventListener("click", (event) => {removeOverlay()});
})

function removeOverlay() {
    document.getElementById("dialog_overlay").classList.remove("show");
}