import { Pipe, PipeTransform } from '@angular/core';
import { Observable } from 'rxjs';

import { AuthenticationService } from '../services/authentication.service';

@Pipe({
    name: 'hasRole',
})
export class HasRolePipe implements PipeTransform {
    constructor(private authService: AuthenticationService) { }

    public transform(values: string[]): Observable<boolean> {
        return this.authService.isAllowed(values);
    }
}
