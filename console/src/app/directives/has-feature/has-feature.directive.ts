import { Directive, Input, TemplateRef, ViewContainerRef } from '@angular/core';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';

@Directive({
  selector: '[cnslHasFeature]',
})
export class HasFeatureDirective {
  private hasView: boolean = false;
  @Input() public set hasFeature(features: string[] | RegExp[] | undefined) {
    if (features && features.length > 0) {
      this.authService.canUseFeature(features).subscribe((isAllowed) => {
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
