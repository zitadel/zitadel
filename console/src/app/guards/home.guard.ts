import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { map } from 'rxjs/operators';

import { GrpcAuthService } from '../services/grpc-auth.service';

export const homeGuard: CanActivateFn = (route) => {
  const authService = inject(GrpcAuthService);
  const router = inject(Router);

  // Check if user has any roles (using the same logic as roleGuard)
  return authService.isAllowed(route.data['roles'], route.data['requiresAll']).pipe(
    map((hasRoles) => {
      if (!hasRoles) {
        // User has no roles, redirect to /users/me
        router.navigate(['/users/me']);
        return false;
      }
      return true;
    }),
  );
};
