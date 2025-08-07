import { Component, EventEmitter, Input, Output } from '@angular/core';

@Component({
  selector: 'cnsl-create-layout',
  templateUrl: './create-layout.component.html',
  styleUrls: ['./create-layout.component.scss'],
})
export class CreateLayoutComponent {
  @Input() currentCreateStep: number = 1;
  @Input() createSteps: number = 1;
  @Input() title: string = '';
  @Output() closed: EventEmitter<void> = new EventEmitter();
  constructor() {}

  close() {
    this.closed.emit();
  }
}
