---
title: Update and Scale ZITADEL
sidebar_label: Update and Scale
---

## TL;DR

For getting started easily, you can just combine the init, setup and runtime phases
by executing the `zitadel` binary with the argument `start-from-init`.
However, to improve day-two-operations,
we recommend you run the init phase and the setup phase
separately using `zitadel init` and  `zitadel setup` respectively.
Both the init and the setup phases are idempotent,
meaning you can run them as many times as you want whithout unexpectedly changing any state.
Nevertheless, they just __need__ to be run at certain points in time, whereas further runs are obsolete.
The init phase has to run once over the full ZITADEL lifecycle, before the setup and runtime phases run first.
The setup phase has to run each time a new ZITADEL version is deployed.
After init and setup is done, you can start the runtime phase, that immediately accepts traffic, using `zitadel start`.

## Scaling ZITADEL

For production setups, we recommend that you [run ZITADEL in a highly available mode](/docs/self-hosting/manage/production).
As all shared state is managed in the database,
the ZITADEL binary itself is stateless.
You can easily run as many ZITADEL binaries in parallel as you want.
If you use [Knative](/docs/self-hosting/deploy/knative)
or [Google Cloud Run](https://cloud.google.com/run) (which uses Knative, too),
you can even scale to zero.
Especially if you use an autoscaler that scales to zero,
it is crucial that you minimize startup times.

## Updating ZITADEL

When want to deploy a new ZITADEL version,
the new versions setup phase takes care of database migrations.
You generally want to run the job in a controlled manner,
so multiple executions don’t interfere with each other.
Also, after the setup is done,
rolling out a new ZITADEL version is much faster
when the runtime processes are just executed with `zitadel start`.

## Separating Init and Setup from the Runtime

If you use the [official ZITADEL Helm chart](/docs/self-hosting/deploy/kubernetes),
then you can stop reading now.
The init and setup phases are already separated and executed in dedicated Kubernetes Jobs.
If you use another orchestrator and want to separate the phases yourself,
you should know what happens during these phases.

### The Init Phase

The command `zitadel init` ensures the database connection is ready to use for the subsequent phases.
It just needs to be executed once over ZITADELs full life cycle,
when you install ZITADEL from scratch.
During `zitadel init`, for connecting to your database,
ZITADEL uses the privileged and preexisting database user configured in `Database.postgres.Admin.Username`.
, `zitadel init` ensures the following:
- If it doesn’t exist already, it creates a database with the configured database name.
- If it doesn’t exist already, it creates the unprivileged user use configured in `Database.postgres.User.Username`.
  Subsequent phases connect to the database with this user's credentials only.
- If not already done, it grants the necessary permissions ZITADEL needs to the non privileged user.
- If they don’t exist already, it creates all schemas and some basic tables.

The init phase is idempotent if executed with the same binary version.

### The Setup Phase

During `zitadel setup`, ZITADEL creates projection tables and migrates existing data, if `--init-projections=true` is set.
Depending on the ZITADEL version and the runtime resources,
this step can take several minutes.
When deploying a new ZITADEL version,
make sure the setup phase runs before you roll out the new `zitadel start` processes.
The setup phase is executed in subsequent steps
whereas a new version's execution takes over where the last execution stopped.

Some configuration changes are only applied during the setup phase, like ExternalDomain, ExternalPort and ExternalSecure.

The setup phase is idempotent if executed with the same binary version.

### The Runtime Phase

The `zitadel start` command assumes the database is already initialized and set up.
It starts serving requests within fractions of a second.
Beware, in the background, out-of-date projections
[recompute their state by replaying all missed events](/docs/concepts/eventstore/implementation#projections).
If a new ZITADEL version is deployed, this can take quite a long time,
depending on the amount of events to catch up.
You probably should consider providing `--init-projections=true`-flag to the [Setup Phase](#the-setup-phase) to shift the synchronization time to previous steps and delay the startup phase until events are caught up.