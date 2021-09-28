import { Component } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { take } from 'rxjs/operators';

enum UserType {
  USER = 'user',
  MACHINE = 'machine'
}
@Component({
  selector: 'app-project-list',
  templateUrl: './project-list.component.html',
  styleUrls: ['./project-list.component.scss'],
})
export class ProjectListComponent {
  public UserType: any = UserType;
  public type: UserType = UserType.USER;

  constructor(public translate: TranslateService, activatedRoute: ActivatedRoute) {
    activatedRoute.data.pipe(take(1)).subscribe(params => {
      const { type } = params;
      this.type = type;
    });
  }
}
