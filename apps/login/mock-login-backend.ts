// Inspired by https://kentcdodds.com/blog/stop-mocking-fetch
// These handlers simulate the login backend
// They are used in tests and for local development
import {DefaultBodyType, PathParams, ResponseComposition, RestContext, RestRequest, rest} from 'msw'
import {setupServer} from 'msw/node'

const handlers = [
  rest.post('/verifyemail', async (req, res, ctx) => {
//    checkAuthorized(req,res,ctx)
    return res(ctx.json({
        sequence: 111,
        changeDate: "2023-01-01T00:00:00.000Z",
        resourceOwner: "111111111111111111"
    }))
  }),
]

const checkAuthorized = (req: RestRequest<DefaultBodyType, PathParams<string>>, res: ResponseComposition<DefaultBodyType>, ctx: RestContext) => {
  if (!req.headers.has('Authorization')){
    const err = "Not authorized"
    res(ctx.status(401), ctx.json({message: err}))
    throw err
  }
}

const mockLoginBackend = setupServer(...handlers)
export {mockLoginBackend, rest}