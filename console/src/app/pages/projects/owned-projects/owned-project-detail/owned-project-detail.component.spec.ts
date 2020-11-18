import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { OwnedProjectDetailComponent } from './owned-project-detail.component';


describe('ProjectDetailComponent', () => {
    let component: OwnedProjectDetailComponent;
    let fixture: ComponentFixture<OwnedProjectDetailComponent>;

    beforeEach(waitForAsync(() => {
        TestBed.configureTestingModule({
            declarations: [OwnedProjectDetailComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(OwnedProjectDetailComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
