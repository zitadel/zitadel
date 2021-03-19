import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { ClientKeysComponent } from './client-keys.component';

describe('ClientKeysComponent', () => {
    let component: ClientKeysComponent;
    let fixture: ComponentFixture<ClientKeysComponent>;

    beforeEach(waitForAsync(() => {
        TestBed.configureTestingModule({
            declarations: [ClientKeysComponent],
        })
            .compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(ClientKeysComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
