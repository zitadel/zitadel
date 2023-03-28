import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';

import { PolicyComponentServiceType } from '../policies/policy-component-types.enum';
import { SidenavSetting } from '../sidenav/sidenav.component';

@Component({
  selector: 'cnsl-settings-list',
  templateUrl: './settings-list.component.html',
  styleUrls: ['./settings-list.component.scss'],
})
export class SettingsListComponent implements OnChanges {
  @Input() public title: string = '';
  @Input() public description: string = '';
  @Input() public serviceType!: PolicyComponentServiceType;
  @Input() public selectedId: string = '';
  @Input() public settingsList: SidenavSetting[] = [];
  public currentSetting: string | undefined = '';
  public PolicyComponentServiceType: any = PolicyComponentServiceType;
  constructor() {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes['selectedId']?.currentValue) {
      this.currentSetting =
        this.settingsList && this.settingsList.find((l) => l.id === changes['selectedId'].currentValue)
          ? changes['selectedId'].currentValue
          : '';
    } else {
      this.currentSetting = this.settingsList ? this.settingsList[0].id : '';
    }
  }
}
