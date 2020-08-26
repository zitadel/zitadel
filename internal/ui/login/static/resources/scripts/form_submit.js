function disableSubmit(checks, button) {
    let form = document.getElementsByTagName('form')[0];
    let inputs = form.getElementsByTagName('input');
    for (i = 0; i < inputs.length; i++) {
        button.disabled = true;
        inputs[i].addEventListener('input', function () {
            if (checks != undefined) {
                if (checks() === false) {
                    button.disabled = true;
                    return
                }
            }
            if (checkRequired(form, inputs) === false) {
                button.disabled = true;
                return
            }
            button.disabled = false;
        });
    }
}

function checkRequired(form, inputs) {
    for (i = 0; i < inputs.length; i++) {
        if (inputs[i].required && inputs[i].value == '') {
            return false
        }
    }
    return true;
}