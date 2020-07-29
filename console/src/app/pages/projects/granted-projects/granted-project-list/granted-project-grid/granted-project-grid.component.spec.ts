import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { GrantedProjectGridComponent } from './granted-project-grid.component';

describe('GridComponent', () => {
    let component: GrantedProjectGridComponent;
    let fixture: ComponentFixture<GrantedProjectGridComponent>;

    beforeEach(async(() => {
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
