import { animate, animateChild, query, stagger, style, transition, trigger } from '@angular/animations';
import { SelectionModel } from '@angular/cdk/collections';
import { Component, EventEmitter, Input, OnChanges, Output, SimpleChanges } from '@angular/core';
import { Router } from '@angular/router';
import { ProjectState, ProjectType, ProjectView } from 'src/app/proto/generated/management_pb';
import { AuthService } from 'src/app/services/auth.service';

@Component({
    selector: 'app-owned-project-grid',
    templateUrl: './owned-project-grid.component.html',
    styleUrls: ['./owned-project-grid.component.scss'],
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

    constructor(private router: Router, private authService: AuthService) {
        this.selection.changed.subscribe(selection => {
            this.setPrefixedItem('pinned-projects', JSON.stringify(
                this.selection.selected.map(item => item.projectId),
            )).then(() => {
                const filtered = this.notPinned.filter(item => item === selection.added.find(i => i === item));
                filtered.forEach((f, i) => {
                    this.notPinned.splice(i, 1);
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
                        // this.notPinned.splice(index, 1);
                        return true;
                    }
                });
                this.selection.select(...toSelect);

                const toNotPinned: ProjectView.AsObject[] = this.items.filter((item, index) => {
                    if (!array.includes(item.projectId)) {
                        return true;
                    }
                });
                this.notPinned = toNotPinned;
            }
        });
    }

    private async getPrefixedItem(key: string): Promise<string | null> {
        const prefix = (await this.authService.GetActiveOrg()).id;
        return localStorage.getItem(`${prefix}:${key}`);
    }

    private async setPrefixedItem(key: string, value: any): Promise<void> {
        const prefix = (await this.authService.GetActiveOrg()).id;
        return localStorage.setItem(`${prefix}:${key}`, value);
    }

    public navigateToProject(id: string, event: any): void {
        if (event && event.srcElement && event.srcElement.localName !== 'button') {
            this.router.navigate(['/projects', id]);
        }
    }

    public closeGridView(): void {
        this.changedView.emit(true);
    }
}
