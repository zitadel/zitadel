import { Location } from '@angular/common';
import { Component, OnDestroy } from '@angular/core';
import { ActivatedRoute, Params } from '@angular/router';
import { Subscription } from 'rxjs';
import { UserTarget } from 'src/app/modules/search-user-autocomplete/search-user-autocomplete.component';
import { UserGrantContext } from 'src/app/modules/user-grants/user-grants-datasource';
import { Org } from 'src/app/proto/generated/zitadel/org_pb';
import { GrantedProject, Project, Role } from 'src/app/proto/generated/zitadel/project_pb';
import { User } from 'src/app/proto/generated/zitadel/user_pb';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { StorageKey, StorageLocation, StorageService } from 'src/app/services/storage.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
  selector: 'cnsl-user-grant-create',
  templateUrl: './user-grant-create.component.html',
  styleUrls: ['./user-grant-create.component.scss'],
})
export class UserGrantCreateComponent implements OnDestroy {
  public context!: UserGrantContext;

  public org!: Org.AsObject;
  public userIds: string[] = [];

  public projectId: string = '';
  public project!: GrantedProject.AsObject | Project.AsObject;

  public grantId: string = '';
  public rolesList: string[] = [];

  public STEPS: number = 2; // project, roles
  public currentCreateStep: number = 1;

  public filterValue: string = '';

  private subscription: Subscription = new Subscription();

  public UserGrantContext: any = UserGrantContext;

  public user!: User.AsObject;
  public UserTarget: any = UserTarget;

  public editState: boolean = false;

  constructor(
    private userService: ManagementService,
    private toast: ToastService,
    private _location: Location,
    private route: ActivatedRoute,
    private mgmtService: ManagementService,
    private storage: StorageService,
    breadcrumbService: BreadcrumbService,
  ) {
    breadcrumbService.setBreadcrumb([
      new Breadcrumb({
        type: BreadcrumbType.IAM,
        name: 'IAM',
        routerLink: ['/system'],
      }),
      new Breadcrumb({
        type: BreadcrumbType.ORG,
        routerLink: ['/org'],
      }),
    ]);
    this.subscription = this.route.params.subscribe((params: Params) => {
      const { projectid, grantid, userid } = params;
      this.context = UserGrantContext.NONE;

      this.projectId = projectid;
      this.grantId = grantid;
      this.userIds = userid ? [userid] : [];

      if (this.projectId && !this.grantId) {
        this.context = UserGrantContext.OWNED_PROJECT;

        this.mgmtService
          .getProjectByID(this.projectId)
          .then((resp) => {
            if (resp.project) {
              this.project = resp.project;
            }
          })
          .catch((error: any) => {
            this.toast.showError(error);
          });
      } else if (this.projectId && this.grantId) {
        this.context = UserGrantContext.GRANTED_PROJECT;
        this.mgmtService
          .getGrantedProjectByID(this.projectId, this.grantId)
          .then((resp) => {
            if (resp.grantedProject) {
              this.project = resp.grantedProject;
            }
          })
          .catch((error: any) => {
            this.toast.showError(error);
          });
      } else if (this.userIds && this.userIds.length === 1) {
        this.context = UserGrantContext.USER;
        this.mgmtService
          .getUserByID(this.userIds[0])
          .then((resp) => {
            if (resp.user) {
              this.user = resp.user;
            }
          })
          .catch((error: any) => {
            this.toast.showError(error);
          });
      }
    });

    const temporg = this.storage.getItem<Org.AsObject>(StorageKey.organization, StorageLocation.session);
    if (temporg) {
      this.org = temporg;
    }
  }

  public close(): void {
    this._location.back();
  }

  public addGrant(): void {
    switch (this.context) {
      case UserGrantContext.OWNED_PROJECT:
        const prom = this.userIds.map((id) => this.userService.addUserGrant(id, this.rolesList, this.projectId));
        Promise.all(prom)
          .then(() => {
            this.toast.showInfo('PROJECT.GRANT.TOAST.PROJECTGRANTADDED', true);
            this.close();
          })
          .catch((error: any) => {
            this.toast.showError(error);
            this.close();
          });
        break;
      case UserGrantContext.GRANTED_PROJECT:
        const promp = this.userIds.map((id) =>
          this.userService.addUserGrant(id, this.rolesList, this.projectId, this.grantId),
        );
        Promise.all(promp)
          .then(() => {
            this.toast.showInfo('PROJECT.GRANT.TOAST.PROJECTGRANTUSERGRANTADDED', true);
            this.close();
          })
          .catch((error: any) => {
            this.toast.showError(error);
            this.close();
          });
        break;
      case UserGrantContext.USER:
        let grantId: string = '';

        if ((this.project as GrantedProject.AsObject)?.grantId) {
          grantId = (this.project as GrantedProject.AsObject).grantId;
        }

        const promu = this.userIds.map((id) => this.userService.addUserGrant(id, this.rolesList, this.projectId, grantId));
        Promise.all(promu)
          .then(() => {
            this.toast.showInfo('PROJECT.GRANT.TOAST.PROJECTGRANTUSERGRANTADDED', true);
            this.close();
          })
          .catch((error: any) => {
            this.toast.showError(error);
            this.close();
          });
        break;
      case UserGrantContext.NONE:
        let tempGrantId: string = '';

        if ((this.project as GrantedProject.AsObject)?.grantId) {
          tempGrantId = (this.project as GrantedProject.AsObject).grantId;
        }

        const promn = this.userIds.map((id) =>
          this.userService.addUserGrant(id, this.rolesList, this.projectId, tempGrantId),
        );
        Promise.all(promn)
          .then(() => {
            this.toast.showInfo('PROJECT.GRANT.TOAST.PROJECTGRANTUSERGRANTADDED', true);
            this.close();
          })
          .catch((error: any) => {
            this.toast.showError(error);
            this.close();
          });
        break;
    }
  }

  public selectProject(project: Project.AsObject | GrantedProject.AsObject | any): void {
    this.project = project;
    this.projectId = project.id || project.projectId;
  }

  public selectUsers(user: User.AsObject[]): void {
    if (user && user.length) {
      this.userIds = (user as User.AsObject[]).map((u) => u.id);
    }
  }

  public selectRoles(roles: Role.AsObject[]): void {
    this.rolesList = roles.map((role) => role.key);
  }

  public next(): void {
    this.currentCreateStep++;
  }

  public previous(): void {
    this.currentCreateStep--;
  }

  public ngOnDestroy(): void {
    this.subscription.unsubscribe();
  }
}
