import { Pipe, PipeTransform } from '@angular/core';
import { delay, map, Observable } from 'rxjs';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { NewAuthService } from '../../services/new-auth.service';
import { toObservable } from '@angular/core/rxjs-interop';
import { filter } from 'rxjs/operators';

@Pipe({
  name: 'hasRole',
  standalone: false,
})
export class HasRolePipe implements PipeTransform {
  private readonly permissions$ = toObservable(this.newAuthService.listMyZitadelPermissionsQuery().data);

  constructor(
    private readonly authService: GrpcAuthService,
    private readonly newAuthService: NewAuthService,
  ) {}

  public transform(values: string[], requiresAll: boolean = false): Observable<boolean> {
    return this.permissions$.pipe(
      map((permissions) => {
        if (!permissions) {
          return undefined;
        }
        return this.authService.hasRoles(permissions, values, requiresAll);
      }),
      filter((hasRole): hasRole is Exclude<typeof hasRole, undefined> => hasRole !== undefined),
      delay(0),
    );
  }
}
