import { Pipe, PipeTransform } from '@angular/core';
import { Observable } from 'rxjs';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';

@Pipe({
  name: 'hasRole',
})
export class HasRolePipe implements PipeTransform {
  constructor(private authService: GrpcAuthService) {}

  public transform(values: string[], requresAll: boolean = false): Observable<boolean> {
    return this.authService.isAllowed(values, requresAll);
  }
}
