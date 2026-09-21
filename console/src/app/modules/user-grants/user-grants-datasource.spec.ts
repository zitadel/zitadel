import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ManagementService } from 'src/app/services/mgmt.service';

import { UserGrantContext, UserGrantsDataSource } from './user-grants-datasource';

describe('UserGrantsDataSource', () => {
  let mgmtService: jasmine.SpyObj<ManagementService>;
  let authService: jasmine.SpyObj<GrpcAuthService>;
  let dataSource: UserGrantsDataSource;

  function response(ids: Array<string>): any {
    return {
      resultList: ids.map((id) => ({ id })),
      details: { totalResult: ids.length },
    };
  }

  /** A promise the test resolves by hand, to control the order responses arrive in. */
  function deferred(): { promise: Promise<any>; resolve: (value: any) => void } {
    let resolve!: (value: any) => void;
    const promise = new Promise<any>((r) => (resolve = r));
    return { promise, resolve };
  }

  beforeEach(() => {
    mgmtService = jasmine.createSpyObj<ManagementService>('ManagementService', ['listUserGrants']);
    authService = jasmine.createSpyObj<GrpcAuthService>('GrpcAuthService', ['listMyUserGrants']);
    dataSource = new UserGrantsDataSource(authService, mgmtService);
  });

  // Regression: the search box fires overlapping requests while typing. Responses are not
  // guaranteed to arrive in order, and the table used to show whichever landed last.
  it('ignores a response that a newer request superseded', async () => {
    const first = deferred();
    const second = deferred();
    mgmtService.listUserGrants.and.returnValues(first.promise, second.promise);

    dataSource.loadGrants(UserGrantContext.NONE, 0, 50, {});
    dataSource.loadGrants(UserGrantContext.NONE, 0, 50, {});

    second.resolve(response(['new']));
    await second.promise;

    first.resolve(response(['stale']));
    await first.promise;
    await Promise.resolve();

    expect(dataSource.grantsSubject.value.map((g: any) => g.id)).toEqual(['new']);
  });

  it('applies a response when no newer request was started', async () => {
    const only = deferred();
    mgmtService.listUserGrants.and.returnValue(only.promise);

    dataSource.loadGrants(UserGrantContext.NONE, 0, 50, {});
    only.resolve(response(['a', 'b']));
    await only.promise;
    await Promise.resolve();

    expect(dataSource.grantsSubject.value.map((g: any) => g.id)).toEqual(['a', 'b']);
    expect(dataSource.totalResult).toBe(2);
  });

  it('does not let a superseded failure clear the table', async () => {
    const first = deferred();
    const second = deferred();
    mgmtService.listUserGrants.and.returnValues(first.promise, second.promise);

    dataSource.loadGrants(UserGrantContext.NONE, 0, 50, {});
    dataSource.loadGrants(UserGrantContext.NONE, 0, 50, {});

    second.resolve(response(['new']));
    await second.promise;

    first.resolve(Promise.reject(new Error('stale failure')));
    await Promise.resolve();
    await Promise.resolve();

    expect(dataSource.grantsSubject.value.map((g: any) => g.id)).toEqual(['new']);
  });
});
