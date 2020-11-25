import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { ResendEmailDialogComponent } from './resend-email-dialog.component';

describe('ResendEmailDialogComponent', () => {
    let component: ResendEmailDialogComponent;
    let fixture: ComponentFixture<ResendEmailDialogComponent>;

    beforeEach(waitForAsync(() => {
        TestBed.configureTestingModule({
            declarations: [ResendEmailDialogComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(ResendEmailDialogComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
