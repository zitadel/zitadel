// Polyfill "window.fetch" used in the React component.
import 'whatwg-fetch'

import * as mockRouter from 'next-router-mock';
import { mockLoginBackend } from './mock-login-backend'

// Inspired by https://github.com/scottrippey/next-router-mock/issues/67#issuecomment-1564906960
const useRouter = mockRouter.useRouter;

const MockNextNavigation = {
  ...mockRouter,
  usePathname: () => {
    const router = useRouter();
    return router.pathname;
  },
  useSearchParams: () => {
    const router = useRouter();
    const path = router.asPath.split('?')?.[1] ?? '';
    return new URLSearchParams(path);
  },
};

// jest.mock('next/navigation', () => MockNextNavigation);
// jest.mock('next/router', () => mockRouter)

beforeAll(() => {
  mockLoginBackend.listen()
})
beforeEach(() => mockLoginBackend.resetHandlers())
afterAll(() => mockLoginBackend.close())