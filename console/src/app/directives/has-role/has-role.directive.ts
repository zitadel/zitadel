import { Directive, Input, TemplateRef, ViewContainerRef } from '@angular/core';
import { AuthUserService } from 'src/app/services/auth-user.service';


@Directive({
    selector: '[appHasRole]',
})

export class HasRoleDirective {
    private hasView: boolean = false;
    @Input() public set appHasRole(roles: string[]) {
        if (roles && roles.length > 0) {
            this.userService.isAllowed(roles).subscribe(isAllowed => {
                if (isAllowed && !this.hasView) {
                    this.viewContainerRef.clear();
                    this.viewContainerRef.createEmbeddedView(this.templateRef);
                } else if (this.hasView) {
                    console.log('User blocked!', roles, isAllowed);
                    this.viewContainerRef.clear();
                    this.hasView = false;
                }
            });
        }
    }

    constructor(
        private userService: AuthUserService,
        protected templateRef: TemplateRef<any>,
        protected viewContainerRef: ViewContainerRef,
    ) { }
}
