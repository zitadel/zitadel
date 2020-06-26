import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { GrantedProjectsComponent } from './granted-projects.component';

describe('GrantedProjectsComponent', () => {
    let component: GrantedProjectsComponent;
    let fixture: ComponentFixture<GrantedProjectsComponent>;

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            declarations: [GrantedProjectsComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(GrantedProjectsComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
