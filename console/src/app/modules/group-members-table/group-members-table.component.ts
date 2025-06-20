import { SelectionModel } from '@angular/cdk/collections';
import { Component, EventEmitter, Input, OnDestroy, OnInit, Output, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatTable } from '@angular/material/table';
import { Observable, Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';
import { Member } from 'src/app/proto/generated/zitadel/member_pb';
import { Type } from 'src/app/proto/generated/zitadel/user_pb';
import { AddMemberRolesDialogComponent } from '../add-member-roles-dialog/add-member-roles-dialog.component';
import { PageEvent, PaginatorComponent } from '../paginator/paginator.component';
import { GroupMembersDataSource } from 'src/app/pages/groups/group-detail/group-detail/group-members-datasource';
import { WarnDialogComponent } from '../warn-dialog/warn-dialog.component';

type MemberDatasource = GroupMembersDataSource;

@Component({
  selector: 'cnsl-group-members-table',
  templateUrl: './group-members-table.component.html',
  styleUrls: ['./group-members-table.component.scss'],
})
export class GroupMembersTableComponent implements OnInit, OnDestroy {
  public INITIALPAGESIZE: number = 25;
  @Input() public canDelete: boolean | null = false;
  @Input() public canWrite: boolean | null = false;
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

  private destroyed: Subject<void> = new Subject();
  public displayedColumns: string[] = ['select', 'userId', 'displayName', 'loginname', 'email'];
  public UserType: any = Type;

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
        titleKey: 'GROUP.MEMBER.DELETE_TITLE',
        descriptionKey: 'GROUP.MEMBER.DELETE_DESCRIPTION',
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
