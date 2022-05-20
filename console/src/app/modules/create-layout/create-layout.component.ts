import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';

@Component({
  selector: 'cnsl-create-layout',
  templateUrl: './create-layout.component.html',
  styleUrls: ['./create-layout.component.scss'],
})
export class CreateLayoutComponent implements OnInit {
  @Input() currentCreateStep: number = 1;
  @Input() createSteps: number = 1;
  @Input() title: string = '';
  @Output() closed: EventEmitter<void> = new EventEmitter();
  constructor() {}

  ngOnInit(): void {}

  close() {
    this.closed.emit();
  }
}
