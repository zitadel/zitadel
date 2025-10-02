import { SelectionModel } from '@angular/cdk/collections';
import { Component, EventEmitter, Input, OnDestroy, OnInit, Output, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatTable } from '@angular/material/table';
import { combineLatestWith, firstValueFrom, Observable, ReplaySubject, Subject } from 'rxjs';
import { map, startWith, takeUntil } from 'rxjs/operators';
import { InstanceMembersDataSource } from 'src/app/pages/instance/instance-members/instance-members-datasource';
import { OrgMembersDataSource } from 'src/app/pages/orgs/org-members/org-members-datasource';
import { ProjectGrantMembersDataSource } from 'src/app/pages/projects/owned-projects/project-grant-detail/project-grant-members-datasource';
import { Member } from 'src/app/proto/generated/zitadel/member_pb';
import { getMembershipColor } from 'src/app/utils/color';

import { Type } from 'src/app/proto/generated/zitadel/user_pb';
import { AddMemberRolesDialogComponent } from '../add-member-roles-dialog/add-member-roles-dialog.component';
import { PageEvent, PaginatorComponent } from '../paginator/paginator.component';
import { ProjectMembersDataSource } from '../project-members/project-members-datasource';
import { WarnDialogComponent } from '../warn-dialog/warn-dialog.component';

type MemberDatasource =
  | OrgMembersDataSource
  | ProjectMembersDataSource
  | ProjectGrantMembersDataSource
  | InstanceMembersDataSource;

@Component({
  selector: 'cnsl-members-table',
  templateUrl: './members-table.component.html',
  styleUrls: ['./members-table.component.scss'],
  standalone: false,
})
export class MembersTableComponent implements OnInit, OnDestroy {
  public INITIALPAGESIZE: number = 25;
  @Input()
  public set canWrite(value: boolean | null) {
    this.canWrite$.next(!!value);
  }

  @Input()
  public set canDelete(value: boolean | null) {
    this.canDelete$.next(!!value);
  }

  @ViewChild(PaginatorComponent) public paginator!: PaginatorComponent;
  @ViewChild(MatTable) public table!: MatTable<Member.AsObject>;
  @Input() public dataSource?: MemberDatasource;
  public selection: SelectionModel<any> = new SelectionModel<any>(true, []);
  @Input() public memberRoleOptions: string[] = [];
  @Input() public factoryLoadFunc!: Function;
  @Input() public refreshTrigger!: Observable<void>;
  @Output() public updateRoles: EventEmitter<{ member: Member.AsObject; change: string[] }> = new EventEmitter();
  @Output() public changedSelection: EventEmitter<any[]> = new EventEmitter();
  @Output() public deleteMember: EventEmitter<Member.AsObject> = new EventEmitter();

  protected readonly displayedColumns$: Observable<string[]>;
  protected readonly canWrite$ = new ReplaySubject<boolean>(1);
  protected readonly canDelete$ = new ReplaySubject<boolean>(1);
  private destroyed: Subject<void> = new Subject();
  public UserType: any = Type;

  constructor(private dialog: MatDialog) {
    this.selection.changed.pipe(takeUntil(this.destroyed)).subscribe((_) => {
      this.changedSelection.emit(this.selection.selected);
    });

    this.displayedColumns$ = this.getDisplayedColumns();
  }

  public ngOnInit(): void {
    this.refreshTrigger.pipe(takeUntil(this.destroyed)).subscribe(() => {
      this.changePage(this.paginator);
    });
  }

  private getDisplayedColumns() {
    const defaultColumns = ['select', 'userId', 'displayName', 'loginname', 'email', 'roles'];
    return this.canWrite$.pipe(
      combineLatestWith(this.canDelete$),
      map(([canWrite, canDelete]) => {
        if (canWrite || canDelete) {
          return [...defaultColumns, 'actions'];
        }
        return defaultColumns;
      }),
      startWith(defaultColumns),
    );
  }

  public ngOnDestroy(): void {
    this.destroyed.next();
  }

  public getColor(role: string) {
    return getMembershipColor(role)[500];
  }

  public removeRole(member: Member.AsObject, role: string) {
    if (member.rolesList.length === 1) {
      this.triggerDeleteMember(member);
    } else {
      const dialogRef = this.dialog.open(WarnDialogComponent, {
        data: {
          confirmKey: 'ACTIONS.DELETE',
          cancelKey: 'ACTIONS.CANCEL',
          titleKey: 'ROLES.DIALOG.DELETE_TITLE',
          descriptionKey: 'ROLES.DIALOG.DELETE_DESCRIPTION',
        },
        width: '400px',
      });

      dialogRef.afterClosed().subscribe((resp) => {
        if (resp) {
          const newRoles = Object.assign([], member.rolesList);

          const index = newRoles.findIndex((r) => r === role);
          if (index > -1) {
            newRoles.splice(index, 1);
            member.rolesList = newRoles;
            this.updateRoles.emit({ member: member, change: newRoles });
          }
        }
      });
    }
  }

  public async addRole(member: Member.AsObject) {
    if (!(await firstValueFrom(this.canWrite$))) {
      return;
    }

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
    const numRows = this.dataSource?.membersSubject.value.length;
    return numSelected === numRows;
  }

  public masterToggle(): void {
    this.isAllSelected()
      ? this.selection.clear()
      : this.dataSource?.membersSubject.value.forEach((row: Member.AsObject) => this.selection.select(row));
  }

  public changePage(event?: PageEvent): any {
    this.selection.clear();
    return this.factoryLoadFunc(event ?? this.paginator);
  }

  public triggerDeleteMember(member: any): void {
    const dialogRef = this.dialog.open(WarnDialogComponent, {
      data: {
        confirmKey: 'ACTIONS.DELETE',
        cancelKey: 'ACTIONS.CANCEL',
        titleKey: 'MEMBER.DIALOG.DELETE_TITLE',
        descriptionKey: 'MEMBER.DIALOG.DELETE_DESCRIPTION',
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((resp) => {
      if (resp) {
        this.deleteMember.emit(member);
      }
    });
  }
}
