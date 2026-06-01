import { Location } from '@angular/common';
import { Component, DestroyRef, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ToastService } from 'src/app/services/toast.service';
import { Org } from '@zitadel/proto/zitadel/org_pb';
import { NewMgmtService } from 'src/app/services/new-mgmt.service';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';

const ROUTEPARAM = 'projectid';

@Component({
  selector: 'cnsl-project-grant-create',
  templateUrl: './project-grant-create.component.html',
  styleUrls: ['./project-grant-create.component.scss'],
  standalone: false,
})
export class ProjectGrantCreateComponent implements OnInit {
  public org?: Org;
  public projectId: string = '';
  public grantId: string = '';
  public rolesKeyList: string[] = [];

  public createSteps: number = 2;
  public currentCreateStep: number = 1;

  constructor(
    private readonly route: ActivatedRoute,
    private readonly toast: ToastService,
    private readonly mgmtService: NewMgmtService,
    private readonly _location: Location,
    private readonly breadcrumbService: BreadcrumbService,
    private readonly destroyRef: DestroyRef,
  ) {}

  public ngOnInit(): void {
    this.route.params.pipe(takeUntilDestroyed(this.destroyRef)).subscribe((params) => {
      this.projectId = params['projectid'];

      const breadcrumbs = [
        new Breadcrumb({
          type: BreadcrumbType.ORG,
          routerLink: ['/org'],
        }),
        new Breadcrumb({
          type: BreadcrumbType.PROJECT,
          name: '',
          param: { key: ROUTEPARAM, value: this.projectId },
          routerLink: ['/projects', this.projectId],
        }),
      ];
      this.breadcrumbService.setBreadcrumb(breadcrumbs);
    });
  }

  public close(): void {
    this._location.back();
  }

  public addGrant(): void {
    if (this.org) {
      this.mgmtService
        .addProjectGrant({ grantedOrgId: this.org.id, projectId: this.projectId, roleKeys: this.rolesKeyList })
        .then(() => {
          this.close();
        })
        .catch((error) => {
          this.toast.showError(error);
        });
    }
  }

  public selectOrg(org: Org): void {
    this.org = org;
  }

  public selectRoles(roles: string[]): void {
    this.rolesKeyList = roles;
  }

  public next(): void {
    this.currentCreateStep++;
  }

  public previous(): void {
    this.currentCreateStep--;
  }
}
