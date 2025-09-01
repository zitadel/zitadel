import { Location } from '@angular/common';
import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { MessageInitShape } from '@bufbuild/protobuf';
import { AddProjectRequestSchema } from '@zitadel/proto/zitadel/management_pb';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ToastService } from 'src/app/services/toast.service';
import { NewMgmtService } from 'src/app/services/new-mgmt.service';
import { UserService } from 'src/app/services/user.service';

@Component({
  selector: 'cnsl-project-create',
  templateUrl: './project-create.component.html',
  styleUrls: ['./project-create.component.scss'],
})
export class ProjectCreateComponent {
  protected readonly project: MessageInitShape<typeof AddProjectRequestSchema> = {
    name: '',
    admins: [
      {
        userId: this.userService.userId(),
      },
    ],
  };

  constructor(
    private readonly router: Router,
    private readonly toast: ToastService,
    private readonly newMgmtService: NewMgmtService,
    private readonly _location: Location,
    private readonly userService: UserService,
    breadcrumbService: BreadcrumbService,
  ) {
    const bread: Breadcrumb = {
      type: BreadcrumbType.ORG,
      routerLink: ['/org'],
    };
    breadcrumbService.setBreadcrumb([bread]);
  }

  public saveProject(): void {
    this.newMgmtService
      .addProject(this.project)
      .then((resp) => {
        this.toast.showInfo('PROJECT.TOAST.CREATED', true);
        return this.router.navigate(['projects', resp.id], { queryParams: { new: true } });
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public close(): void {
    this._location.back();
  }
}
