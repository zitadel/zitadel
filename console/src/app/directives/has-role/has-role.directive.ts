import { DestroyRef, Directive, Input, TemplateRef, ViewContainerRef } from '@angular/core';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';

@Directive({
  selector: '[cnslHasRole]',
})
export class HasRoleDirective {
  private hasView: boolean = false;
  @Input() public set hasRole(roles: string[] | RegExp[] | undefined) {
    if (roles && roles.length > 0) {
      this.authService
        .isAllowed(roles)
        .pipe(takeUntilDestroyed(this.destroyRef))
        .subscribe((isAllowed) => {
          if (isAllowed && !this.hasView) {
            if (this.viewContainerRef.length !== 0) {
              this.viewContainerRef.clear();
            }
            this.viewContainerRef.createEmbeddedView(this.templateRef);
          } else {
            this.viewContainerRef.clear();
            this.hasView = false;
          }
        });
    } else {
      if (!this.hasView) {
        if (this.viewContainerRef.length !== 0) {
          this.viewContainerRef.clear();
        }
        this.viewContainerRef.createEmbeddedView(this.templateRef);
      }
    }
  }

  constructor(
    private authService: GrpcAuthService,
    protected templateRef: TemplateRef<any>,
    protected viewContainerRef: ViewContainerRef,
    private readonly destroyRef: DestroyRef,
  ) {}
}
