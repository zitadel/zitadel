import { inject } from '@angular/core';
import { CanActivateFn } from '@angular/router';

import { GrpcAuthService } from '../services/grpc-auth.service';

export const roleGuard: CanActivateFn = (route) => {
  const authService = inject(GrpcAuthService);
  return authService.isAllowed(route.data['roles'], route.data['requiresAll']);
};
