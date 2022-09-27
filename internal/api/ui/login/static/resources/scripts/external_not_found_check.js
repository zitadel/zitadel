let buttons1 = document.getElementsByName("linkbutton");

for (let i = 0; i < buttons1.length; i++) {
    disableSubmit(undefined, buttons1[i]);
}

let buttons2 = document.getElementsByName("autoregisterbutton");
for (let i = 0; i < buttons2.length; i++) {
    disableSubmit(undefined, buttons2[i]);
}
