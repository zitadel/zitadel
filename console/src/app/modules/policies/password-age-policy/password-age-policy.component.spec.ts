import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { PasswordAgePolicyComponent } from './password-age-policy.component';

describe('PasswordAgePolicyComponent', () => {
    let component: PasswordAgePolicyComponent;
    let fixture: ComponentFixture<PasswordAgePolicyComponent>;

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            declarations: [PasswordAgePolicyComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(PasswordAgePolicyComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
