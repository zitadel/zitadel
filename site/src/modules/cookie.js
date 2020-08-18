export function getCookie(name, cookies) {
    if (cookies == null) {
        if (typeof window === 'undefined') {
            return undefined;
        }
        cookies = document.cookie;
    }

    const kv = cookies.split(';').find((part) => part.trim().startsWith(name));

    if (!kv) return undefined;

    const cookieValue = kv.split('=')[1];
    if (!cookieValue) return undefined;

    return decodeURIComponent(cookieValue.trim());
}

export function setCookie(name, value, options = {}) {
    if (options.expires instanceof Date) {
        options.expires = options.expires.toUTCString();
    }

    let updatedCookie = {
        [encodeURIComponent(name)]: encodeURIComponent(value),
        sameSite: 'strict',
        ...options,
    };

    document.cookie = Object.entries(updatedCookie)
        .map((kv) => kv.join('='))
        .join(';');
}