const usesDarkTheme = window.matchMedia('(prefers-color-scheme: dark)').matches;
if (usesDarkTheme) {
    document.documentElement.classList.replace('lgn-light-theme', 'lgn-dark-theme');
} else {
    document.documentElement.classList.replace('lgn-dark-theme', 'lgn-light-theme');
}
