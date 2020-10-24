import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { PolicyGridComponent } from './policy-grid.component';

describe('PolicyGridComponent', () => {
    let component: PolicyGridComponent;
    let fixture: ComponentFixture<PolicyGridComponent>;

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            declarations: [PolicyGridComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(PolicyGridComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
