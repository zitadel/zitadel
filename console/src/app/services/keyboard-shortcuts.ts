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

export const SYSTEM: KeyboardShortcut = {
  i18nKey: 'KEYBOARDSHORTCUTS.SHORTCUTS.SYSTEM',
  link: ['/system'],
  keyboardKeys: ['g', 's'],
};

export const ORG: KeyboardShortcut = {
  i18nKey: 'KEYBOARDSHORTCUTS.SHORTCUTS.ORG',
  link: ['/org'],
  keyboardKeys: ['g', 'o'],
};

export const ME: KeyboardShortcut = {
  i18nKey: 'KEYBOARDSHORTCUTS.SHORTCUTS.ME',
  link: ['/users/me'],
  keyboardKeys: ['m', 'e'],
};

export const PROJECTS: KeyboardShortcut = {
  i18nKey: 'KEYBOARDSHORTCUTS.SHORTCUTS.PROJECTS',
  link: ['/projects'],
  keyboardKeys: ['g', 'p'],
};

export const USERS: KeyboardShortcut = {
  i18nKey: 'KEYBOARDSHORTCUTS.SHORTCUTS.USERS',
  link: ['/users'],
  keyboardKeys: ['g', 'u'],
};

export const USERGRANTS: KeyboardShortcut = {
  i18nKey: 'KEYBOARDSHORTCUTS.SHORTCUTS.USERGRANTS',
  link: ['/grants'],
  keyboardKeys: ['g', 'a'],
};

export const ACTIONS: KeyboardShortcut = {
  i18nKey: 'KEYBOARDSHORTCUTS.SHORTCUTS.ACTIONS',
  link: ['/actions'],
  keyboardKeys: ['g', 'f'],
};

export const DOMAINS: KeyboardShortcut = {
  i18nKey: 'KEYBOARDSHORTCUTS.SHORTCUTS.DOMAINS',
  link: ['/domains'],
  keyboardKeys: ['g', 'd'],
};

export const SIDEWIDESHORTCUTS = [ME, HOME, SYSTEM, ORG, PROJECTS, USERS, USERGRANTS, ACTIONS, DOMAINS];

export const ORGSWITCHER: KeyboardShortcut = {
  i18nKey: 'KEYBOARDSHORTCUTS.SHORTCUTS.ORGSWITCHER',
  link: ['/org'],
  keyboardKeys: ['/'],
};

export const ORGSHORTCUTS = [ORGSWITCHER];
