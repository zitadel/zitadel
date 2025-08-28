import { computed, Injector, Pipe, PipeTransform } from '@angular/core';
import { delay, Observable } from 'rxjs';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { NewAuthService } from '../../services/new-auth.service';
import { toObservable } from '@angular/core/rxjs-interop';
import { filter } from 'rxjs/operators';

@Pipe({
  name: 'hasRole',
})
export class HasRolePipe implements PipeTransform {
  private readonly permissions = this.newAuthService.listMyZitadelPermissionsQuery();

  constructor(
    private readonly authService: GrpcAuthService,
    private readonly newAuthService: NewAuthService,
    private readonly injector: Injector,
  ) {}

  public transform(values: string[], requiresAll: boolean = false): Observable<boolean> {
    const signal = computed(() => {
      const permissions = this.permissions.data();
      if (!permissions) {
        return undefined;
      }
      return this.authService.hasRoles(permissions, values, requiresAll);
    });

    return toObservable(signal, {
      injector: this.injector,
    }).pipe(
      filter((hasRole): hasRole is Exclude<typeof hasRole, undefined> => hasRole !== undefined),
      delay(0),
    );
  }
}
