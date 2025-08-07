function removeOverlay(overlay) {
    if (overlay.classList.contains("show")) {
        overlay.classList.remove("show");
        document.removeEventListener("mousemove", onMouseMove);
        document.removeEventListener("click", onClick);
    }
}

function onMouseMove() {
    const overlay = document.getElementById("dialog_overlay");
    if (overlay) {
        removeOverlay(overlay);
    }
}

function onClick() {
    const overlay = document.getElementById("dialog_overlay");
    if (overlay) {
        removeOverlay(overlay);
    }
}

window.addEventListener('DOMContentLoaded', () => {
    const overlay = document.getElementById("dialog_overlay");
    if (overlay && overlay.classList.contains("show")) {
        setTimeout(() => removeOverlay(overlay), 5000);
        document.addEventListener("mousemove", onMouseMove);
        document.addEventListener("click", onClick);
    }
});
