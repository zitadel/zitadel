import UsernameForm from "#/ui/UsernameForm";

export default function Page() {
  return (
    <div className="flex flex-col items-center space-y-4">
      <h1>Welcome back!</h1>
      <p className="ztdl-p">Enter your login data.</p>

      <UsernameForm />
    </div>
  );
}
