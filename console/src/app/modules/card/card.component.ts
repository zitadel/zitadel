import { animate, style, transition, trigger } from '@angular/animations';
import { Component, Input } from '@angular/core';

@Component({
  selector: 'cnsl-card',
  templateUrl: './card.component.html',
  styleUrls: ['./card.component.scss'],
  animations: [
    trigger('openClose', [
      transition(':enter', [
        style({ height: '0', opacity: 0 }),
        animate('150ms ease-in-out', style({ height: '*', opacity: 1 })),
      ]),
      transition(':leave', [animate('150ms ease-in-out', style({ height: '0', opacity: 0 }))]),
    ]),
  ],
})
export class CardComponent {
  @Input() public expanded: boolean = true;
  @Input() public warn: boolean = false;
  @Input() public title: string = '';
  @Input() public description: string = '';
  @Input() public animate: boolean = false;
  @Input() public nomargin?: boolean = false;
  @Input() public stretch: boolean = false;
}
