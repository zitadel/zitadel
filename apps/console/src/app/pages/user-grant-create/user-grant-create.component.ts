import { Location } from '@angular/common';
import { Component, OnDestroy } from '@angular/core';
import { ActivatedRoute, Params } from '@angular/router';
import { Subject, takeUntil } from 'rxjs';
import { ProjectType } from 'src/app/modules/project-members/project-members-datasource';
import { UserTarget } from 'src/app/modules/search-user-autocomplete/search-user-autocomplete.component';
import { UserGrantContext } from 'src/app/modules/user-grants/user-grants-datasource';
import { Org } from 'src/app/proto/generated/zitadel/org_pb';
import { GrantedProject, Project } from 'src/app/proto/generated/zitadel/project_pb';
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

  public org?: Org.AsObject;
  public userIds: string[] = [];

  public project?: Project.AsObject;
  public grantedProject?: GrantedProject.AsObject;

  public rolesList: string[] = [];

  public createSteps: number = 2;
  public currentCreateStep: number = 1;

  public UserGrantContext: any = UserGrantContext;

  public user?: User.AsObject;
  public UserTarget: any = UserTarget;

  public editState: boolean = false;
  private destroy$: Subject<void> = new Subject();

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
        type: BreadcrumbType.ORG,
        routerLink: ['/org'],
      }),
    ]);
    this.route.params.pipe(takeUntil(this.destroy$)).subscribe((params: Params) => {
      const { projectid, grantid, userid } = params;
      this.context = UserGrantContext.NONE;

      this.userIds = userid ? [userid] : [];

      if (projectid && !grantid) {
        this.context = UserGrantContext.OWNED_PROJECT;

        this.mgmtService
          .getProjectByID(projectid)
          .then((resp) => {
            if (resp.project) {
              this.project = resp.project;
            }
          })
          .catch((error: any) => {
            this.toast.showError(error);
          });
      } else if (projectid && grantid) {
        this.context = UserGrantContext.GRANTED_PROJECT;
        this.mgmtService
          .getGrantedProjectByID(projectid, grantid)
          .then((resp) => {
            if (resp.grantedProject) {
              this.grantedProject = resp.grantedProject;
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
        const prom = this.userIds.map((id) => this.userService.addUserGrant(id, this.rolesList, this.project?.id));
        Promise.all(prom)
          .then(() => {
            this.toast.showInfo('GRANTS.TOAST.UPDATED', true);
            this.close();
          })
          .catch((error: any) => {
            this.toast.showError(error);
            this.close();
          });
        break;
      case UserGrantContext.GRANTED_PROJECT:
        const promp = this.userIds.map((id) =>
          this.userService.addUserGrant(id, this.rolesList, this.grantedProject?.projectId, this.grantedProject?.grantId),
        );
        Promise.all(promp)
          .then(() => {
            this.toast.showInfo('GRANTS.TOAST.UPDATED', true);
            this.close();
          })
          .catch((error: any) => {
            this.toast.showError(error);
            this.close();
          });
        break;
      case UserGrantContext.USER:
        let grantId: string = '';
        let grantedProjectId: string = '';

        if (this.grantedProject?.grantId) {
          grantId = this.grantedProject.grantId;
          grantedProjectId = this.grantedProject.projectId;
        }

        const promu = this.userIds.map((id) =>
          this.userService.addUserGrant(
            id,
            this.rolesList,
            this.project?.id ? this.project.id : grantedProjectId ? grantedProjectId : '',
            grantId,
          ),
        );
        Promise.all(promu)
          .then(() => {
            this.toast.showInfo('GRANTS.TOAST.UPDATED', true);
            this.close();
          })
          .catch((error: any) => {
            this.toast.showError(error);
            this.close();
          });
        break;
      case UserGrantContext.NONE:
        let tempGrantId: string = '';

        if (this.grantedProject?.grantId) {
          tempGrantId = this.grantedProject.grantId;
        }

        const promn = this.userIds.map((id) =>
          this.userService.addUserGrant(
            id,
            this.rolesList,
            this.project ? this.project.id : this.grantedProject ? this.grantedProject.projectId : '',
            tempGrantId,
          ),
        );
        Promise.all(promn)
          .then(() => {
            this.toast.showInfo('GRANTS.TOAST.UPDATED', true);
            this.close();
          })
          .catch((error: any) => {
            this.toast.showError(error);
            this.close();
          });
        break;
    }
  }

  public selectProject(project: Project.AsObject | GrantedProject.AsObject, type: ProjectType): void {
    if (type === ProjectType.PROJECTTYPE_OWNED) {
      this.project = project as Project.AsObject;
    } else if (type === ProjectType.PROJECTTYPE_GRANTED) {
      this.grantedProject = project as GrantedProject.AsObject;
    }
  }

  public selectUsers(user: User.AsObject[]): void {
    if (user && user.length) {
      this.userIds = (user as User.AsObject[]).map((u) => u.id);
    }
  }

  public selectRoles(roleKeys: string[]): void {
    this.rolesList = roleKeys;
  }

  public next(): void {
    this.currentCreateStep++;
  }

  public previous(): void {
    this.currentCreateStep--;
  }

  public ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }
}
