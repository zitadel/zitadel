import { Component } from '@angular/core';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';

@Component({
  selector: 'cnsl-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss'],
})
export class HomeComponent {
  public dark: boolean = true;

  constructor(public authService: GrpcAuthService) {
    const theme = localStorage.getItem('theme');
    this.dark = theme === 'dark-theme' ? true : theme === 'light-theme' ? false : true;
  }
}
