import { Component, OnInit } from '@angular/core';
import { PrivacyPolicy } from 'src/app/proto/generated/zitadel/policy_pb';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { faXTwitter } from '@fortawesome/free-brands-svg-icons';

@Component({
  selector: 'cnsl-footer',
  templateUrl: './footer.component.html',
  styleUrls: ['./footer.component.scss'],
})
export class FooterComponent {
  public faXTwitter = faXTwitter;
  constructor(public authService: GrpcAuthService) {}
}
