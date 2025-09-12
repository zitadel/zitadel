import { Component, effect, Input, OnInit, signal } from '@angular/core';

import { PolicyComponentServiceType } from '../policies/policy-component-types.enum';
import { SidenavSetting } from '../sidenav/sidenav.component';

@Component({
  selector: 'cnsl-settings-list',
  templateUrl: './settings-list.component.html',
  styleUrls: ['./settings-list.component.scss'],
})
export class SettingsListComponent implements OnInit {
  @Input({ required: true }) public serviceType!: PolicyComponentServiceType;
  @Input() public set selectedId(selectedId: string) {
    this.selectedId$.set(selectedId);
  }
  @Input({ required: true }) public settingsList: SidenavSetting[] = [];

  protected setting = signal<SidenavSetting | null>(null);
  private selectedId$ = signal<string | undefined>(undefined);
  protected PolicyComponentServiceType: any = PolicyComponentServiceType;

  constructor() {
    effect(
      () => {
        const selectedId = this.selectedId$();
        if (!selectedId) {
          return;
        }

        const setting = this.settingsList.find(({ id }) => id === selectedId);
        if (!setting) {
          return;
        }
        this.setting.set(setting);
      },
      { allowSignalWrites: true },
    );
  }

  ngOnInit(): void {
    const firstSetting = this.settingsList[0];
    if (!firstSetting || this.setting()) {
      return;
    }
    this.setting.set(firstSetting);
  }
}
