import { Injectable } from '@angular/core';
import { ActivatedRouteSnapshot, CanActivate, RouterStateSnapshot } from '@angular/router';
import { Observable } from 'rxjs';

import { AuthUserService } from '../services/auth-user.service';

@Injectable({
    providedIn: 'root',
})
export class RoleGuard implements CanActivate {

    constructor(private userService: AuthUserService) { }

    public canActivate(
        route: ActivatedRouteSnapshot,
        state: RouterStateSnapshot,
    ): Observable<boolean> {
        return this.userService.isAllowed(route.data['roles']);
    }
}
