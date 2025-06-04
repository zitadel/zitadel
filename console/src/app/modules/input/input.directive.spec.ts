import { Component, ElementRef, NgZone } from '@angular/core';
import { TestBed, ComponentFixture } from '@angular/core/testing';
import { InputDirective } from './input.directive';
import { Platform } from '@angular/cdk/platform';
import { NgControl, NgForm, FormGroupDirective } from '@angular/forms';
import { ErrorStateMatcher } from '@angular/material/core';
import { AutofillMonitor } from '@angular/cdk/text-field';
import { MatFormField } from '@angular/material/form-field';
import { MAT_INPUT_VALUE_ACCESSOR } from '@angular/material/input';
import { of } from 'rxjs';
import { By } from '@angular/platform-browser';

@Component({
  template: `<input appInputDirective />`,
})
class TestHostComponent {}

describe('InputDirective', () => {
  let fixture: ComponentFixture<TestHostComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [InputDirective, TestHostComponent],
      providers: [
        { provide: ElementRef, useValue: new ElementRef(document.createElement('input')) },
        Platform,
        { provide: NgControl, useValue: null },
        { provide: NgForm, useValue: null },
        { provide: FormGroupDirective, useValue: null },
        ErrorStateMatcher,
        { provide: MAT_INPUT_VALUE_ACCESSOR, useValue: null },
        {
          provide: AutofillMonitor,
          useValue: { monitor: () => of(), stopMonitoring: () => {} },
        },
        NgZone,
        { provide: MatFormField, useValue: null },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(TestHostComponent);
    fixture.detectChanges();
  });

  it('should create an instance', () => {
    const directiveEl = fixture.debugElement.query(By.directive(InputDirective));
    expect(directiveEl).toBeTruthy();
  });
});
