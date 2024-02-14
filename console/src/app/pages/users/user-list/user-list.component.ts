import { Component } from '@angular/core';
import { ActivatedRoute, Params } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { take } from 'rxjs/operators';
import { Type } from 'src/app/proto/generated/zitadel/user_pb';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';

@Component({
  selector: 'cnsl-user-list',
  templateUrl: './user-list.component.html',
  styleUrls: ['./user-list.component.scss'],
})
export class UserListComponent {
  public Type: any = Type;
  public type: Type = Type.TYPE_HUMAN;

  constructor(
    public translate: TranslateService,
    activatedRoute: ActivatedRoute,
    breadcrumbService: BreadcrumbService,
  ) {
    activatedRoute.queryParams.pipe(take(1)).subscribe((params: Params) => {
      const { type } = params;
      if (type && type === 'human') {
        this.type = Type.TYPE_HUMAN;
      } else if (type && type === 'machine') {
        this.type = Type.TYPE_MACHINE;
      }
    });

    const bread: Breadcrumb = {
      type: BreadcrumbType.ORG,
      routerLink: ['/org'],
    };
    breadcrumbService.setBreadcrumb([bread]);
  }
}
