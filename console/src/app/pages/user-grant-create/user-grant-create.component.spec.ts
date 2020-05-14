import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { UserGrantCreateComponent } from './user-grant-create.component';

describe('UserGrantCreateComponent', () => {
    let component: UserGrantCreateComponent;
    let fixture: ComponentFixture<UserGrantCreateComponent>;

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            declarations: [UserGrantCreateComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(UserGrantCreateComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
