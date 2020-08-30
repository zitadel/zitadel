import { Component } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { take } from 'rxjs/operators';

export enum UserType {
    HUMAN = 'human',
    MACHINE = 'machine',
}
@Component({
    selector: 'app-user-list',
    templateUrl: './user-list.component.html',
    styleUrls: ['./user-list.component.scss'],
})
export class UserListComponent {
    public UserType: any = UserType;
    public type: UserType = UserType.HUMAN;
    constructor(public translate: TranslateService, activatedRoute: ActivatedRoute) {
        activatedRoute.data.pipe(take(1)).subscribe(params => {
            const { type } = params;
            this.type = type;
        });
    }
}
