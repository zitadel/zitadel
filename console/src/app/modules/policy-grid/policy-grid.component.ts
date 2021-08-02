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
  public tags: Set<string> = new Set(POLICIES.map(p => p.tags).flat());

  @Input() public tagForFilter: string = '';
  @Input() public currentPolicy!: GridPolicy;

  public get filteredPolicies(): GridPolicy[] {
    if (this.tagForFilter) {
      return POLICIES.filter(p => p !== this.currentPolicy && p.tags.includes(this.tagForFilter));
    } else {
      return POLICIES.filter(p => p !== this.currentPolicy);
    }
  }
}
