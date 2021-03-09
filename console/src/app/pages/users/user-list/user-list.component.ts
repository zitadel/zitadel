import { Component } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { take } from 'rxjs/operators';
import { Type } from 'src/app/proto/generated/zitadel/user_pb';

@Component({
    selector: 'app-user-list',
    templateUrl: './user-list.component.html',
    styleUrls: ['./user-list.component.scss'],
})
export class UserListComponent {
    public Type: any = Type;
    public type: Type = Type.TYPE_HUMAN;

    constructor(public translate: TranslateService, activatedRoute: ActivatedRoute) {
        activatedRoute.data.pipe(take(1)).subscribe(params => {
            const { type } = params;
            this.type = type;
        });
    }
}
