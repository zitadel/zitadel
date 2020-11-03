import { Component, Input } from '@angular/core';
import { PolicyComponentType } from 'src/app/modules/policies/policy-component-types.enum';

export enum PolicyGridType {
    ORG,
    IAM,
}

@Component({
    selector: 'app-policy-grid',
    templateUrl: './policy-grid.component.html',
    styleUrls: ['./policy-grid.component.scss'],
})
export class PolicyGridComponent {
    @Input() public type!: PolicyGridType;
    public PolicyComponentType: any = PolicyComponentType;
    public PolicyGridType: any = PolicyGridType;
    constructor() { }
}
