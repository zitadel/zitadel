import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { AuthUserMfaComponent } from './auth-user-mfa.component';

describe('AuthUserMfaComponent', () => {
    let component: AuthUserMfaComponent;
    let fixture: ComponentFixture<AuthUserMfaComponent>;

    beforeEach(waitForAsync(() => {
        TestBed.configureTestingModule({
            declarations: [AuthUserMfaComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(AuthUserMfaComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
