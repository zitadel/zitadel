import { animate, animateChild, keyframes, query, stagger, style, transition, trigger } from '@angular/animations';
import { SelectionModel } from '@angular/cdk/collections';
import { Component, EventEmitter, Input, OnDestroy, OnInit, Output } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { Router } from '@angular/router';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { BehaviorSubject, Observable, Subject, takeUntil } from 'rxjs';
import { ProjectType } from 'src/app/modules/project-members/project-members-datasource';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import { Org } from 'src/app/proto/generated/zitadel/org_pb';
import { GrantedProject, Project, ProjectState } from 'src/app/proto/generated/zitadel/project_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { StorageKey, StorageLocation, StorageService } from 'src/app/services/storage.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
  selector: 'cnsl-project-grid',
  templateUrl: './project-grid.component.html',
  styleUrls: ['./project-grid.component.scss'],
  animations: [
    trigger('cardAnimation', [
      transition('* => *', [query('@animate', stagger('100ms', animateChild()), { optional: true })]),
    ]),
    trigger('animate', [
      transition(':enter', [
        animate(
          '.2s ease-in',
          keyframes([
            style({ opacity: 0, transform: 'translateY(-50%)', offset: 0 }),
            style({ opacity: 0.5, transform: 'translateY(-10px) scale(1.1)', offset: 0.3 }),
            style({ opacity: 1, transform: 'translateY(0)', offset: 1 }),
          ]),
        ),
      ]),
      transition(':leave', [
        animate(
          '.2s ease-out',
          keyframes([
            style({ opacity: 1, transform: 'scale(1.1)', offset: 0 }),
            style({ opacity: 0.5, transform: 'scale(.5)', offset: 0.3 }),
            style({ opacity: 0, transform: 'scale(0)', offset: 1 }),
          ]),
        ),
      ]),
    ]),
  ],
})
export class ProjectGridComponent implements OnInit, OnDestroy {
  public projectList: Array<Project.AsObject | GrantedProject.AsObject> = [];
  public totalResult: number = 0;
  public viewTimestamp!: Timestamp.AsObject;

  @Input() public projectType$: BehaviorSubject<any> = new BehaviorSubject(ProjectType.PROJECTTYPE_OWNED);
  @Output() public emitAddProject: EventEmitter<void> = new EventEmitter();

  public notPinned: Array<Project.AsObject | GrantedProject.AsObject> = [];

  private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public loading$: Observable<boolean> = this.loadingSubject.asObservable();

  public selection: SelectionModel<Project.AsObject | GrantedProject.AsObject> = new SelectionModel<
    Project.AsObject | GrantedProject.AsObject
  >(true, []);

  public ProjectState: any = ProjectState;
  public ProjectType: any = ProjectType;
  @Input() public zitadelProjectId: string = '';
  private destroy$: Subject<void> = new Subject();

  constructor(
    private router: Router,
    private dialog: MatDialog,
    private storage: StorageService,
    private mgmtService: ManagementService,
    private toast: ToastService,
  ) {}

  public listenForSelectionChanges(): void {
    this.selection.changed.pipe(takeUntil(this.destroy$)).subscribe((selection) => {
      if (this.projectType$.value === ProjectType.PROJECTTYPE_OWNED) {
        this.setPrefixedItem(
          'pinned-projects',
          JSON.stringify(this.selection.selected.map((item) => (item as Project.AsObject).id)),
        ).then(() => {
          selection.added.forEach((item) => {
            const index = (this.notPinned as Array<Project.AsObject>).findIndex(
              (i) => i.id === (item as Project.AsObject).id,
            );
            this.notPinned.splice(index, 1);
          });
          this.notPinned.push(...(selection.removed as Project.AsObject[]));
        });
      } else if (this.projectType$.value === ProjectType.PROJECTTYPE_GRANTED) {
        this.setPrefixedItem(
          'pinned-granted-projects',
          JSON.stringify(this.selection.selected.map((item) => (item as GrantedProject.AsObject).projectId)),
        ).then(() => {
          selection.added.forEach((item) => {
            const index = (this.notPinned as Array<GrantedProject.AsObject>).findIndex(
              (i) => i.projectId === (item as GrantedProject.AsObject).projectId,
            );
            this.notPinned.splice(index, 1);
          });
          this.notPinned.push(...(selection.removed as GrantedProject.AsObject[]));
        });
      }
    });
  }

  public ngOnInit(): void {
    this.projectType$.pipe(takeUntil(this.destroy$)).subscribe((type) => {
      this.getData(type).then(() => {
        this.notPinned = Object.assign([], this.projectList);
        this.reorganizeItems(type);
      });
    });
  }

  public ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  private async getData(type: ProjectType, limit?: number, offset?: number): Promise<void> {
    this.loadingSubject.next(true);
    switch (type) {
      case ProjectType.PROJECTTYPE_OWNED:
        return this.mgmtService
          .listProjects(limit, offset)
          .then((resp) => {
            this.projectList = resp.resultList;
            if (resp.details?.totalResult) {
              this.totalResult = resp.details.totalResult;
            } else {
              this.totalResult = 0;
            }
            if (this.totalResult > 10) {
              // trigger change to table
            }
            if (resp.details?.viewTimestamp) {
              this.viewTimestamp = resp.details?.viewTimestamp;
            }

            this.loadingSubject.next(false);
          })
          .catch((error) => {
            this.toast.showError(error);
            this.loadingSubject.next(false);
          });

      case ProjectType.PROJECTTYPE_GRANTED:
        return this.mgmtService
          .listGrantedProjects(limit, offset)
          .then((resp) => {
            this.projectList = resp.resultList;
            if (resp.details?.totalResult) {
              this.totalResult = resp.details.totalResult;
            } else {
              this.totalResult = 0;
            }
            if (this.totalResult > 10) {
              // trigger change to table
            }
            if (resp.details?.viewTimestamp) {
              this.viewTimestamp = resp.details?.viewTimestamp;
            }

            this.loadingSubject.next(false);
          })
          .catch((error) => {
            console.error(error);
            this.toast.showError(error);
            this.loadingSubject.next(false);
          });
    }
  }

  public addItem(): void {
    this.emitAddProject.emit();
  }

  public reorganizeItems(type: ProjectType): void {
    this.selection = new SelectionModel<Project.AsObject | GrantedProject.AsObject>(true, []);
    this.listenForSelectionChanges();

    switch (type) {
      case ProjectType.PROJECTTYPE_OWNED:
        this.getPrefixedItem('pinned-projects').then((storageEntry) => {
          if (storageEntry) {
            const array: string[] = JSON.parse(storageEntry);
            const toSelect: Project.AsObject[] = (this.projectList as Project.AsObject[]).filter((item) => {
              if (array.includes(item.id)) {
                return true;
              } else {
                return false;
              }
            });

            this.selection.select(...toSelect);
          }
        });
        break;
      case ProjectType.PROJECTTYPE_GRANTED:
        this.getPrefixedItem('pinned-granted-projects').then((storageEntry) => {
          if (storageEntry) {
            const array: string[] = JSON.parse(storageEntry);
            const toSelect: GrantedProject.AsObject[] = (this.projectList as GrantedProject.AsObject[]).filter((item) => {
              if (array.includes(item.projectId)) {
                return true;
              } else {
                return false;
              }
            });

            this.selection.select(...toSelect);
          }
        });
        break;
    }
  }

  private async getPrefixedItem(key: string): Promise<string | null> {
    const org = this.storage.getItem<Org.AsObject>(StorageKey.organization, StorageLocation.session) as Org.AsObject;
    return localStorage.getItem(`${org?.id}:${key}`);
  }

  private async setPrefixedItem(key: string, value: any): Promise<void> {
    const org = this.storage.getItem<Org.AsObject>(StorageKey.organization, StorageLocation.session) as Org.AsObject;
    return localStorage.setItem(`${org.id}:${key}`, value);
  }

  public navigateToProject(type: ProjectType, item: Project.AsObject | GrantedProject.AsObject, event: any): void {
    if (event && event.srcElement && event.srcElement.localName !== 'button') {
      if (type === ProjectType.PROJECTTYPE_OWNED) {
        this.router.navigate(['/projects', (item as Project.AsObject).id]);
      } else if (type === ProjectType.PROJECTTYPE_GRANTED) {
        this.router.navigate([
          '/granted-projects',
          (item as GrantedProject.AsObject).projectId,
          'grant',
          (item as GrantedProject.AsObject).grantId,
        ]);
      }
    }
  }

  public toggle(item: Project.AsObject, event: any): void {
    event.stopPropagation();
    this.selection.toggle(item);
  }

  public deleteProject(event: any, item: Project.AsObject): void {
    event.stopPropagation();
    const dialogRef = this.dialog.open(WarnDialogComponent, {
      data: {
        confirmKey: 'ACTIONS.DELETE',
        cancelKey: 'ACTIONS.CANCEL',
        titleKey: 'PROJECT.PAGES.DIALOG.DELETE.TITLE',
        descriptionKey: 'PROJECT.PAGES.DIALOG.DELETE.DESCRIPTION',
        confirmationKey: 'PROJECT.PAGES.DIALOG.DELETE.TYPENAME',
        confirmation: item.name,
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((resp) => {
      if (resp && item.id !== this.zitadelProjectId) {
        this.mgmtService
          .removeProject(item.id)
          .then(() => {
            this.toast.showInfo('PROJECT.TOAST.DELETED', true);
            const index = this.projectList.findIndex((iter) => (iter as Project.AsObject).id === item.id);
            if (index > -1) {
              this.projectList.splice(index, 1);
            }

            const indexSelection = this.selection.selected.findIndex((iter) => (iter as Project.AsObject).id === item.id);
            if (indexSelection > -1) {
              this.selection.selected.splice(indexSelection, 1);
            }

            const indexPinned = (this.notPinned as Array<Project.AsObject>).findIndex((iter) => iter.id === item.id);
            if (indexPinned > -1) {
              this.notPinned.splice(indexPinned, 1);
            }
          })
          .catch((error) => {
            this.toast.showError(error);
          });
      }
    });
  }
}
