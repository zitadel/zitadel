import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { MfaTableComponent } from './mfa-table.component';

describe('MfaTableComponent', () => {
    let component: MfaTableComponent;
    let fixture: ComponentFixture<MfaTableComponent>;

    beforeEach(waitForAsync(() => {
        TestBed.configureTestingModule({
            declarations: [MfaTableComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(MfaTableComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
