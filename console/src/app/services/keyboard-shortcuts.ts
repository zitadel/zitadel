export interface KeyboardShortcut {
  keyboardKeys: string[];
  link: any[];
  i18nKey: string;
}

export const HOME: KeyboardShortcut = {
  i18nKey: 'KEYBOARDSHORTCUTS.SHORTCUTS.HOME',
  link: ['/'],
  keyboardKeys: ['g', 'h'],
};

export const ORG: KeyboardShortcut = {
  i18nKey: 'KEYBOARDSHORTCUTS.SHORTCUTS.ORG',
  link: ['/org'],
  keyboardKeys: ['g', 'o'],
};

export const PROJECTS: KeyboardShortcut = {
  i18nKey: 'KEYBOARDSHORTCUTS.SHORTCUTS.PROJECTS',
  link: ['/projects'],
  keyboardKeys: ['g', 'p'],
};

export const USERS: KeyboardShortcut = {
  i18nKey: 'KEYBOARDSHORTCUTS.SHORTCUTS.USERS',
  link: ['/'],
  keyboardKeys: ['g', 'u'],
};

export const USERGRANTS: KeyboardShortcut = {
  i18nKey: 'KEYBOARDSHORTCUTS.SHORTCUTS.USERGRANTS',
  link: ['/'],
  keyboardKeys: ['g', 'a'],
};

export const ACTIONS: KeyboardShortcut = {
  i18nKey: 'KEYBOARDSHORTCUTS.SHORTCUTS.ACTIONS',
  link: ['/'],
  keyboardKeys: ['g', 'f'],
};

export const SIDEWIDESHORTCUTS = [HOME, ORG, PROJECTS, USERS, USERGRANTS, ACTIONS];
