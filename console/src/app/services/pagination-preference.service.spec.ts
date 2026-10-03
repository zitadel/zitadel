import { TestBed } from '@angular/core/testing';

import { PaginationPreferenceService } from './pagination-preference.service';
import { StorageLocation, StorageService } from './storage.service';

describe('PaginationPreferenceService', () => {
  let service: PaginationPreferenceService;
  let storage: jasmine.SpyObj<StorageService>;

  const OPTIONS = [10, 25, 50, 100, 250];

  beforeEach(() => {
    storage = jasmine.createSpyObj<StorageService>('StorageService', ['getItem', 'setItem']);

    TestBed.configureTestingModule({
      providers: [PaginationPreferenceService, { provide: StorageService, useValue: storage }],
    });

    service = TestBed.inject(PaginationPreferenceService);
  });

  describe('get', () => {
    it('returns the stored size when it is one of the offered options', () => {
      storage.getItem.and.returnValue(100);

      expect(service.get('user-list', 25, OPTIONS)).toBe(100);
      expect(storage.getItem).toHaveBeenCalledOnceWith('pagesize-user-list', StorageLocation.local);
    });

    it('falls back when nothing is stored', () => {
      storage.getItem.and.returnValue(null);

      expect(service.get('user-list', 25, OPTIONS)).toBe(25);
    });

    it('falls back when the stored size is no longer offered', () => {
      // a size written before the option set changed must not reach the request
      storage.getItem.and.returnValue(20);

      expect(service.get('user-list', 25, OPTIONS)).toBe(25);
    });

    it('falls back on a non numeric value', () => {
      storage.getItem.and.returnValue('fifty' as unknown as number);

      expect(service.get('user-list', 25, OPTIONS)).toBe(25);
    });

    it('falls back on a non positive value', () => {
      storage.getItem.and.returnValue(0);

      expect(service.get('user-list', 25, OPTIONS)).toBe(25);
    });

    it('falls back when storage throws', () => {
      storage.getItem.and.throwError('storage blocked');

      expect(service.get('user-list', 25, OPTIONS)).toBe(25);
    });

    it('accepts any positive integer when no options are given', () => {
      storage.getItem.and.returnValue(37);

      expect(service.get('user-list', 25)).toBe(37);
    });
  });

  describe('set', () => {
    it('writes the size to local storage under a prefixed key', () => {
      service.set('user-grants', 250);

      expect(storage.setItem).toHaveBeenCalledOnceWith('pagesize-user-grants', 250, StorageLocation.local);
    });

    it('swallows storage errors so the table keeps working', () => {
      storage.setItem.and.throwError('quota exceeded');

      expect(() => service.set('user-grants', 250)).not.toThrow();
    });
  });
});
