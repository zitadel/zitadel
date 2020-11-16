import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { IdpCreateComponent } from './idp-create.component';

describe('IdpCreateComponent', () => {
    let component: IdpCreateComponent;
    let fixture: ComponentFixture<IdpCreateComponent>;

    beforeEach(waitForAsync(() => {
        TestBed.configureTestingModule({
            declarations: [IdpCreateComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(IdpCreateComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
