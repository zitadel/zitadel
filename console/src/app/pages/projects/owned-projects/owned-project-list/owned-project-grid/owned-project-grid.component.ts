import { animate, animateChild, keyframes, query, stagger, style, transition, trigger } from '@angular/animations';
import { SelectionModel } from '@angular/cdk/collections';
import { Component, EventEmitter, Input, OnChanges, Output, SimpleChanges } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { Router } from '@angular/router';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import { Org } from 'src/app/proto/generated/zitadel/auth_pb';
import { ProjectState, ProjectType, ProjectView } from 'src/app/proto/generated/zitadel/management_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { StorageKey, StorageService } from 'src/app/services/storage.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
    selector: 'app-owned-project-grid',
    templateUrl: './owned-project-grid.component.html',
    styleUrls: ['./owned-project-grid.component.scss'],
    animations: [
        trigger('cardAnimation', [
            transition('* => *', [
                query('@animate', stagger('100ms', animateChild()), { optional: true }),
            ]),
        ]),
        trigger('animate', [
            transition(':enter', [
                animate('.2s ease-in', keyframes([
                    style({ opacity: 0, transform: 'translateY(-50%)', offset: 0 }),
                    style({ opacity: .5, transform: 'translateY(-10px) scale(1.1)', offset: 0.3 }),
                    style({ opacity: 1, transform: 'translateY(0)', offset: 1 }),
                ])),
            ]),
            transition(':leave', [
                animate('.2s ease-out', keyframes([
                    style({ opacity: 1, transform: 'scale(1.1)', offset: 0 }),
                    style({ opacity: .5, transform: 'scale(.5)', offset: 0.3 }),
                    style({ opacity: 0, transform: 'scale(0)', offset: 1 }),
                ])),
            ]),
        ]),
    ],
})
export class OwnedProjectGridComponent implements OnChanges {
    @Input() items: Array<ProjectView.AsObject> = [];
    public notPinned: Array<ProjectView.AsObject> = [];

    @Output() newClicked: EventEmitter<boolean> = new EventEmitter();
    @Output() changedView: EventEmitter<boolean> = new EventEmitter();
    @Input() loading: boolean = false;

    public selection: SelectionModel<ProjectView.AsObject> = new SelectionModel<ProjectView.AsObject>(true, []);

    public showNewProject: boolean = false;
    public ProjectState: any = ProjectState;
    public ProjectType: any = ProjectType;
    @Input() public zitadelProjectId: string = '';
    constructor(
        private router: Router,
        private dialog: MatDialog,
        private storage: StorageService,
        private mgmtService: ManagementService,
        private toast: ToastService,
    ) {
        this.selection.changed.subscribe(selection => {
            this.setPrefixedItem('pinned-projects', JSON.stringify(
                this.selection.selected.map(item => item.projectId),
            )).then(() => {
                selection.added.forEach(item => {
                    const index = this.notPinned.findIndex(i => i.projectId === item.projectId);
                    this.notPinned.splice(index, 1);
                });
                this.notPinned.push(...selection.removed);
            });
        });
    }

    public selectItem(item: ProjectView.AsObject, event?: any): void {
        if (event && !event.target.classList.contains('mat-icon')) {
            this.router.navigate(['/projects', item.projectId]);
        } else if (!event) {
            this.router.navigate(['/projects', item.projectId]);
        }
    }

    public addItem(): void {
        this.newClicked.emit(true);
    }

    public ngOnChanges(changes: SimpleChanges): void {
        if (changes.items.currentValue && changes.items.currentValue.length > 0) {
            this.notPinned = Object.assign([], this.items);
            this.reorganizeItems();
        }
    }

    public reorganizeItems(): void {
        this.getPrefixedItem('pinned-projects').then(storageEntry => {
            if (storageEntry) {
                const array: string[] = JSON.parse(storageEntry);
                const toSelect: ProjectView.AsObject[] = this.items.filter((item, index) => {
                    if (array.includes(item.projectId)) {
                        return true;
                    }
                });
                this.selection.select(...toSelect);
            }
        });
    }

    private async getPrefixedItem(key: string): Promise<string | null> {
        const org = this.storage.getItem<Org.AsObject>(StorageKey.organization) as Org.AsObject;
        return localStorage.getItem(`${org.id}:${key}`);
    }

    private async setPrefixedItem(key: string, value: any): Promise<void> {
        const org = this.storage.getItem<Org.AsObject>(StorageKey.organization) as Org.AsObject;
        return localStorage.setItem(`${org.id}:${key}`, value);
    }

    public navigateToProject(id: string, event: any): void {
        if (event && event.srcElement && event.srcElement.localName !== 'button') {
            this.router.navigate(['/projects', id]);
        }
    }

    public closeGridView(): void {
        this.changedView.emit(true);
    }

    public toggle(item: ProjectView.AsObject, event: any): void {
        event.stopPropagation();
        this.selection.toggle(item);
    }

    public deleteProject(event: any, item: ProjectView.AsObject): void {
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

        dialogRef.afterClosed().subscribe(resp => {
            if (resp && item.projectId !== this.zitadelProjectId) {
                this.mgmtService.RemoveProject(item.projectId).then(() => {
                    this.toast.showInfo('PROJECT.TOAST.DELETED', true);
                    const index = this.items.findIndex(iter => iter.projectId === item.projectId);
                    if (index > -1) {
                        this.items.splice(index, 1);
                    }

                    const indexSelection = this.selection.selected.findIndex(iter => iter.projectId === item.projectId);
                    if (indexSelection > -1) {
                        this.selection.selected.splice(indexSelection, 1);
                    }

                    const indexPinned = this.notPinned.findIndex(iter => iter.projectId === item.projectId);
                    if (indexPinned > -1) {
                        this.notPinned.splice(indexPinned, 1);
                    }
                }).catch(error => {
                    this.toast.showError(error);
                });
            }
        });
    }
}
