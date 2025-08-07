let form = document.getElementsByTagName('form')[0];
let editButton = document.getElementById('edit');
editButton.addEventListener('click', function () {
    form.submit();
});

