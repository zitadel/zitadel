import { ChangeDetectionStrategy, Component, Input } from '@angular/core';

@Component({
  selector: 'cnsl-card',
  templateUrl: './card.component.html',
  styleUrls: ['./card.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class CardComponent {
  @Input() public title: string = '';
  @Input() public description: string = '';
  @Input() public expanded: boolean = true;
  @Input() public warn: boolean = false;
  @Input() public nomargin: boolean = false;
}
