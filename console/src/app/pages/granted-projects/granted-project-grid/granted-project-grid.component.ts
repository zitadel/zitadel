import { animate, animateChild, query, stagger, style, transition, trigger } from '@angular/animations';
import { SelectionModel } from '@angular/cdk/collections';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { Router } from '@angular/router';
import { ProjectGrantView, ProjectState, ProjectType } from 'src/app/proto/generated/management_pb';

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
export class GrantedProjectGridComponent {
    @Input() items: Array<ProjectGrantView.AsObject> = [];
    @Output() newClicked: EventEmitter<boolean> = new EventEmitter();
    @Output() changedView: EventEmitter<boolean> = new EventEmitter();
    @Input() loading: boolean = false;
    public selection: SelectionModel<string> = new SelectionModel<string>(true, []);


    public showNewProject: boolean = false;
    public ProjectState: any = ProjectState;
    public ProjectType: any = ProjectType;

    constructor(private router: Router) {
        const storageEntry = localStorage.getItem('pinned-granted-projects');
        if (storageEntry) {
            const array = JSON.parse(storageEntry);
            this.selection.select(...array);
        }

        this.selection.changed.subscribe(selection => {
            console.log(this.selection.selected);
            localStorage.setItem('pinned-granted-projects', JSON.stringify(this.selection.selected));
        });
    }

    public selectItem(item: ProjectGrantView.AsObject, event?: any): void {
        if (event && !event.target.classList.contains('mat-icon')) {
            this.router.navigate(['granted-projects', item.projectId, 'grant', `${item.id}`]);
        } else if (!event) {
            this.router.navigate(['granted-projects', item.projectId, 'grant', `${item.id}`]);
        }
    }

    public addItem(): void {
        this.newClicked.emit(true);
    }
}
