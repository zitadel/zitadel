import { SelectionModel } from '@angular/cdk/collections';
import { Component, EventEmitter, Input, OnDestroy, OnInit, Output, ViewChild } from '@angular/core';
import { MatTable } from '@angular/material/table';
import { Router } from '@angular/router';
import { Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';
import { Membership } from 'src/app/proto/generated/zitadel/user_pb';
import { AdminService } from 'src/app/services/admin.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { OverlayWorkflowService } from 'src/app/services/overlay/overlay-workflow.service';
import { OrgContextChangedWorkflowOverlays } from 'src/app/services/overlay/workflows';
import { ToastService } from 'src/app/services/toast.service';
import { getMembershipColor } from 'src/app/utils/color';

import { PageEvent, PaginatorComponent } from '../paginator/paginator.component';
import { MembershipsDataSource } from './memberships-datasource';
import { NewOrganizationService } from '../../services/new-organization.service';
import { Organization } from '@zitadel/proto/zitadel/org/v2/org_pb';

@Component({
  selector: 'cnsl-memberships-table',
  templateUrl: './memberships-table.component.html',
  styleUrls: ['./memberships-table.component.scss'],
})
export class MembershipsTableComponent implements OnInit, OnDestroy {
  public INITIALPAGESIZE: number = 25;
  @ViewChild(PaginatorComponent) public paginator!: PaginatorComponent;
  @ViewChild(MatTable) public table!: MatTable<Membership.AsObject>;
  @Input() public userId: string = '';
  public dataSource: MembershipsDataSource = new MembershipsDataSource(this.authService, this.mgmtService);
  public selection: SelectionModel<Membership.AsObject> = new SelectionModel<Membership.AsObject>(true, []);

  @Output() public changedSelection: EventEmitter<any[]> = new EventEmitter();
  @Output() public deleteMembership: EventEmitter<Membership.AsObject> = new EventEmitter();

  private destroyed: Subject<void> = new Subject();
  public membershipRoleOptions: string[] = [];

  public displayedColumns: string[] = ['displayName', 'type', 'rolesList'];
  public membershipToEdit: string = '';
  public loadingRoles: boolean = false;

  constructor(
    private authService: GrpcAuthService,
    private toastService: ToastService,
    private mgmtService: ManagementService,
    private adminService: AdminService,
    private toast: ToastService,
    private router: Router,
    private workflowService: OverlayWorkflowService,
    private readonly newOrganizationService: NewOrganizationService,
  ) {
    this.selection.changed.pipe(takeUntil(this.destroyed)).subscribe((_) => {
      this.changedSelection.emit(this.selection.selected);
    });
  }

  public ngOnInit(): void {
    this.changePage(this.paginator);
  }

  public loadRoles(membership: Membership.AsObject, opened: boolean): void {
    if (opened) {
      this.loadingRoles = true;

      if (membership.orgId && !membership.projectId && !membership.projectGrantId) {
        this.membershipToEdit = `${membership.orgId}${membership.projectId}${membership.projectGrantId}`;
        this.mgmtService
          .listOrgMemberRoles()
          .then((resp) => {
            this.membershipRoleOptions = resp.resultList;
            this.loadingRoles = false;
          })
          .catch((error) => {
            this.toastService.showError(error);
            this.loadingRoles = false;
          });
      } else if (membership.projectGrantId) {
        this.membershipToEdit = `${membership.orgId}${membership.projectId}${membership.projectGrantId}`;
        this.mgmtService
          .listProjectGrantMemberRoles()
          .then((resp) => {
            this.membershipRoleOptions = resp.resultList;
            this.loadingRoles = false;
          })
          .catch((error) => {
            this.toastService.showError(error);
            this.loadingRoles = false;
          });
      } else if (membership.projectId) {
        this.membershipToEdit = `${membership.orgId}${membership.projectId}${membership.projectGrantId}`;
        this.mgmtService
          .listProjectMemberRoles()
          .then((resp) => {
            this.membershipRoleOptions = resp.resultList;
            this.loadingRoles = false;
          })
          .catch((error) => {
            this.toastService.showError(error);
            this.loadingRoles = false;
          });
      } else if (membership.iam) {
        this.membershipToEdit = `IAM`;
        this.adminService
          .listIAMMemberRoles()
          .then((resp) => {
            this.membershipRoleOptions = resp.rolesList;
            this.loadingRoles = false;
          })
          .catch((error) => {
            this.toastService.showError(error);
            this.loadingRoles = false;
          });
      }
    }
  }

  public async goto(membership: Membership.AsObject) {
    const orgId = this.newOrganizationService.orgId();

    if (membership.orgId && !membership.projectId && !membership.projectGrantId) {
      // only shown on auth user, or if currentOrg === resourceOwner
      try {
        const membershipOrg = await this.newOrganizationService.setOrgId(membership.orgId);
        await this.router.navigate(['/org/members']);
        this.startOrgContextWorkflow(membershipOrg, orgId);
      } catch (error) {
        this.toast.showInfo('USER.MEMBERSHIPS.NOPERMISSIONTOEDIT', true);
      }
    } else if (membership.projectGrantId && membership.details?.resourceOwner) {
      // only shown on auth user
      try {
        const membershipOrg = await this.newOrganizationService.setOrgId(membership.details?.resourceOwner);
        await this.router.navigate(['/granted-projects', membership.projectId, 'grants', membership.projectGrantId]);
        this.startOrgContextWorkflow(membershipOrg, orgId);
      } catch (error) {
        this.toast.showInfo('USER.MEMBERSHIPS.NOPERMISSIONTOEDIT', true);
      }
    } else if (membership.projectId && membership.details?.resourceOwner) {
      // only shown on auth user, or if currentOrg === resourceOwner
      try {
        const membershipOrg = await this.newOrganizationService.setOrgId(membership.details?.resourceOwner);
        await this.router.navigate(['/projects', membership.projectId, 'members']);
        this.startOrgContextWorkflow(membershipOrg, orgId);
      } catch (error) {
        this.toast.showInfo('USER.MEMBERSHIPS.NOPERMISSIONTOEDIT', true);
      }
    } else if (membership.iam) {
      // only shown on auth user
      await this.router.navigate(['/instance/members']);
    }
  }

  private startOrgContextWorkflow(membershipOrg: Organization, currentOrgId?: string | null): void {
    if (!currentOrgId || (membershipOrg.id && currentOrgId && currentOrgId !== membershipOrg.id)) {
      setTimeout(() => {
        this.workflowService.startWorkflow(OrgContextChangedWorkflowOverlays, null);
      }, 1000);
    }
  }

  public getType(membership: Membership.AsObject): string {
    if (membership.orgId && !membership.projectId && !membership.projectGrantId) {
      return 'Organization';
    } else if (membership.projectGrantId) {
      return 'Project Grant';
    } else if (membership.projectId) {
      return 'Project';
    } else if (membership.iam) {
      return 'IAM';
    } else {
      return '';
    }
  }

  public ngOnDestroy(): void {
    this.destroyed.next();
  }

  public isAllSelected(): boolean {
    const numSelected = this.selection.selected.length;
    const numRows = this.dataSource.membershipsSubject.value.length;
    return numSelected === numRows;
  }

  public masterToggle(): void {
    this.isAllSelected()
      ? this.selection.clear()
      : this.dataSource.membershipsSubject.value.forEach((row) => this.selection.select(row));
  }

  public changePage(event?: PageEvent): any {
    this.selection.clear();
    return this.userId
      ? this.dataSource.loadMemberships(this.userId, event?.pageIndex ?? 0, event?.pageSize ?? this.INITIALPAGESIZE)
      : this.dataSource.loadMyMemberships(event?.pageIndex ?? 0, event?.pageSize ?? this.INITIALPAGESIZE);
  }

  public isCurrentMembership(membership: Membership.AsObject): boolean {
    return (
      this.membershipToEdit ===
      (membership.iam ? 'IAM' : `${membership.orgId}${membership.projectId}${membership.projectGrantId}`)
    );
  }

  public getColor(role: string): string {
    const color = getMembershipColor(role);
    return color[500];
  }

  public removeRole(membership: Membership.AsObject, role: string): void {
    const newRoles = Object.assign([], membership.rolesList);
    const index = newRoles.findIndex((r) => r === role);
    if (index > -1) {
      newRoles.splice(index);
      if (membership.orgId) {
        console.log('org member', membership.userId, newRoles);
        this.mgmtService
          .updateOrgMember(membership.userId, newRoles)
          .then(() => {
            this.toast.showInfo('USER.MEMBERSHIPS.UPDATED', true);
            this.changePage(this.paginator);
          })
          .catch((error) => {
            this.toastService.showError(error);
          });
      } else if (membership.projectGrantId) {
        this.mgmtService
          .updateProjectGrantMember(membership.projectId, membership.projectGrantId, membership.userId, newRoles)
          .then(() => {
            this.toast.showInfo('USER.MEMBERSHIPS.UPDATED', true);
            this.changePage(this.paginator);
          })
          .catch((error) => {
            this.toastService.showError(error);
          });
      } else if (membership.projectId) {
        console.log(membership.projectId, membership.userId, newRoles);
        this.mgmtService
          .updateProjectMember(membership.projectId, membership.userId, newRoles)
          .then(() => {
            this.toast.showInfo('USER.MEMBERSHIPS.UPDATED', true);
            this.changePage(this.paginator);
          })
          .catch((error) => {
            this.toastService.showError(error);
          });
      } else if (membership.iam) {
        this.adminService
          .updateIAMMember(membership.userId, newRoles)
          .then(() => {
            this.toast.showInfo('USER.MEMBERSHIPS.UPDATED', true);
            this.changePage(this.paginator);
          })
          .catch((error) => {
            this.toastService.showError(error);
          });
      }
    }
  }
}
