import { inject, Injectable } from '@angular/core';
import { GrpcService } from './grpc.service';

@Injectable({
  providedIn: 'root',
})
export class AuthorizationService {
  private readonly grpcService = inject(GrpcService);

  public deleteAuthorization(id: string) {
    return this.grpcService.authorization.deleteAuthorization({ id });
  }
}
