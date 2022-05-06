import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';

import { PolicyComponentServiceType } from '../policies/policy-component-types.enum';
import { SidenavSetting } from '../sidenav/sidenav.component';
import {
    BRANDING,
    GENERAL,
    IDP,
    LOGIN,
    LOGINTEXTS,
    MESSAGETEXTS,
    NOTIFICATIONPROVIDERS,
    NOTIFICATIONS,
    PRIVACYPOLICY,
} from './settings';

@Component({
  selector: 'cnsl-settings-list',
  templateUrl: './settings-list.component.html',
  styleUrls: ['./settings-list.component.scss'],
})
export class SettingsListComponent implements OnChanges {
  @Input() public title: string = '';
  @Input() public description: string = '';
  @Input() public serviceType!: PolicyComponentServiceType;
  @Input() public selectedId: string = 'general';
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

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.selectedId.currentValue) {
      this.currentSetting = changes.selectedId.currentValue;
    } else {
      this.currentSetting = 'general';
    }
  }

  private changeSelection(small: boolean): void {
    if (small) {
      this.currentSetting = undefined;
    } else {
      this.currentSetting = this.currentSetting === undefined ? 'general' : this.currentSetting;
    }
  }
}
