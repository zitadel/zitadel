import { RegisterTOTPResponse } from "@zitadel/server";

export default function TOTPRegister({
  uri,
  secret,
}: {
  uri: string;
  secret: string;
}) {
  return <div>{uri}</div>;
}
