import { computed, Injector, Pipe, PipeTransform } from '@angular/core';
import { delay, Observable } from 'rxjs';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { NewAuthService } from '../../services/new-auth.service';
import { toObservable } from '@angular/core/rxjs-interop';

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
    const signal = computed(() => this.authService.hasRoles(this.permissions.data() ?? [], values, requiresAll));

    return toObservable(signal, {
      injector: this.injector,
    }).pipe(delay(0));
  }
}
