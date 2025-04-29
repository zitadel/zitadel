import { Component } from '@angular/core';
import { LabelPolicy } from 'src/app/proto/generated/zitadel/policy_pb';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
const LABELPOLICY_LOCALSTORAGE_KEY = 'labelPolicyOnSignout';

@Component({
  selector: 'cnsl-signedout',
  templateUrl: './signedout.component.html',
  styleUrls: ['./signedout.component.scss'],
})
export class SignedoutComponent {
  public dark: boolean = true;

  public labelpolicy?: LabelPolicy.AsObject;
  public queryParams = { state: '' };
  constructor(authService: GrpcAuthService) {
    const theme = localStorage.getItem('theme');
    this.dark = theme === 'dark-theme' ? true : theme === 'light-theme' ? false : true;

    const lP = localStorage.getItem(LABELPOLICY_LOCALSTORAGE_KEY);

    if (!lP) {
      authService.labelPolicyLoading$.next(false);
      return;
    }

    const parsed = JSON.parse(lP);
    localStorage.removeItem(LABELPOLICY_LOCALSTORAGE_KEY);

    if (!parsed) {
      return;
    }

    this.labelpolicy = parsed;
    // todo: figure this one out
    // authService.labelpolicy.next(parsed);
    authService.labelPolicyLoading$.next(false);
  }
}
