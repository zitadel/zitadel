import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { AuthFactorDialogComponent } from './auth-factor-dialog.component';

describe('CodeDialogComponent', () => {
    let component: AuthFactorDialogComponent;
    let fixture: ComponentFixture<AuthFactorDialogComponent>;

    beforeEach(waitForAsync(() => {
        TestBed.configureTestingModule({
            declarations: [AuthFactorDialogComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(AuthFactorDialogComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
