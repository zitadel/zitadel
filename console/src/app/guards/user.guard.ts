import { Injectable } from '@angular/core';
import { ActivatedRouteSnapshot, CanActivate, Router, RouterStateSnapshot } from '@angular/router';
import { Observable } from 'rxjs';
import { map, tap } from 'rxjs/operators';

import { GrpcAuthService } from '../services/grpc-auth.service';


@Injectable({
    providedIn: 'root',
})
export class UserGuard implements CanActivate {
    constructor(private authService: GrpcAuthService, private router: Router) { }

    public canActivate(
        route: ActivatedRouteSnapshot,
        state: RouterStateSnapshot,
    ): Observable<boolean> | Promise<boolean> | boolean {
        return this.authService.user.pipe(
            map(user => user?.id !== route.params.id),
            tap((isNotMe) => {
                if (!isNotMe) {
                    this.router.navigate(['/users', 'me']);
                }
            }),
        );
    }
}
