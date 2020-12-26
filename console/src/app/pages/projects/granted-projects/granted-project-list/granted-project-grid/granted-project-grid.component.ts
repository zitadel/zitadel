import { animate, animateChild, query, stagger, style, transition, trigger } from '@angular/animations';
import { SelectionModel } from '@angular/cdk/collections';
import { Component, EventEmitter, Input, OnChanges, Output, SimpleChanges } from '@angular/core';
import { Router } from '@angular/router';
import { Org } from 'src/app/proto/generated/zitadel/auth_pb';
import { ProjectGrantView, ProjectState, ProjectType } from 'src/app/proto/generated/zitadel/management_pb';
import { StorageKey, StorageService } from 'src/app/services/storage.service';

@Component({
    selector: 'app-granted-project-grid',
    templateUrl: './granted-project-grid.component.html',
    styleUrls: ['./granted-project-grid.component.scss'],
    animations: [
        trigger('list', [
            transition(':enter', [
                query('@animate',
                    stagger(100, animateChild()),
                ),
            ]),
        ]),
        trigger('animate', [
            transition(':enter', [
                style({ opacity: 0, transform: 'translateY(-100%)' }),
                animate('100ms', style({ opacity: 1, transform: 'translateY(0)' })),
            ]),
            transition(':leave', [
                style({ opacity: 1, transform: 'translateY(0)' }),
                animate('100ms', style({ opacity: 0, transform: 'translateY(100%)' })),
            ]),
        ]),
    ],
})
export class GrantedProjectGridComponent implements OnChanges {
    @Input() items: Array<ProjectGrantView.AsObject> = [];
    public notPinned: Array<ProjectGrantView.AsObject> = [];
    @Output() newClicked: EventEmitter<boolean> = new EventEmitter();
    @Output() changedView: EventEmitter<boolean> = new EventEmitter();
    @Input() loading: boolean = false;
    public selection: SelectionModel<ProjectGrantView.AsObject> = new SelectionModel<ProjectGrantView.AsObject>(true, []);

    public showNewProject: boolean = false;
    public ProjectState: any = ProjectState;
    public ProjectType: any = ProjectType;

    constructor(private storage: StorageService, private router: Router) {
        this.selection.changed.subscribe(selection => {
            this.setPrefixedItem('pinned-granted-projects', JSON.stringify(
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

    public addItem(): void {
        this.newClicked.emit(true);
    }

    public ngOnChanges(changes: SimpleChanges): void {
        if (changes.items?.currentValue && changes.items.currentValue.length > 0) {
            this.notPinned = Object.assign([], this.items);
            this.reorganizeItems();
        }
    }

    public reorganizeItems(): void {
        this.getPrefixedItem('pinned-granted-projects').then(storageEntry => {
            if (storageEntry) {
                const array: string[] = JSON.parse(storageEntry);
                const toSelect: ProjectGrantView.AsObject[] = this.items.filter((item, index) => {
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

    public navigateToProject(projectId: string, id: string, event: any): void {
        if (event && event.srcElement && event.srcElement.localName !== 'button') {
            this.router.navigate(['/granted-projects', projectId, 'grant', id]);
        }
    }
}
