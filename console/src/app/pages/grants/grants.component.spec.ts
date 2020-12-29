import { ComponentFixture, TestBed } from '@angular/core/testing';

import { GrantsComponent } from './grants.component';

describe('GrantsComponent', () => {
    let component: GrantsComponent;
    let fixture: ComponentFixture<GrantsComponent>;

    beforeEach(async () => {
        await TestBed.configureTestingModule({
            declarations: [GrantsComponent],
        })
            .compileComponents();
    });

    beforeEach(() => {
        fixture = TestBed.createComponent(GrantsComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
