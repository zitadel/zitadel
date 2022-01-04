import { SelectionModel } from '@angular/cdk/collections';
import { Component, EventEmitter, Input, OnDestroy, OnInit, Output, ViewChild } from '@angular/core';
import { MatSelectChange } from '@angular/material/select';
import { MatTable } from '@angular/material/table';
import { Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';
import { Membership } from 'src/app/proto/generated/zitadel/user_pb';
import { AdminService } from 'src/app/services/admin.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { PageEvent, PaginatorComponent } from '../paginator/paginator.component';
import { MembershipsDataSource } from './memberships-datasource';

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
  public dataSource!: MembershipsDataSource;
  public selection: SelectionModel<any> = new SelectionModel<any>(true, []);

  @Output() public updateRoles: EventEmitter<{ member: Membership.AsObject; change: MatSelectChange }> = new EventEmitter();
  @Output() public changedSelection: EventEmitter<any[]> = new EventEmitter();
  @Output() public deleteMembership: EventEmitter<Membership.AsObject> = new EventEmitter();

  private destroyed: Subject<void> = new Subject();
  public membershipRoleOptions: string[] = [];

  public displayedColumns: string[] = ['select', 'displayName', 'orgId', 'rolesList'];
  public membershipToEdit: string = '';

  constructor(
    private authService: GrpcAuthService,
    private toastService: ToastService,
    private mgmtService: ManagementService,
    private adminService: AdminService,
  ) {
    this.dataSource = new MembershipsDataSource(this.authService, this.mgmtService);

    this.selection.changed.pipe(takeUntil(this.destroyed)).subscribe((_) => {
      this.changedSelection.emit(this.selection.selected);
    });
  }

  public ngOnInit(): void {
    // this.refreshTrigger.pipe(takeUntil(this.destroyed)).subscribe(() => {
    this.changePage(this.paginator);
    // });
  }

  public loadRoles(membership: Membership.AsObject): void {
    if (membership.orgId) {
      this.membershipToEdit = `${membership.orgId}${membership.projectId}${membership.projectGrantId}`;
      this.mgmtService
        .listOrgMemberRoles()
        .then((resp) => {
          this.membershipRoleOptions = resp.resultList;
        })
        .catch((error) => {
          this.toastService.showError(error);
        });
    } else if (membership.projectGrantId) {
      this.membershipToEdit = `${membership.orgId}${membership.projectId}${membership.projectGrantId}`;
      this.mgmtService
        .listProjectMemberRoles()
        .then((resp) => {
          this.membershipRoleOptions = resp.resultList;
        })
        .catch((error) => {
          this.toastService.showError(error);
        });
    } else if (membership.projectId) {
      this.membershipToEdit = `${membership.orgId}${membership.projectId}${membership.projectGrantId}`;
      this.mgmtService
        .listProjectGrantMemberRoles()
        .then((resp) => {
          this.membershipRoleOptions = resp.resultList;
        })
        .catch((error) => {
          this.toastService.showError(error);
        });
    } else if (membership.iam) {
      this.membershipToEdit = `IAM`;
      this.adminService
        .listIAMMemberRoles()
        .then((resp) => {
          console.log(resp);
          this.membershipRoleOptions = resp.rolesList;
        })
        .catch((error) => {
          this.toastService.showError(error);
        });
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
}
