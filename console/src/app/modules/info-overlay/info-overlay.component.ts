import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'cnsl-info-overlay',
  templateUrl: './info-overlay.component.html',
  styleUrls: ['./info-overlay.component.scss'],
})
export class InfoOverlayComponent implements OnInit {
  public show: boolean = false;
  constructor() {}

  ngOnInit(): void {}

  public dismiss() {
    this.show = false;
  }
}
