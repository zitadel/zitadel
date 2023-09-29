import { Injectable } from '@angular/core';
import { ActivatedRouteSnapshot, Router, RouterStateSnapshot } from '@angular/router';
import { Observable, from, of } from 'rxjs';
import { map, switchMap, take, tap } from 'rxjs/operators';

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
    const user = this.authService.userSubject.getValue();
    const isMe = user?.id === route.params['id']
    if (isMe) {
      this.router.navigate(['/users', 'me']);

    }
    return !isMe;
  }
}
