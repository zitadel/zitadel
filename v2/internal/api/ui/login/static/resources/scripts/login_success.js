document.addEventListener('DOMContentLoaded', function () {
    autoSubmit();
});

function autoSubmit() {
    let form = document.getElementsByTagName('form')[0];
    if (form) {
        let button = document.getElementById("redirect-button");
        if (button) {
            button.addEventListener("click", function (event) {
                location.reload();
                event.preventDefault();
            });
        }
        form.submit();
    }
}
