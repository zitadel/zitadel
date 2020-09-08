import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { IdpTableComponent } from './idp-table.component';

describe('UserTableComponent', () => {
    let component: IdpTableComponent;
    let fixture: ComponentFixture<IdpTableComponent>;

    beforeEach(async(() => {
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
