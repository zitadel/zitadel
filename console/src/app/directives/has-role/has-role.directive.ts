import { Directive, Input, OnDestroy, TemplateRef, ViewContainerRef } from '@angular/core';
import { Subject, takeUntil } from 'rxjs';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';

@Directive({
  selector: '[cnslHasRole]',
})
export class HasRoleDirective implements OnDestroy {
  private destroy$: Subject<void> = new Subject();
  private hasView: boolean = false;
  @Input() public set hasRole(roles: string[] | RegExp[] | undefined) {
    if (roles && roles.length > 0) {
      this.authService
        .isAllowed(roles)
        .pipe(takeUntil(this.destroy$))
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
  ) {}

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }
}
