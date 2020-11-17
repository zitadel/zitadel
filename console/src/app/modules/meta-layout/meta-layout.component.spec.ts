import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { MetaLayoutComponent } from './meta-layout.component';

describe('MetaLayoutComponent', () => {
    let component: MetaLayoutComponent;
    let fixture: ComponentFixture<MetaLayoutComponent>;

    beforeEach(waitForAsync(() => {
        TestBed.configureTestingModule({
            declarations: [MetaLayoutComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(MetaLayoutComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
