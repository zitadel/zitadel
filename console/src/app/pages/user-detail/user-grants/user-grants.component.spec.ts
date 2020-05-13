import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { UserGrantsComponent } from './user-grants.component';

describe('UserGrantsComponent', () => {
    let component: UserGrantsComponent;
    let fixture: ComponentFixture<UserGrantsComponent>;

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            declarations: [UserGrantsComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(UserGrantsComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
