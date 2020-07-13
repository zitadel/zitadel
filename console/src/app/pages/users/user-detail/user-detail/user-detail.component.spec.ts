import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { UserDetailComponent } from './user-detail.component';

describe('UserDetailComponent', () => {
    let component: UserDetailComponent;
    let fixture: ComponentFixture<UserDetailComponent>;

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            declarations: [UserDetailComponent],
        }).compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(UserDetailComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
