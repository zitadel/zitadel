document.addEventListener('DOMContentLoaded', function () {
    autoSubmit();
});

function autoSubmit() {
    let form = document.getElementsByTagName('form')[0];
    if (form) {
        form.submit();
    }
}
