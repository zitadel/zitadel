import { Component } from '@angular/core';
import { TranslateService } from '@ngx-translate/core';

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

    constructor(public translate: TranslateService) { }
}
