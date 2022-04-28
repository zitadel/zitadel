import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { DomainVerificationComponent } from './domain-verification.component';

describe('DomainVerificationComponent', () => {
    let component: DomainVerificationComponent;
    let fixture: ComponentFixture<DomainVerificationComponent>;

    beforeEach(waitForAsync(() => {
        TestBed.configureTestingModule({
            declarations: [DomainVerificationComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(DomainVerificationComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
