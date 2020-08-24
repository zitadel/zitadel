import { Directive, Input, TemplateRef, ViewContainerRef } from '@angular/core';
import { AuthenticationService } from 'src/app/services/authentication.service';


@Directive({
    selector: '[appHasRole]',
})

export class HasRoleDirective {
    private hasView: boolean = false;
    @Input() public set appHasRole(roles: string[] | RegExp[]) {
        if (roles && roles.length > 0) {
            this.authService.isAllowed(roles).subscribe(isAllowed => {
                if (isAllowed && !this.hasView) {
                    this.viewContainerRef.clear();
                    this.viewContainerRef.createEmbeddedView(this.templateRef);
                } else if (this.hasView) {
                    this.viewContainerRef.clear();
                    this.hasView = false;
                }
            });
        }
    }

    constructor(
        private authService: AuthenticationService,
        protected templateRef: TemplateRef<any>,
        protected viewContainerRef: ViewContainerRef,
    ) { }
}
