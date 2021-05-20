import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AdditionalOriginsComponent } from './additional-origins.component';

describe('AdditionalOriginsComponent', () => {
    let component: AdditionalOriginsComponent;
    let fixture: ComponentFixture<AdditionalOriginsComponent>;

    beforeEach(async () => {
        await TestBed.configureTestingModule({
            declarations: [AdditionalOriginsComponent],
        })
            .compileComponents();
    });

    beforeEach(() => {
        fixture = TestBed.createComponent(AdditionalOriginsComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
