import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { FailedEventsComponent } from './failed-events.component';

describe('FailedEventsComponent', () => {
    let component: FailedEventsComponent;
    let fixture: ComponentFixture<FailedEventsComponent>;

    beforeEach(waitForAsync(() => {
        TestBed.configureTestingModule({
            declarations: [FailedEventsComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(FailedEventsComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
