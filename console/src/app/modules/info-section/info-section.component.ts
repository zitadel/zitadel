import { Component, Input } from '@angular/core';

export enum InfoSectionType {
  INFO = 'INFO',
  SUCCESS = 'SUCCESS',
  WARN = 'WARN',
  ALERT = 'ALERT',
}

@Component({
  selector: 'cnsl-info-section',
  templateUrl: './info-section.component.html',
  styleUrls: ['./info-section.component.scss'],
})
export class InfoSectionComponent {
  @Input() type: InfoSectionType = InfoSectionType.INFO;
  @Input() fitWidth: boolean = false;
}
