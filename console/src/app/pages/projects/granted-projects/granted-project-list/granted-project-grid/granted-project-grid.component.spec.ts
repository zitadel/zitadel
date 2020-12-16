import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { GrantedProjectGridComponent } from './granted-project-grid.component';

describe('GridComponent', () => {
    let component: GrantedProjectGridComponent;
    let fixture: ComponentFixture<GrantedProjectGridComponent>;

    beforeEach(waitForAsync(() => {
        TestBed.configureTestingModule({
            declarations: [GrantedProjectGridComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(GrantedProjectGridComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
