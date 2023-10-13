import { Injectable } from '@angular/core';
import { ActivatedRouteSnapshot, Router, RouterStateSnapshot } from '@angular/router';
import { map, Observable, take } from 'rxjs';

import { GrpcAuthService } from '../services/grpc-auth.service';

@Injectable({
  providedIn: 'root',
})
export class UserGuard {
  constructor(
    private authService: GrpcAuthService,
    private router: Router,
  ) {}

  public canActivate(
    route: ActivatedRouteSnapshot,
    state: RouterStateSnapshot,
  ): Observable<boolean> | Promise<boolean> | boolean {
    return this.authService.user.pipe(
      take(1),
      map((user) => {
        const isMe = user?.id === route.params['id'];
        if (isMe) {
          this.router.navigate(['/users', 'me']);
        }
        return !isMe;
      }),
    );
  }
}
