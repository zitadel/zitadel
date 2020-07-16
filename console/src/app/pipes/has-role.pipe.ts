import { Pipe, PipeTransform } from '@angular/core';
import { Observable } from 'rxjs';

import { AuthUserService } from '../services/auth-user.service';

@Pipe({
    name: 'hasRole',
})
export class HasRolePipe implements PipeTransform {

    constructor(private authUserService: AuthUserService) { }

    public transform(values: string[], each: boolean = false): Observable<boolean> {
        return this.authUserService.isAllowed(values, each);
    }
}
