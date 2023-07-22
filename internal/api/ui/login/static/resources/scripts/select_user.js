document.addEventListener('DOMContentLoaded', function () {
    const container = document.getElementsByClassName('lgn-account-selection')[0];
    container.addEventListener('change', function (event) {
        const t = event.target;
        if (t.classList.contains('lgn-login-as')) {
            const btn = t.closest('.lgn-account-container').getElementsByClassName('not-org-user')[0];
            if (btn) {
                if (t.checked) {
                    const title = btn.getAttribute('title');
                    btn.removeAttribute('title');
                    if (title) {
                        btn.setAttribute('_title', title);
                    }
                    btn.removeAttribute('disabled');
                } else {
                    const title = btn.getAttribute('_title');
                    btn.removeAttribute('_title');
                    if (title) {
                        btn.setAttribute('title', title);
                    }
                    btn.setAttribute('disabled', 'disabled');
                }
            }
        }
    });
});