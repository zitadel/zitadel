import { ChangeDetectionStrategy, Component } from '@angular/core';

@Component({
  selector: 'cnsl-header-button',
  templateUrl: './header-button.component.html',
  styleUrls: ['./header-button.component.scss'],
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class HeaderButtonComponent {}
