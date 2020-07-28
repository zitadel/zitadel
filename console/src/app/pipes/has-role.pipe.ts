import { Pipe, PipeTransform } from '@angular/core';
import { Observable } from 'rxjs';

import { AuthService } from '../services/auth.service';

@Pipe({
    name: 'hasRole',
})
export class HasRolePipe implements PipeTransform {

    constructor(private authService: AuthService) { }

    public transform(values: string[], each: boolean = false): Observable<boolean> {
        return this.authService.isAllowed(values, each);
    }
}
