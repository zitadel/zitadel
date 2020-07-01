import { Pipe, PipeTransform } from '@angular/core';

import { AuthUserService } from '../services/auth-user.service';

@Pipe({
    name: 'hasRole',
    pure: false,
})
export class HasRolePipe implements PipeTransform {

    constructor(private authUserService: AuthUserService) { }

    public transform(values: string[], each: boolean = false): any {
        return this.authUserService.isAllowed(values, each);
    }
}
