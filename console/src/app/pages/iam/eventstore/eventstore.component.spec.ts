import { ComponentFixture, TestBed } from '@angular/core/testing';

import { EventstoreComponent } from './eventstore.component';

describe('EventstoreComponent', () => {
    let component: EventstoreComponent;
    let fixture: ComponentFixture<EventstoreComponent>;

    beforeEach(async () => {
        await TestBed.configureTestingModule({
            declarations: [EventstoreComponent],
        })
            .compileComponents();
    });

    beforeEach(() => {
        fixture = TestBed.createComponent(EventstoreComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
