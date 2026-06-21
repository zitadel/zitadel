import { ChangeDetectionStrategy, Component, DestroyRef, signal } from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { MatDialog } from '@angular/material/dialog';
import { MessageInitShape } from '@bufbuild/protobuf';
import { FieldName, Group } from '@zitadel/proto/zitadel/group/v2/group_pb';
import { CreateGroupRequestSchema, UpdateGroupRequestSchema } from '@zitadel/proto/zitadel/group/v2/group_service_pb';
import { PageEvent } from '@angular/material/paginator';
import { BehaviorSubject, combineLatest, defer, firstValueFrom, lastValueFrom, Observable, of, ReplaySubject, shareReplay } from 'rxjs';
import { catchError, filter, map, startWith, switchMap, tap } from 'rxjs/operators';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';
import { GroupService } from 'src/app/services/group.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ToastService } from 'src/app/services/toast.service';

import { GroupDialogComponent } from './group-dialog/group-dialog.component';
import { GroupGrantsDialogComponent } from './group-grants-dialog/group-grants-dialog.component';
import { GroupMembersDialogComponent } from './group-members-dialog/group-members-dialog.component';

@Component({
  selector: 'cnsl-groups',
  templateUrl: './groups.component.html',
  styleUrls: ['./groups.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
  standalone: false,
})
export class GroupsComponent {
  protected readonly groups$: Observable<Group[]>;
  protected readonly refresh$ = new ReplaySubject<true>(1);
  protected readonly page$ = new BehaviorSubject<{ pageIndex: number; pageSize: number }>({ pageIndex: 0, pageSize: 20 });
  protected readonly totalResult = signal(0);
  private readonly orgId$: Observable<string>;

  constructor(
    private readonly groupService: GroupService,
    private readonly authService: GrpcAuthService,
    private readonly toast: ToastService,
    private readonly dialog: MatDialog,
    private readonly destroyRef: DestroyRef,
    breadcrumbService: BreadcrumbService,
  ) {
    const bread: Breadcrumb = {
      type: BreadcrumbType.ORG,
      routerLink: ['/org'],
    };
    breadcrumbService.setBreadcrumb([bread]);
    this.orgId$ = this.getOrgId$();
    this.groups$ = this.getGroups$();
  }

  private getOrgId$(): Observable<string> {
    return defer(() => this.authService.getActiveOrg()).pipe(
      switchMap((org) => this.authService.activeOrgChanged.pipe(startWith(org))),
      map((org) => org?.id),
      filter(Boolean),
      shareReplay({ refCount: true, bufferSize: 1 }),
    );
  }

  private getGroups$(): Observable<Group[]> {
    return combineLatest([this.orgId$, this.page$]).pipe(
      switchMap(([orgId, page]) =>
        this.refresh$.pipe(
          startWith(true),
          switchMap(() =>
            this.groupService.listGroups({
              filters: [{ filter: { case: 'organizationId', value: { id: orgId } } }],
              pagination: {
                offset: BigInt(page.pageIndex * page.pageSize),
                limit: page.pageSize,
                asc: false,
              },
              sortingColumn: FieldName.CREATION_DATE,
            }),
          ),
        ),
      ),
      tap((resp) => this.totalResult.set(Number(resp.pagination?.totalResult ?? 0))),
      map(({ groups }) => groups),
      catchError((err) => {
        this.toast.showError(err);
        return of([]);
      }),
    );
  }

  protected changePage(event: PageEvent): void {
    this.page$.next({ pageIndex: event.pageIndex, pageSize: event.pageSize });
  }

  protected async openGroupDialog(group?: Group): Promise<void> {
    const orgId = await firstValueFrom(this.orgId$);

    const request$ = this.dialog
      .open<
        GroupDialogComponent,
        { group?: Group; organizationId: string },
        MessageInitShape<typeof CreateGroupRequestSchema | typeof UpdateGroupRequestSchema>
      >(GroupDialogComponent, {
        width: '450px',
        data: { group, organizationId: orgId },
      })
      .afterClosed()
      .pipe(takeUntilDestroyed(this.destroyRef));

    const request = await lastValueFrom(request$);
    if (!request) {
      return;
    }

    try {
      if (group) {
        await this.groupService.updateGroup(request as MessageInitShape<typeof UpdateGroupRequestSchema>);
        this.toast.showInfo('GROUPS.TOAST.UPDATED', true);
      } else {
        await this.groupService.createGroup(request as MessageInitShape<typeof CreateGroupRequestSchema>);
        this.toast.showInfo('GROUPS.TOAST.CREATED', true);
      }
      await new Promise((res) => setTimeout(res, 1000));
      this.refresh$.next(true);
    } catch (error) {
      this.toast.showError(error);
    }
  }

  protected openMembersDialog(group: Group): void {
    this.dialog
      .open<GroupMembersDialogComponent, { group: Group }, boolean>(GroupMembersDialogComponent, {
        width: '550px',
        data: { group },
      })
      .afterClosed()
      .pipe(takeUntilDestroyed(this.destroyRef))
      .subscribe((changed) => {
        if (changed) {
          this.refresh$.next(true);
        }
      });
  }

  protected openGrantsDialog(group: Group): void {
    this.dialog.open<GroupGrantsDialogComponent, { group: Group }>(GroupGrantsDialogComponent, {
      width: '550px',
      data: { group },
    });
  }

  protected async deleteGroup(group: Group): Promise<void> {
    const confirmed$ = this.dialog
      .open(WarnDialogComponent, {
        width: '400px',
        data: {
          confirmKey: 'ACTIONS.DELETE',
          cancelKey: 'ACTIONS.CANCEL',
          titleKey: 'GROUPS.DIALOG.DELETE.TITLE',
          descriptionKey: 'GROUPS.DIALOG.DELETE.DESCRIPTION',
        },
      })
      .afterClosed()
      .pipe(takeUntilDestroyed(this.destroyRef));

    const confirmed = await lastValueFrom(confirmed$);
    if (!confirmed) {
      return;
    }

    try {
      await this.groupService.deleteGroup({ id: group.id });
      this.toast.showInfo('GROUPS.TOAST.DELETED', true);
      await new Promise((res) => setTimeout(res, 1000));
      this.refresh$.next(true);
    } catch (error) {
      this.toast.showError(error);
    }
  }
}
