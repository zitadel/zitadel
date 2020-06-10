import { animate, animateChild, query, stagger, style, transition, trigger } from '@angular/animations';
import { SelectionModel } from '@angular/cdk/collections';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { Router } from '@angular/router';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { GrantedProject, Project, ProjectState } from 'src/app/proto/generated/management_pb';
import { ProjectService } from 'src/app/services/project.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
    selector: 'app-project-grid',
    templateUrl: './project-grid.component.html',
    styleUrls: ['./project-grid.component.scss'],
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
export class ProjectGridComponent {
    @Input() items: Array<GrantedProject.AsObject> = [];
    @Output() newClicked: EventEmitter<boolean> = new EventEmitter();
    @Output() changedView: EventEmitter<boolean> = new EventEmitter();
    @Input() loading: boolean = false;

    public selection: SelectionModel<GrantedProject.AsObject> = new SelectionModel<GrantedProject.AsObject>(true, []);
    public selectedIndex: number = -1;

    public showNewProject: boolean = false;
    public ProjectState: any = ProjectState;

    constructor(private router: Router, private projectService: ProjectService, private toast: ToastService) { }

    public selectItem(item: GrantedProject.AsObject, event?: any): void {
        if (event && !event.target.classList.contains('mat-icon')) {
            if (item.grantId) {
                this.router.navigate(['projects', item.id, 'grant', `${item.grantId}`]);
            } else {
                this.router.navigate(['/projects', item.id]);
            }
        } else if (!event) {
            if (item.grantId) {
                this.router.navigate(['projects', item.id, 'grant', `${item.grantId}`]);
            } else {
                this.router.navigate(['/projects', item.id]);
            }
        }
    }

    public addItem(): void {
        this.newClicked.emit(true);
    }

    public dateFromTimestamp(date: Timestamp.AsObject): any {
        const ts: Date = new Date(date.seconds * 1000 + date.nanos / 1000);
        return ts;
    }

    public reactivateProjects(selected: Project.AsObject[]): void {
        Promise.all([selected.map(proj => {
            return this.projectService.ReactivateProject(proj.id);
        })]).then(() => {
            this.toast.showInfo('Successful reactivated all projects');
        }).catch(error => {
            this.toast.showError(error.message);
        });
    }

    public deactivateProjects(selected: Project.AsObject[]): void {
        Promise.all([selected.map(proj => {
            return this.projectService.DeactivateProject(proj.id);
        })]).then(() => {
            this.toast.showInfo('Successful deactivated all projects');
        }).catch(error => {
            this.toast.showError(error.message);
        });
    }
}
