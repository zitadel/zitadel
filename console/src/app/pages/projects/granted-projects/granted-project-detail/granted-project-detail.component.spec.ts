import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { GrantedProjectDetailComponent } from './granted-project-detail.component';

describe('GrantedProjectDetailComponent', () => {
    let component: GrantedProjectDetailComponent;
    let fixture: ComponentFixture<GrantedProjectDetailComponent>;

    beforeEach(waitForAsync(() => {
        TestBed.configureTestingModule({
            declarations: [GrantedProjectDetailComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(GrantedProjectDetailComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
