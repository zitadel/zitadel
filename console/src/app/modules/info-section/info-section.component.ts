import { ChangeDetectionStrategy, Component, Input } from '@angular/core';

export enum InfoSectionType {
  INFO = 'INFO',
  WARN = 'WARN',
  ALERT = 'ALERT',
}

@Component({
  selector: 'cnsl-info-section',
  templateUrl: './info-section.component.html',
  styleUrls: ['./info-section.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class InfoSectionComponent {
  @Input() type: InfoSectionType = InfoSectionType.INFO;
  @Input() fitWidth: boolean = false;

  protected readonly infoSectionType = InfoSectionType;
}
