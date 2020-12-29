import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { UserMfaComponent } from './user-mfa.component';

describe('UserMfaComponent', () => {
    let component: UserMfaComponent;
    let fixture: ComponentFixture<UserMfaComponent>;

    beforeEach(waitForAsync(() => {
        TestBed.configureTestingModule({
            declarations: [UserMfaComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(UserMfaComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
