import { Component, Input } from '@angular/core';

@Component({
  selector: 'cnsl-top-view',
  templateUrl: './top-view.component.html',
  styleUrls: ['./top-view.component.scss'],
})
export class TopViewComponent {
  @Input() public title: string = '';
  @Input() public sub: string = '';
  @Input() public stateTooltip: string = '';
  @Input() public isActive: boolean = false;
  @Input() public isInactive: boolean = false;
  @Input() public hasActions: boolean = false;
  @Input() public backRouterLink!: any[];
  @Input() public backQueryParams!: any;
  @Input() public docLink: string = '';

  constructor() {}
}
