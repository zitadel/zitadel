import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { UserMfaComponent } from './user-mfa.component';

describe('UserMfaComponent', () => {
    let component: UserMfaComponent;
    let fixture: ComponentFixture<UserMfaComponent>;

    beforeEach(async(() => {
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
