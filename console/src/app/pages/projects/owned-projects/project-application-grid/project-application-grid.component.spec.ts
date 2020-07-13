import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ProjectApplicationGridComponent } from './project-application-grid.component';

describe('AppGridComponent', () => {
    let component: ProjectApplicationGridComponent;
    let fixture: ComponentFixture<ProjectApplicationGridComponent>;

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            declarations: [ProjectApplicationGridComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(ProjectApplicationGridComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
