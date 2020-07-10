import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { GrantedProjectDetailComponent } from './granted-project-detail.component';

describe('GrantedProjectDetailComponent', () => {
    let component: GrantedProjectDetailComponent;
    let fixture: ComponentFixture<GrantedProjectDetailComponent>;

    beforeEach(async(() => {
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
