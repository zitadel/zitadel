import { animate, animateChild, keyframes, query, stagger, style, transition, trigger } from '@angular/animations';
import { SelectionModel } from '@angular/cdk/collections';
import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { Router } from '@angular/router';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { BehaviorSubject, Observable } from 'rxjs';
import { ProjectType } from 'src/app/modules/project-members/project-members-datasource';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import { Org } from 'src/app/proto/generated/zitadel/org_pb';
import { Project, ProjectState } from 'src/app/proto/generated/zitadel/project_pb';
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
export class ProjectGridComponent implements OnInit {
  public ownedProjectList: Array<Project.AsObject> = [];
  public totalResult: number = 0;
  public viewTimestamp!: Timestamp.AsObject;

  @Input() public projectType: ProjectType = ProjectType.PROJECTTYPE_OWNED;
  @Output() public emitAddProject: EventEmitter<void> = new EventEmitter();

  public notPinned: Array<Project.AsObject> = [];

  private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public loading$: Observable<boolean> = this.loadingSubject.asObservable();

  public selection: SelectionModel<Project.AsObject> = new SelectionModel<Project.AsObject>(true, []);

  public ProjectState: any = ProjectState;
  @Input() public zitadelProjectId: string = '';

  constructor(
    private router: Router,
    private dialog: MatDialog,
    private storage: StorageService,
    private mgmtService: ManagementService,
    private toast: ToastService,
  ) {
    this.selection.changed.subscribe((selection) => {
      this.setPrefixedItem('pinned-projects', JSON.stringify(this.selection.selected.map((item) => item.id))).then(() => {
        selection.added.forEach((item) => {
          const index = this.notPinned.findIndex((i) => i.id === item.id);
          this.notPinned.splice(index, 1);
        });
        this.notPinned.push(...selection.removed);
      });
    });
  }

  public ngOnInit(): void {
    this.getData().then(() => {
      console.log(this.ownedProjectList);
      this.notPinned = Object.assign([], this.ownedProjectList);
      this.reorganizeItems();
    });
  }

  private async getData(limit?: number, offset?: number): Promise<void> {
    this.loadingSubject.next(true);
    return this.mgmtService
      .listProjects(limit, offset)
      .then((resp) => {
        this.ownedProjectList = resp.resultList;
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

  public selectItem(item: Project.AsObject, event?: any): void {
    if (event && !event.target.classList.contains('mat-icon')) {
      this.router.navigate(['/projects', item.id]);
    } else if (!event) {
      this.router.navigate(['/projects', item.id]);
    }
  }

  public addItem(): void {
    this.emitAddProject.emit();
  }

  public reorganizeItems(): void {
    this.getPrefixedItem('pinned-projects').then((storageEntry) => {
      if (storageEntry) {
        const array: string[] = JSON.parse(storageEntry);
        const toSelect: Project.AsObject[] = this.ownedProjectList.filter((item) => {
          if (array.includes(item.id)) {
            return true;
          } else {
            return false;
          }
        });
        this.selection.select(...toSelect);
      }
    });
  }

  private async getPrefixedItem(key: string): Promise<string | null> {
    const org = this.storage.getItem<Org.AsObject>(StorageKey.organization, StorageLocation.session) as Org.AsObject;
    return localStorage.getItem(`${org?.id}:${key}`);
  }

  private async setPrefixedItem(key: string, value: any): Promise<void> {
    const org = this.storage.getItem<Org.AsObject>(StorageKey.organization, StorageLocation.session) as Org.AsObject;
    return localStorage.setItem(`${org.id}:${key}`, value);
  }

  public navigateToProject(id: string, event: any): void {
    if (event && event.srcElement && event.srcElement.localName !== 'button') {
      this.router.navigate(['/projects', id]);
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
      },
      width: '400px',
    });

    dialogRef.afterClosed().subscribe((resp) => {
      if (resp && item.id !== this.zitadelProjectId) {
        this.mgmtService
          .removeProject(item.id)
          .then(() => {
            this.toast.showInfo('PROJECT.TOAST.DELETED', true);
            const index = this.ownedProjectList.findIndex((iter) => iter.id === item.id);
            if (index > -1) {
              this.ownedProjectList.splice(index, 1);
            }

            const indexSelection = this.selection.selected.findIndex((iter) => iter.id === item.id);
            if (indexSelection > -1) {
              this.selection.selected.splice(indexSelection, 1);
            }

            const indexPinned = this.notPinned.findIndex((iter) => iter.id === item.id);
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
