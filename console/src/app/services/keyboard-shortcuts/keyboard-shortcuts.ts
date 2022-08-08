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

export const INSTANCE: KeyboardShortcut = {
  i18nKey: 'KEYBOARDSHORTCUTS.SHORTCUTS.INSTANCE',
  link: ['/instance'],
  keyboardKeys: ['g', 'i'],
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

export const ORGSETTINGS: KeyboardShortcut = {
  i18nKey: 'KEYBOARDSHORTCUTS.SHORTCUTS.ORGSETTINGS',
  link: ['/org-settings'],
  keyboardKeys: ['g', 's'],
  permissions: ['org.read'],
};

export const SIDEWIDESHORTCUTS = [ME, HOME, INSTANCE, ORG, PROJECTS, USERS, USERGRANTS, ACTIONS, DOMAINS, ORGSETTINGS];

export const ORGSWITCHER: KeyboardShortcut = {
  i18nKey: 'KEYBOARDSHORTCUTS.SHORTCUTS.ORGSWITCHER',
  link: ['/org'],
  keyboardKeys: ['/'],
};

export const ORGSHORTCUTS = [ORGSWITCHER];
