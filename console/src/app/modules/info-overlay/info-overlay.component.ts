import { Component } from '@angular/core';

@Component({
  selector: 'cnsl-info-overlay',
  templateUrl: './info-overlay.component.html',
  styleUrls: ['./info-overlay.component.scss'],
})
export class InfoOverlayComponent {
  public show: boolean = false;
  constructor() {}

  public dismiss() {
    this.show = false;
  }
}
