import { Component, Input, OnInit } from '@angular/core';

@Component({
  selector: 'cnsl-top-view',
  templateUrl: './top-view.component.html',
  styleUrls: ['./top-view.component.scss'],
})
export class TopViewComponent implements OnInit {
  @Input() public title: string = '';
  @Input() public sub: string = '';
  @Input() public stateTooltip: string = '';
  @Input() public isActive: boolean = false;
  @Input() public isInactive: boolean = false;
  @Input() public backRouterLink!: any[];
  @Input() public docLink: string = '';

  constructor() {}

  ngOnInit(): void {}
}
