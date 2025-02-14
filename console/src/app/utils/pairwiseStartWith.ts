import { Observable } from 'rxjs';
import { map, pairwise, startWith } from 'rxjs/operators';

export function pairwiseStartWith<T, R>(start: T) {
  return (source: Observable<R>) =>
    source.pipe(
      startWith(start),
      pairwise(),
      map(([prev, curr]) => [prev, curr] as [T | R, R]),
    );
}
