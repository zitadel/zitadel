import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { RefreshTableComponent } from './refresh-table.component';

describe('RefreshTableComponent', () => {
    let component: RefreshTableComponent;
    let fixture: ComponentFixture<RefreshTableComponent>;

    beforeEach(waitForAsync(() => {
        TestBed.configureTestingModule({
            declarations: [RefreshTableComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(RefreshTableComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
