import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { IdpCreateComponent } from './idp-create.component';

describe('IdpCreateComponent', () => {
    let component: IdpCreateComponent;
    let fixture: ComponentFixture<IdpCreateComponent>;

    beforeEach(async(() => {
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
