import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { IdpComponent } from './idp.component';

describe('IdComponent', () => {
    let component: IdpComponent;
    let fixture: ComponentFixture<IdpComponent>;

    beforeEach(waitForAsync(() => {
        TestBed.configureTestingModule({
            declarations: [IdpComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(IdpComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
