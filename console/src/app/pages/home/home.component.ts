import { Component } from '@angular/core';
import { AuthenticationService } from 'src/app/services/authentication.service';

@Component({
    selector: 'app-home',
    templateUrl: './home.component.html',
    styleUrls: ['./home.component.scss'],
})
export class HomeComponent {
    public dark: boolean = true;
    constructor(public authService: AuthenticationService) {
        const theme = localStorage.getItem('theme');
        this.dark = theme === 'dark-theme' ? true : theme === 'light-theme' ? false : true;

    }
}
