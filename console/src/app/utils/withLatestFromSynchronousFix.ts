import { combineLatestWith, distinctUntilChanged, Observable, ObservableInput, ObservableInputTuple } from 'rxjs';
import { map } from 'rxjs/operators';

// withLatestFrom does not work in this case, so we use
// combineLatestWith + distinctUntilChanged
// here the problem is described in more detail
// https://github.com/ReactiveX/rxjs/issues/7068
export const withLatestFromSynchronousFix =
  <T, A extends readonly unknown[]>(...secondaries$: [...ObservableInputTuple<A>]) =>
  (primary$: Observable<T>) =>
    primary$.pipe(
      // we add the index, so we can distinguish
      // primary submissions from each other,
      // and then we can only emit when primary changes
      map((primary, i) => <const>[primary, i]),
      combineLatestWith(...secondaries$),
      distinctUntilChanged(undefined!, ([[_, i]]) => i),
      map(([[primary], ...secondaries]) => <const>[primary, ...secondaries]),
    );
