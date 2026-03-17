import { GenericContainer, type ImagePullPolicy, Wait, type StartedTestContainer } from "testcontainers";

const MOCK_IMAGE = "zitadel/login-api-mock:local";
const STUBS_PORT = 22220;
const MOCK_PORT = 22222;

/** Never pull – the image must already exist locally (built by @zitadel/login-api-mock:build). */
class NeverPullPolicy implements ImagePullPolicy {
  shouldPull(): boolean {
    return false;
  }
}

let container: StartedTestContainer | undefined;

/**
 * Starts the grpc-mock container using testcontainers.
 *
 * Expects the image to already be built (via `@zitadel/login-api-mock:build`).
 * Uses fixed port binding (22220, 22222) because the login app is started
 * separately and connects to the mock at a fixed ZITADEL_API_URL (port 22222).
 *
 * Testcontainers handles readiness waiting and automatic cleanup via Ryuk
 * on process exit.
 */
export async function startMockContainer(): Promise<{
  stubsUrl: string;
  mockHost: string;
  mockPort: number;
}> {
  if (container) {
    return {
      stubsUrl: `http://localhost:${STUBS_PORT}/v1/stubs`,
      mockHost: "localhost",
      mockPort: MOCK_PORT,
    };
  }

  console.log("Starting grpc-mock container via testcontainers...");

  container = await new GenericContainer(MOCK_IMAGE)
    .withPullPolicy(new NeverPullPolicy())
    .withExposedPorts(
      { container: STUBS_PORT, host: STUBS_PORT },
      { container: MOCK_PORT, host: MOCK_PORT },
    )
    .withWaitStrategy(Wait.forListeningPorts())
    .start();

  const stubsUrl = `http://localhost:${STUBS_PORT}/v1/stubs`;
  console.log(`grpc-mock container started: stubs=${stubsUrl} mock=localhost:${MOCK_PORT}`);

  return { stubsUrl, mockHost: "localhost", mockPort: MOCK_PORT };
}

export async function stopMockContainer(): Promise<void> {
  if (container) {
    await container.stop();
    container = undefined;
  }
}
