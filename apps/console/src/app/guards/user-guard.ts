import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { map, take } from 'rxjs';

import { GrpcAuthService } from '../services/grpc-auth.service';

export const userGuard: CanActivateFn = (route) => {
  const authService = inject(GrpcAuthService);
  const router = inject(Router);

  return authService.user.pipe(
    take(1),
    map((user) => {
      const isMe = user?.id === route.params['id'];
      if (isMe) {
        router.navigate(['/users', 'me']).then();
      }
      return !isMe;
    }),
  );
};
