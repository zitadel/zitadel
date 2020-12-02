import { ComponentFixture, TestBed } from '@angular/core/testing';

import { InfoSectionComponent } from './info-section.component';

describe('InfoSectionComponent', () => {
    let component: InfoSectionComponent;
    let fixture: ComponentFixture<InfoSectionComponent>;

    beforeEach(async () => {
        await TestBed.configureTestingModule({
            declarations: [InfoSectionComponent],
        })
            .compileComponents();
    });

    beforeEach(() => {
        fixture = TestBed.createComponent(InfoSectionComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
