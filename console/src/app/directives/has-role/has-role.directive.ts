import { Directive, Input, TemplateRef, ViewContainerRef } from '@angular/core';
import { AuthService } from 'src/app/services/auth.service';


@Directive({
    selector: '[appHasRole]',
})

export class HasRoleDirective {
    private hasView: boolean = false;
    @Input() public isRegexp: boolean = false;
    @Input() public set appHasRole(roles: string[] | RegExp[]) {
        if (roles && roles.length > 0) {
            console.log('isRegexp', this.isRegexp, roles);
            this.authService.isAllowed(roles, false, this.isRegexp).subscribe(isAllowed => {
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
        private authService: AuthService,
        protected templateRef: TemplateRef<any>,
        protected viewContainerRef: ViewContainerRef,
    ) { }
}
