import { Injectable } from '@angular/core';
import { EMPTY, Subject, catchError, merge, shareReplay } from 'rxjs';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { User } from 'src/app/proto/generated/zitadel/user_pb';

/** Keep the sidebar in sync after native profile/avatar changes without altering the shared service. */
@Injectable({ providedIn: 'root' })
export class AccountIdentity {
  private updated = new Subject<User.AsObject>();
  readonly user;
  constructor(native: GrpcAuthService) {
    this.user = merge(native.user.pipe(catchError(() => EMPTY)), this.updated).pipe(
      shareReplay({ refCount: true, bufferSize: 1 }),
    );
  }
  update(user: User.AsObject): void {
    this.updated.next(user);
  }
}
