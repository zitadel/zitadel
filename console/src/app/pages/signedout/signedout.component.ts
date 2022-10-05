import { Component } from '@angular/core';
import { ThemeService } from 'src/app/services/theme.service';

@Component({
  selector: 'cnsl-signedout',
  templateUrl: './signedout.component.html',
  styleUrls: ['./signedout.component.scss'],
})
export class SignedoutComponent {
  public dark: boolean = true;

  constructor(public themeService: ThemeService) {
    const theme = localStorage.getItem('theme');
    this.dark = theme === 'dark-theme' ? true : theme === 'light-theme' ? false : true;
  }
}
