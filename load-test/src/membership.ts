import http from 'k6/http';
import { Trend } from 'k6/metrics';
import url from './url';
import { check, fail } from 'k6';

const addIAMMemberTrend = new Trend('membership_iam_member', true);
export async function addIAMMember(userId: string, roles: string[], accessToken: string): Promise<void> {
  const res = await http.post(
    url('/admin/v1/members'),
    JSON.stringify({
      userId: userId,
      roles: roles,
    }),
    {
      headers: {
        authorization: `Bearer ${accessToken}`,
        'Content-Type': 'application/json',
      },
    },
  );
  check(res, {
    'member added successful': (r) => r.status >= 200 && r.status < 300 || fail(`unable add member: ${JSON.stringify(res)}`),
  });
  addIAMMemberTrend.add(res.timings.duration);
}
