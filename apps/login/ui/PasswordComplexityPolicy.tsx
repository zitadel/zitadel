const fetcher = (url: string) =>
  fetch(url, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    },
  })
    .then((res) => {
      if (res.ok) {
        return res.json();
      } else {
        return res.json().then((error) => {
          throw error;
        });
      }
    })
    .then((resp) => resp.policy);

type Props = {
  password: string;
  equals: boolean;
  isValid: (valid: boolean) => void;
  isMe?: boolean;
  userId?: string;
};

const check = (
  <i className="las la-check text-state-success-light-color dark:text-state-success-dark-color mr-4 text-lg"></i>
);
const cross = (
  <i className="las la-times text-warn-light-500 dark:text-warn-dark-500 mr-4 text-lg"></i>
);
const desc =
  "text-14px leading-4 text-input-light-label dark:text-input-dark-label";

export default function PasswordComplexityPolicy({
  password,
  equals,
  isValid,
  isMe,
  userId,
}: Props) {
  //   const { data: policy } = useSWR<Policy, ClientError>(
  //     `/api/user/passwordpolicy/${isMe ? 'me' : userId}`,
  //     fetcher,
  //   );
  //   if (policy) {
  //     const hasMinLength = password?.length >= policy.minLength;
  //     const hasSymbol = symbolValidator(password);
  //     const hasNumber = numberValidator(password);
  //     const hasUppercase = upperCaseValidator(password);
  //     const hasLowercase = lowerCaseValidator(password);

  //     const policyIsValid =
  //       (policy.hasLowercase ? hasLowercase : true) &&
  //       (policy.hasNumber ? hasNumber : true) &&
  //       (policy.hasUppercase ? hasUppercase : true) &&
  //       (policy.hasSymbol ? hasSymbol : true) &&
  //       hasMinLength;

  //     isValid(policyIsValid);

  //     return (
  //       <div className="mb-4 grid grid-cols-2 gap-x-8 gap-y-2">
  //         <div className="flex flex-row items-center">
  //           {hasMinLength ? check : cross}
  //           <span className={desc}>Password length {policy.minLength}</span>
  //         </div>
  //         <div className="flex flex-row items-center">
  //           {hasSymbol ? check : cross}
  //           <span className={desc}>has Symbol</span>
  //         </div>
  //         <div className="flex flex-row items-center">
  //           {hasNumber ? check : cross}
  //           <span className={desc}>has Number</span>
  //         </div>
  //         <div className="flex flex-row items-center">
  //           {hasUppercase ? check : cross}
  //           <span className={desc}>has uppercase</span>
  //         </div>
  //         <div className="flex flex-row items-center">
  //           {hasLowercase ? check : cross}
  //           <span className={desc}>has lowercase</span>
  //         </div>
  //         <div className="flex flex-row items-center">
  //           {equals ? check : cross}
  //           <span className={desc}>equals</span>
  //         </div>
  //       </div>
  //     );
  //   } else {
  return null;
  //   }
}
