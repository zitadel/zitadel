import { Injectable } from '@angular/core';
import { ActivatedRouteSnapshot, RouterStateSnapshot } from '@angular/router';
import { Observable } from 'rxjs';

import { GrpcAuthService } from '../services/grpc-auth.service';

@Injectable({
  providedIn: 'root',
})
export class RoleGuard {
  constructor(private authService: GrpcAuthService) {}

  public canActivate(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<boolean> {
    return this.authService.isAllowed(route.data['roles'], route.data['requiresAll']);
  }
}
