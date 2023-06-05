import { mockLoginBackend } from './mock-login-backend';

mockLoginBackend.printHandlers()
mockLoginBackend.listen()

console.log("listened")