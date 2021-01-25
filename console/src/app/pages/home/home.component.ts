import { Component } from '@angular/core';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';

@Component({
    selector: 'app-home',
    templateUrl: './home.component.html',
    styleUrls: ['./home.component.scss'],
})
export class HomeComponent {
    public dark: boolean = true;
    public firstStepsDismissed: boolean = false;
    constructor(public authService: GrpcAuthService) {
        const theme = localStorage.getItem('theme');
        this.dark = theme === 'dark-theme' ? true : theme === 'light-theme' ? false : true;

        this.firstStepsDismissed = localStorage.getItem('firstStartDismissed') == 'true' ? true : false;
    }

    dismissFirstSteps(event: Event): void {
        event.preventDefault();
        localStorage.setItem('firstStartDismissed', 'true');
        this.firstStepsDismissed = true;
    }
}
