import { animate, style, transition, trigger } from '@angular/animations';
import { Component, Input } from '@angular/core';
import { PolicyComponentServiceType, PolicyComponentType } from 'src/app/modules/policies/policy-component-types.enum';

import { SETTINGLINKS, SettingLinks } from './policies';

@Component({
  selector: 'cnsl-policy-grid',
  templateUrl: './policy-grid.component.html',
  styleUrls: ['./policy-grid.component.scss'],
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
export class PolicyGridComponent {
  @Input() public type!: PolicyComponentServiceType;
  @Input() public tag: string = '';
  public PolicyComponentType: any = PolicyComponentType;
  public PolicyComponentServiceType: any = PolicyComponentServiceType;
  public SETTINGS: SettingLinks[] = SETTINGLINKS;
}
