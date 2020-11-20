import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { IdpTableComponent } from './idp-table.component';

describe('UserTableComponent', () => {
    let component: IdpTableComponent;
    let fixture: ComponentFixture<IdpTableComponent>;

    beforeEach(waitForAsync(() => {
        TestBed.configureTestingModule({
            declarations: [IdpTableComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(IdpTableComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
