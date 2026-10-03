import { ComponentFixture, TestBed } from '@angular/core/testing';

import { PaginatorComponent } from './paginator.component';
import { PaginationPreferenceService } from 'src/app/services/pagination-preference.service';

describe('PaginatorComponent', () => {
  let component: PaginatorComponent;
  let fixture: ComponentFixture<PaginatorComponent>;
  let preference: jasmine.SpyObj<PaginationPreferenceService>;

  beforeEach(async () => {
    preference = jasmine.createSpyObj<PaginationPreferenceService>('PaginationPreferenceService', ['get', 'set']);

    await TestBed.configureTestingModule({
      declarations: [PaginatorComponent],
      providers: [{ provide: PaginationPreferenceService, useValue: preference }],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(PaginatorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  describe('position label', () => {
    it('is one based so the first row reads as 1, not 0', () => {
      component.length = 161;
      component.pageSize = 20;
      component.pageIndex = 0;

      expect(component.displayStartIndex).toBe(1);
      expect(component.endIndex).toBe(20);
    });

    it('reads 0 - 0 for an empty result set', () => {
      component.length = 0;
      component.pageSize = 20;
      component.pageIndex = 0;

      expect(component.displayStartIndex).toBe(0);
      expect(component.endIndex).toBe(0);
    });

    it('caps the end index at the total on a partial last page', () => {
      component.length = 161;
      component.pageSize = 20;
      component.pageIndex = 8;

      expect(component.displayStartIndex).toBe(161);
      expect(component.endIndex).toBe(161);
    });

    it('keeps the raw offset zero based for requests', () => {
      component.length = 161;
      component.pageSize = 20;
      component.pageIndex = 2;

      expect(component.startIndex).toBe(40);
    });
  });

  describe('nextPossible', () => {
    it('is false on the last page when the total is an exact multiple of the page size', () => {
      component.length = 800;
      component.pageSize = 50;
      component.pageIndex = 15;

      expect(component.nextPossible).toBeFalse();
    });

    it('is true while full pages remain', () => {
      component.length = 803;
      component.pageSize = 50;
      component.pageIndex = 15;

      expect(component.nextPossible).toBeTrue();
    });

    it('is false once the partial last page is reached', () => {
      component.length = 803;
      component.pageSize = 50;
      component.pageIndex = 16;

      expect(component.nextPossible).toBeFalse();
    });

    it('is false when everything fits on one page', () => {
      component.length = 12;
      component.pageSize = 25;
      component.pageIndex = 0;

      expect(component.nextPossible).toBeFalse();
    });

    it('is false for an empty result set', () => {
      component.length = 0;
      component.pageSize = 25;
      component.pageIndex = 0;

      expect(component.nextPossible).toBeFalse();
    });
  });

  describe('updatePageSize', () => {
    it('returns to the first page and emits the new size', () => {
      const emitted: Array<number> = [];
      component.page.subscribe((event) => emitted.push(event.pageIndex));
      component.pageIndex = 4;

      component.updatePageSize(100);

      expect(component.pageSize).toBe(100);
      expect(component.pageIndex).toBe(0);
      expect(emitted).toEqual([0]);
    });

    it('remembers the size when a persist key is set', () => {
      component.persistKey = 'user-list';

      component.updatePageSize(250);

      expect(preference.set).toHaveBeenCalledOnceWith('user-list', 250);
    });

    it('does not touch storage without a persist key', () => {
      component.updatePageSize(250);

      expect(preference.set).not.toHaveBeenCalled();
    });
  });
});
