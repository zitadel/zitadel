import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { DomainVerificationComponent } from './domain-verification.component';

describe('DomainVerificationComponent', () => {
    let component: DomainVerificationComponent;
    let fixture: ComponentFixture<DomainVerificationComponent>;

    beforeEach(async(() => {
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
