import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { PasswordComplexityViewComponent } from './password-complexity-view.component';

describe('PasswordComplexityViewComponent', () => {
    let component: PasswordComplexityViewComponent;
    let fixture: ComponentFixture<PasswordComplexityViewComponent>;

    beforeEach(async(() => {
        TestBed.configureTestingModule({
            declarations: [PasswordComplexityViewComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(PasswordComplexityViewComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
