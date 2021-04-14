import { Injectable } from '@angular/core';
import { ActivatedRouteSnapshot, CanActivate, RouterStateSnapshot } from '@angular/router';
import { Observable } from 'rxjs';
import { filter, switchMap } from 'rxjs/operators';

import { GrpcAuthService } from '../services/grpc-auth.service';

@Injectable({
    providedIn: 'root',
})
export class RoleGuard implements CanActivate {

    constructor(private authService: GrpcAuthService) { }

    public canActivate(
        route: ActivatedRouteSnapshot,
        state: RouterStateSnapshot,
    ): Observable<boolean> {
        return this.authService.fetchedZitadelPermissions.pipe(
            filter((permissionsFetched) => !!permissionsFetched),
        ).pipe(
            switchMap(_ => this.authService.isAllowed(route.data['roles'])),
        );
    }
}
