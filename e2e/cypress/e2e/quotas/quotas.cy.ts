import { systemAuth } from 'support/api/apiauth';
import { instanceUnderTest } from 'support/api/instances';
import { addQuota, removeQuota } from 'support/api/quota';

const authenticatedRequestsUnit = 1;
const executionSeconds = 2;

describe('quotas', () => {
  before(function () {
    systemAuth().as('api');
  });

  before(function () {
    instanceUnderTest(this.api).as('instanceId');
  });

  describe('management', () => {
    describe('add one quota', () => {
      beforeEach(function () {
        removeQuota(this.api, this.instanceId, authenticatedRequestsUnit, false).then((res) => {
          if (!res.isOkStatusCode) {
            expect(res.status).to.equal(404);
          }
        });
      });
      it('should add a quota only once per unit', function () {
        addQuota(this.api, this.instanceId, authenticatedRequestsUnit);
        addQuota(this.api, this.instanceId, authenticatedRequestsUnit, false).then((res) => {
          expect(res.status).to.equal(409);
        });
      });

      describe('add two quotas', () => {
        beforeEach(function () {
          removeQuota(this.api, this.instanceId, executionSeconds, false).then((res) => {
            if (!res.isOkStatusCode) {
              expect(res.status).to.equal(404);
            }
          });
        });
        it('should add a quota for each unit', function () {
          addQuota(this.api, this.instanceId, authenticatedRequestsUnit);
          addQuota(this.api, this.instanceId, executionSeconds);
        });
      });
    });

    describe('edit', () => {
      describe('remove one quota', () => {
        beforeEach(function () {
          addQuota(this.api, this.instanceId, authenticatedRequestsUnit, false).then((res) => {
            if (!res.isOkStatusCode) {
              expect(res.status).to.equal(409);
            }
          });
        });
        it('should remove a quota only once per unit', function () {
          removeQuota(this.api, this.instanceId, authenticatedRequestsUnit);
          removeQuota(this.api, this.instanceId, authenticatedRequestsUnit, false).then((res) => {
            expect(res.status).to.equal(404);
          });
        });

        describe('remove two quotas', () => {
          beforeEach(function () {
            addQuota(this.api, this.instanceId, executionSeconds, false).then((res) => {
              if (!res.isOkStatusCode) {
                expect(res.status).to.equal(409);
              }
            });
          });
          it('should remove a quota for each unit', function () {
            removeQuota(this.api, this.instanceId, authenticatedRequestsUnit);
            removeQuota(this.api, this.instanceId, executionSeconds);
          });
        });
      });
    });
  });

  describe('usage', () => {
    beforeEach(() => {
      cy.task('runSQL', 'TRUNCATE logstore.access;');
    });
  });
});
