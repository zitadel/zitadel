export interface KeyboardShortcut {
  keyboardKeys: string[];
  link: any[];
  i18nKey: string;
  permissions?: string[] | RegExp[];
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
  permissions: ['iam.read'],
};

export const ORG: KeyboardShortcut = {
  i18nKey: 'KEYBOARDSHORTCUTS.SHORTCUTS.ORG',
  link: ['/org'],
  keyboardKeys: ['g', 'o'],
  permissions: ['org.read'],
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
  permissions: ['project.read'],
};

export const USERS: KeyboardShortcut = {
  i18nKey: 'KEYBOARDSHORTCUTS.SHORTCUTS.USERS',
  link: ['/users'],
  keyboardKeys: ['g', 'u'],
  permissions: ['user.read'],
};

export const USERGRANTS: KeyboardShortcut = {
  i18nKey: 'KEYBOARDSHORTCUTS.SHORTCUTS.USERGRANTS',
  link: ['/grants'],
  keyboardKeys: ['g', 'a'],
  permissions: ['usergrant.read'],
};

export const ACTIONS: KeyboardShortcut = {
  i18nKey: 'KEYBOARDSHORTCUTS.SHORTCUTS.ACTIONS',
  link: ['/actions'],
  keyboardKeys: ['g', 'f'],
  permissions: ['org.action.read'],
};

export const DOMAINS: KeyboardShortcut = {
  i18nKey: 'KEYBOARDSHORTCUTS.SHORTCUTS.DOMAINS',
  link: ['/domains'],
  keyboardKeys: ['g', 'd'],
  permissions: ['org.read'],
};

export const SIDEWIDESHORTCUTS = [ME, HOME, SYSTEM, ORG, PROJECTS, USERS, USERGRANTS, ACTIONS, DOMAINS];

export const ORGSWITCHER: KeyboardShortcut = {
  i18nKey: 'KEYBOARDSHORTCUTS.SHORTCUTS.ORGSWITCHER',
  link: ['/org'],
  keyboardKeys: ['/'],
};

export const ORGSHORTCUTS = [ORGSWITCHER];
