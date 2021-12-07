import { Component } from '@angular/core';

@Component({
  selector: 'cnsl-signedout',
  templateUrl: './signedout.component.html',
  styleUrls: ['./signedout.component.scss'],
})
export class SignedoutComponent {
  public dark: boolean = true;

  constructor() {
    const theme = localStorage.getItem('theme');
    this.dark = theme === 'dark-theme' ? true : theme === 'light-theme' ? false : true;
  }
}
