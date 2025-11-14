import { ChangeDetectionStrategy, Component, Input } from '@angular/core';
import { NgIconComponent, provideIcons } from '@ng-icons/core';
import { heroChevronUpDown } from '@ng-icons/heroicons/outline';
import { MatRippleModule } from '@angular/material/core';

@Component({
  selector: 'cnsl-header-button',
  templateUrl: './header-button.component.html',
  styleUrls: ['./header-button.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [NgIconComponent, MatRippleModule],
  providers: [provideIcons({ heroChevronUpDown })],
})
export class HeaderButtonComponent {
  @Input() ariaLabel: string = '';
}
