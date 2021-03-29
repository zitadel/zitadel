import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { FeaturesComponent } from './features.component';

describe('FeaturesComponent', () => {
    let component: FeaturesComponent;
    let fixture: ComponentFixture<FeaturesComponent>;

    beforeEach(waitForAsync(() => {
        TestBed.configureTestingModule({
            declarations: [FeaturesComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(FeaturesComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
