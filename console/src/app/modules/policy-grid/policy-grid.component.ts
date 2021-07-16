import { Component, Input } from '@angular/core';
import { PolicyComponentServiceType, PolicyComponentType } from 'src/app/modules/policies/policy-component-types.enum';

import { GridPolicy, POLICIES } from './policies';

@Component({
  selector: 'app-policy-grid',
  templateUrl: './policy-grid.component.html',
  styleUrls: ['./policy-grid.component.scss'],
})
export class PolicyGridComponent {
  @Input() public type!: PolicyComponentServiceType;
  @Input() public tag: string = '';
  public PolicyComponentType: any = PolicyComponentType;
  public PolicyComponentServiceType: any = PolicyComponentServiceType;
  public POLICIES: GridPolicy[] = POLICIES;
}
