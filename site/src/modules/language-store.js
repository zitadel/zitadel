import { writable } from 'svelte/store';

export const docLanguages = writable(['de', 'en']);

export function storeValue(lngs) {
    console.log('lngs: ' + lngs);
    docLanguages.set(lngs);
}