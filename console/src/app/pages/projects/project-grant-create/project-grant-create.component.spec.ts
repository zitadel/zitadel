import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ProjectGrantCreateComponent } from './project-grant-create.component';

describe('GrantCreateComponent', () => {
    let component: ProjectGrantCreateComponent;
    let fixture: ComponentFixture<ProjectGrantCreateComponent>;

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            declarations: [ProjectGrantCreateComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(ProjectGrantCreateComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
