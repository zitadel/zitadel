import { Component, EventEmitter, Input, Output } from '@angular/core';

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
  @Input() public hasActions: boolean | null = false;
  @Input() public hasContributors: boolean | null = false;
  @Input() public docLink: string = '';
  @Input() public hasBackButton: boolean | null = true;
  @Output() public backClicked: EventEmitter<void> = new EventEmitter();

  constructor() {}

  public backClick(): void {
    this.backClicked.emit();
  }
}
