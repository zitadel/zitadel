export function ConsentScreen({ scope }: { scope?: string[] }) {
  return (
    <div className="flex flex-col items-center space-y-4">
      <h1>Consent</h1>
      <p className="ztdl-p">Please confirm your consent.</p>
      <div className="flex flex-col items-center space-y-4">
        <button className="btn btn-primary">Accept</button>
        <button className="btn btn-secondary">Reject</button>
      </div>
    </div>
  );
}
