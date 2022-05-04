import { Component } from '@angular/core';

import { PolicyComponentServiceType } from '../policies/policy-component-types.enum';
import { SidenavSetting } from '../sidenav/sidenav.component';

const GENERAL: SidenavSetting = { id: 'general', i18nKey: 'SETTINGS.LIST.GENERAL' };
const LOGIN: SidenavSetting = { id: 'login', i18nKey: 'SETTINGS.LIST.LOGIN', groupI18nKey: 'SETTINGS.GROUPS.LOGIN' };
const IDP: SidenavSetting = { id: 'idp', i18nKey: 'SETTINGS.LIST.IDP', groupI18nKey: 'SETTINGS.GROUPS.LOGIN' };
const NOTIFICATIONPROVIDERS: SidenavSetting = {
  id: 'notificationproviders',
  i18nKey: 'SETTINGS.LIST.NOTIFICATIONPROVIDERS',
  groupI18nKey: 'SETTINGS.GROUPS.NOTIFICATIONS',
};

const NOTIFICATIONS: SidenavSetting = {
  id: 'notifications',
  i18nKey: 'SETTINGS.LIST.NOTIFICATIONS',
  groupI18nKey: 'SETTINGS.GROUPS.NOTIFICATIONS',
};
const MESSAGETEXTS: SidenavSetting = {
  id: 'messagetexts',
  i18nKey: 'SETTINGS.LIST.MESSAGETEXTS',
  groupI18nKey: 'SETTINGS.GROUPS.APPEARANCE',
};

const LOGINTEXTS: SidenavSetting = {
  id: 'logintexts',
  i18nKey: 'SETTINGS.LIST.LOGINTEXTS',
  groupI18nKey: 'SETTINGS.GROUPS.APPEARANCE',
};
const PRIVACYPOLICY: SidenavSetting = {
  id: 'privacypolicy',
  i18nKey: 'SETTINGS.LIST.PRIVACYPOLICY',
  groupI18nKey: 'SETTINGS.GROUPS.OTHER',
};
const BRANDING: SidenavSetting = {
  id: 'branding',
  i18nKey: 'SETTINGS.LIST.BRANDING',
  groupI18nKey: 'SETTINGS.GROUPS.APPEARANCE',
};

@Component({
  selector: 'cnsl-settings-list',
  templateUrl: './settings-list.component.html',
  styleUrls: ['./settings-list.component.scss'],
})
export class SettingsListComponent {
  public settingsList: SidenavSetting[] = [
    GENERAL,
    LOGIN,
    IDP,
    NOTIFICATIONS,
    NOTIFICATIONPROVIDERS,
    BRANDING,
    MESSAGETEXTS,
    LOGINTEXTS,
    PRIVACYPOLICY,
  ];
  public currentSetting: string | undefined = 'general';
  public PolicyComponentServiceType: any = PolicyComponentServiceType;
  constructor() {}

  private changeSelection(small: boolean): void {
    if (small) {
      this.currentSetting = undefined;
    } else {
      this.currentSetting = this.currentSetting === undefined ? 'general' : this.currentSetting;
    }
  }
}
