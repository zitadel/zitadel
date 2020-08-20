import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { FailedEventsComponent } from './failed-events.component';

describe('FailedEventsComponent', () => {
    let component: FailedEventsComponent;
    let fixture: ComponentFixture<FailedEventsComponent>;

    beforeEach(async(() => {
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
