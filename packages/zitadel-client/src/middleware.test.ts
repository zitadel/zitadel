import { CallOptions, ClientMiddlewareCall, Metadata, MethodDescriptor } from "nice-grpc-web";
import { authMiddleware } from "./middleware";

describe('authMiddleware', () => {
    const scenarios = [
      {
        name: 'should add authorization if metadata is undefined',
        initialMetadata: undefined,
        expectedMetadata: new Metadata().set("authorization", "Bearer mock-token"),
        token: "mock-token"
      },
      {
        name: 'should add authorization if metadata exists but no authorization',
        initialMetadata: new Metadata().set("other-key", "other-value"),
        expectedMetadata: new Metadata().set("other-key", "other-value").set("authorization", "Bearer mock-token"),
        token: "mock-token"
      },
      {
        name: 'should not modify authorization if it already exists',
        initialMetadata: new Metadata().set("authorization", "Bearer initial-token"),
        expectedMetadata: new Metadata().set("authorization", "Bearer initial-token"),
        token: "mock-token"
      },
    ];

    scenarios.forEach(({ name, initialMetadata, expectedMetadata, token }) => {
      it(name, async () => {

        const mockNext = jest.fn().mockImplementation(async function*() { });
        const mockRequest = {};

        const mockMethodDescriptor: MethodDescriptor = {
            options: {idempotencyLevel: undefined},
            path: '',
            requestStream: false,
            responseStream: false,
        };

        const mockCall: ClientMiddlewareCall<unknown, unknown> = {
            method: mockMethodDescriptor,
            requestStream: false,
            responseStream: false,
            request: mockRequest,
            next: mockNext,
          };
        const options: CallOptions = {
          metadata: initialMetadata
        };

        await authMiddleware(token)(mockCall, options).next();

        expect(mockNext).toHaveBeenCalledTimes(1);
        const actualMetadata = mockNext.mock.calls[0][1].metadata;
        expect(actualMetadata?.get('authorization')).toEqual(expectedMetadata.get('authorization'));
      });
    });
  });
