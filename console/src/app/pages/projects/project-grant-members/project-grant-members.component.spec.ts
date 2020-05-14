import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { MatPaginatorModule } from '@angular/material/paginator';
import { MatSortModule } from '@angular/material/sort';
import { MatTableModule } from '@angular/material/table';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';

import { ProjectGrantMembersComponent } from './project-grant-members.component';

describe('ProjectMembersComponent', () => {
    let component: ProjectGrantMembersComponent;
    let fixture: ComponentFixture<ProjectGrantMembersComponent>;

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            declarations: [ProjectGrantMembersComponent],
            imports: [
                NoopAnimationsModule,
                MatPaginatorModule,
                MatSortModule,
                MatTableModule,
            ],
        }).compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(ProjectGrantMembersComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should compile', () => {
        expect(component).toBeTruthy();
    });
});
