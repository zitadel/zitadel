import { animate, style, transition, trigger } from '@angular/animations';
import { Component, Input, OnInit } from '@angular/core';
import { PolicyComponentServiceType } from 'src/app/modules/policies/policy-component-types.enum';

import { SETTINGLINKS, SettingLinks } from './settinglinks';

@Component({
  selector: 'cnsl-settings-grid',
  templateUrl: './settings-grid.component.html',
  styleUrls: ['./settings-grid.component.scss'],
  animations: [
    trigger('policy', [
      transition(':enter', [
        style({
          opacity: 0.5,
        }),
        animate(
          '.15s ease-in-out',
          style({
            opacity: 1,
          }),
        ),
      ]),
      transition(':leave', [
        style({
          opacity: 1,
        }),
        animate(
          '.15s ease-in-out',
          style({
            opacity: 0.5,
          }),
        ),
      ]),
    ]),
  ],
})
export class SettingsGridComponent implements OnInit {
  @Input() public type!: PolicyComponentServiceType;
  @Input() public tag: string = '';
  public PolicyComponentServiceType: any = PolicyComponentServiceType;
  public SETTINGS: SettingLinks[] = SETTINGLINKS;

  ngOnInit(): void {
    this.SETTINGS = this.SETTINGS.filter((setting) =>
      this.type === PolicyComponentServiceType.MGMT ? !!setting.orgRouterLink : !!setting.iamRouterLink,
    );
  }
}
