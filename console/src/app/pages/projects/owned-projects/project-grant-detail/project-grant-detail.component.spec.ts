import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { ProjectGrantDetailComponent } from './project-grant-detail.component';

describe('GrantComponent', () => {
    let component: ProjectGrantDetailComponent;
    let fixture: ComponentFixture<ProjectGrantDetailComponent>;

    beforeEach(waitForAsync(() => {
        TestBed.configureTestingModule({
            declarations: [ProjectGrantDetailComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(ProjectGrantDetailComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
