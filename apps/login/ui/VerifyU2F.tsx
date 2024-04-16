type Props = {
  loginName: string | undefined;
  sessionId: string | undefined;
  authRequestId?: string;
  organization?: string;
  submit: boolean;
};

export default function VerifyU2F({
  loginName,
  authRequestId,
  organization,
  submit,
}: Props) {
  return <div>Verify U2F</div>;
}
