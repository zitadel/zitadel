import { Location } from '@angular/common';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Subject, takeUntil } from 'rxjs';
import { Org } from 'src/app/proto/generated/zitadel/org_pb';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

const ROUTEPARAM = 'projectid';

@Component({
  selector: 'cnsl-project-grant-create',
  templateUrl: './project-grant-create.component.html',
  styleUrls: ['./project-grant-create.component.scss'],
})
export class ProjectGrantCreateComponent implements OnInit, OnDestroy {
  public org?: Org.AsObject;
  public projectId: string = '';
  public grantId: string = '';
  public rolesKeyList: string[] = [];

  public createSteps: number = 2;
  public currentCreateStep: number = 1;

  private destroy$: Subject<void> = new Subject();
  constructor(
    private route: ActivatedRoute,
    private toast: ToastService,
    private mgmtService: ManagementService,
    private _location: Location,
    private breadcrumbService: BreadcrumbService,
  ) {}

  public ngOnInit(): void {
    this.route.params.pipe(takeUntil(this.destroy$)).subscribe((params) => {
      this.projectId = params.projectid;

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

  public ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  public searchOrg(domain: string): void {
    this.mgmtService
      .getOrgByDomainGlobal(domain)
      .then((ret) => {
        if (ret.org) {
          const tmp = ret.org;
          this.org = tmp;
        }
      })
      .catch((error) => {
        this.toast.showError(error);
      });
  }

  public close(): void {
    this._location.back();
  }

  public addGrant(): void {
    if (this.org) {
      this.mgmtService
        .addProjectGrant(this.org.id, this.projectId, this.rolesKeyList)
        .then(() => {
          this.close();
        })
        .catch((error) => {
          this.toast.showError(error);
        });
    }
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
