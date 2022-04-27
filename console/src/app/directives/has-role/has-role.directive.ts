import { Directive, Input, TemplateRef, ViewContainerRef } from '@angular/core';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';

@Directive({
  selector: '[cnslHasRole]',
})
export class HasRoleDirective {
  private hasView: boolean = false;
  @Input() public set hasRole(roles: string[] | RegExp[] | undefined) {
    if (roles && roles.length > 0) {
      this.authService.isAllowed(roles).subscribe((isAllowed) => {
        if (isAllowed && !this.hasView) {
          this.viewContainerRef.clear();
          this.viewContainerRef.createEmbeddedView(this.templateRef);
        } else {
          this.viewContainerRef.clear();
          this.hasView = false;
        }
      });
    } else {
      if (!this.hasView) {
        this.viewContainerRef.clear();
        this.viewContainerRef.createEmbeddedView(this.templateRef);
      }
    }
  }

  constructor(
    private authService: GrpcAuthService,
    protected templateRef: TemplateRef<any>,
    protected viewContainerRef: ViewContainerRef,
  ) {}
}
