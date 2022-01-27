import { SelectionModel } from '@angular/cdk/collections';
import { Component, EventEmitter, Input, OnDestroy, OnInit, Output, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatTable } from '@angular/material/table';
import { Observable, Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';
import { IamMembersDataSource } from 'src/app/pages/iam/iam-members/iam-members-datasource';
import { OrgMembersDataSource } from 'src/app/pages/orgs/org-members/org-members-datasource';
import {
  ProjectGrantMembersDataSource,
} from 'src/app/pages/projects/owned-projects/project-grant-detail/project-grant-members-datasource';
import { Member } from 'src/app/proto/generated/zitadel/member_pb';

import { AddMemberRolesDialogComponent } from '../add-member-roles-dialog/add-member-roles-dialog.component';
import { getColor } from '../avatar/avatar.component';
import { PageEvent, PaginatorComponent } from '../paginator/paginator.component';
import { ProjectMembersDataSource } from '../project-members/project-members-datasource';

type MemberDatasource =
  | OrgMembersDataSource
  | ProjectMembersDataSource
  | ProjectGrantMembersDataSource
  | IamMembersDataSource;

@Component({
  selector: 'cnsl-members-table',
  templateUrl: './members-table.component.html',
  styleUrls: ['./members-table.component.scss'],
})
export class MembersTableComponent implements OnInit, OnDestroy {
  public INITIALPAGESIZE: number = 25;
  @Input() public canDelete: boolean | null = false;
  @Input() public canWrite: boolean | null = false;
  @ViewChild(PaginatorComponent) public paginator!: PaginatorComponent;
  @ViewChild(MatTable) public table!: MatTable<Member.AsObject>;
  @Input() public dataSource!: MemberDatasource;
  public selection: SelectionModel<any> = new SelectionModel<any>(true, []);
  @Input() public memberRoleOptions: string[] = [];
  @Input() public factoryLoadFunc!: Function;
  @Input() public refreshTrigger!: Observable<void>;
  @Output() public updateRoles: EventEmitter<{ member: Member.AsObject; change: string[] }> = new EventEmitter();
  @Output() public changedSelection: EventEmitter<any[]> = new EventEmitter();
  @Output() public deleteMember: EventEmitter<Member.AsObject> = new EventEmitter();

  private destroyed: Subject<void> = new Subject();
  public displayedColumns: string[] = ['select', 'userId', 'displayName', 'loginname', 'email', 'roles'];

  constructor(private dialog: MatDialog) {
    this.selection.changed.pipe(takeUntil(this.destroyed)).subscribe((_) => {
      this.changedSelection.emit(this.selection.selected);
    });
  }

  public ngOnInit(): void {
    this.refreshTrigger.pipe(takeUntil(this.destroyed)).subscribe(() => {
      this.changePage(this.paginator);
    });

    if (this.canDelete || this.canWrite) {
      this.displayedColumns.push('actions');
    }
  }

  public ngOnDestroy(): void {
    this.destroyed.next();
  }

  public getColor(role: string) {
    return getColor(role);
  }

  public removeRole(member: Member.AsObject, role: string) {
    const newRoles = Object.assign([], member.rolesList);
    const index = newRoles.findIndex((r) => r === role);
    if (index > -1) {
      newRoles.splice(index);
      member.rolesList = newRoles;
      this.updateRoles.emit({ member: member, change: newRoles });
    }
  }

  public addRole(member: Member.AsObject) {
    const dialogRef = this.dialog.open(AddMemberRolesDialogComponent, {
      data: {
        user: member.displayName,
        allRoles: this.memberRoleOptions,
        selectedRoles: member.rolesList,
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((resp) => {
      if (resp && resp.length) {
        member.rolesList = resp;
        this.updateRoles.emit({ member: member, change: resp });
      }
    });
  }

  public isAllSelected(): boolean {
    const numSelected = this.selection.selected.length;
    const numRows = this.dataSource.membersSubject.value.length;
    return numSelected === numRows;
  }

  public masterToggle(): void {
    this.isAllSelected()
      ? this.selection.clear()
      : this.dataSource.membersSubject.value.forEach((row) => this.selection.select(row));
  }

  public changePage(event?: PageEvent): any {
    this.selection.clear();
    return this.factoryLoadFunc(event ?? this.paginator);
  }

  public triggerDeleteMember(member: any): void {
    this.deleteMember.emit(member);
  }
}
