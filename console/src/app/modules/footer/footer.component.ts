import { Component } from '@angular/core';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { faXTwitter } from '@fortawesome/free-brands-svg-icons';

@Component({
  selector: 'cnsl-footer',
  templateUrl: './footer.component.html',
  styleUrls: ['./footer.component.scss'],
  standalone: false,
})
export class FooterComponent {
  public faXTwitter = faXTwitter;
  constructor(public authService: GrpcAuthService) {}
}
