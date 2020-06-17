import { Pipe, PipeTransform } from '@angular/core';

import { PasswordComplexityPolicy } from '../proto/generated/management_pb';
import { OrgService } from '../services/org.service';

@Pipe({
    name: 'passwordPattern',
})
export class PasswordPatternPipe implements PipeTransform {

    constructor(private orgService: OrgService) { }

    transform(policy: PasswordComplexityPolicy.AsObject, ...args: unknown[]): string {
        return this.orgService.getLocalizedComplexityPolicyPatternErrorString(policy);
    }

}
