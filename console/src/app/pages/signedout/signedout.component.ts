import { Component } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { take } from 'rxjs';
import { LabelPolicy } from 'src/app/proto/generated/zitadel/policy_pb';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ThemeService } from 'src/app/services/theme.service';

@Component({
  selector: 'cnsl-signedout',
  templateUrl: './signedout.component.html',
  styleUrls: ['./signedout.component.scss'],
})
export class SignedoutComponent {
  public dark: boolean = true;

  public labelpolicy?: LabelPolicy.AsObject;
  public queryParams = { state: '' };
  constructor(route: ActivatedRoute, themeService: ThemeService, authService: GrpcAuthService) {
    const theme = localStorage.getItem('theme');
    this.dark = theme === 'dark-theme' ? true : theme === 'light-theme' ? false : true;

    route.queryParams.pipe(take(1)).subscribe((queryParams) => {
      const state = queryParams.state;
      if (state) {
        const parsed = JSON.parse(state);
        if (parsed) {
          this.labelpolicy = parsed;
          themeService.applyLabelPolicy(parsed);
          authService.labelpolicy.next(parsed);
        }
      }
    });
  }
}
