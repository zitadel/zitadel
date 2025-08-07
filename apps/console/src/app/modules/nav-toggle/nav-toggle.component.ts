import { Component, EventEmitter, Input, Output } from '@angular/core';

@Component({
  selector: 'cnsl-nav-toggle',
  templateUrl: './nav-toggle.component.html',
  styleUrls: ['./nav-toggle.component.scss'],
})
export class NavToggleComponent {
  @Input() public label: string = '';
  @Input() public count: number | null = 0;
  @Input() public active: boolean = false;
  @Output() public clicked: EventEmitter<void> = new EventEmitter<void>();
  constructor() {}
}
