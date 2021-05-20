const avatars = document.getElementsByClassName('lgn-avatar');
for (let i = 0; i < avatars.length; i++) {
    const displayName = avatars[i].getAttribute('loginname');
    if (displayName) {
        const username = displayName.split('@')[0];
        let separator = '_';
        if (username.includes('-')) {
            separator = '-';
        }
        if (username.includes('.')) {
            separator = '.';
        }
        const split = username.split(separator);
        const initials = split[0].charAt(0) + (split[1] ? split[1].charAt(0) : '');
        avatars[i].getElementsByClassName('initials')[0].innerHTML = initials;

        avatars[i].style.background = this.getColor(displayName);
        // set default white text instead of contrast text mode
        avatars[i].style.color = '#ffffff';
    }
}

function getColor(username) {
    const s = 40;
    const l = 50;
    const l2 = 62.5;
    let hash = 0;
    for (let i = 0; i < username.length; i++) {
        hash = username.charCodeAt(i) + ((hash << 5) - hash);
    }

    const h = hash % 360;
    const col1 = 'hsl(' + h + ', ' + s + '%, ' + l + '%)';
    const col2 = 'hsl(' + h + ', ' + s + '%, ' + l2 + '%)';
    return 'linear-gradient(40deg, ' + col1 + ' 30%, ' + col2 + ')';
}