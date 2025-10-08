import React, { useEffect, useState } from "react";
import OriginalApiItem from "@theme-original/ApiItem";

export default function ApiItemWrapper(props) {
  const [feedback, setFeedback] = useState(null);
  const [comment, setComment] = useState("");
  const [plausibleLoaded, setPlausibleLoaded] = useState(false);

  useEffect(() => {
    if (!document.getElementById("plausible-script")) {
      const script = document.createElement("script");
      script.src = "https://plausible.io/js/plausible.js";
      script.defer = true;
      script.async = true;
      script.dataset.domain = window.location.hostname;
      script.id = "plausible-script";

      script.onload = () => {
        setPlausibleLoaded(true);
      };

      document.body.appendChild(script);
    } else {
      setPlausibleLoaded(true);
    }
  }, []);

  const sendPlausibleEvent = (feedbackValue, commentText = "") => {
    if (plausibleLoaded && typeof window.plausible === "function") {
      console.log("Sending event:", feedbackValue, commentText);
      window.plausible("docs-feedback", {
        props: {
          feedback: feedbackValue,
          comment: commentText,
          page: window.location.pathname,
        },
      });
    } else {
      console.warn("Plausible not loaded yet. Event skipped.");
    }
  };

  const handleYes = () => {
    setFeedback("yes");
    sendPlausibleEvent("yes");
  };

  const handleNo = () => {
    setFeedback("no");
    sendPlausibleEvent("no");
  };

  const handleSubmitComment = () => {
    sendPlausibleEvent("no", comment);
    setFeedback("submitted");
  };

  return (
    <>
      <OriginalApiItem {...props} />

      <div className="mt-10 flex justify-start">
        {feedback === "no" && (
          <div
            className="w-fit rounded-full border border-gray-300 dark:border-gray-700 bg-white dark:bg-[#ffffff10] px-2 py-2 shadow-sm transition-all duration-300"
            style={{
              display: "flex",
              flexDirection: "column",
            }}
          >
            <div className="flex items-center justify-between flex-wrap gap-3 flex-grow">
              <p
                className="font-base ml-4 mr-4 my-0"
                style={{ color: "var(--ifm-menu-color)" }}
              >
                Was this page useful?
              </p>
              <div className="flex gap-3">
                <button
                  onClick={handleNo}
                  disabled={!plausibleLoaded}
                  className="group bg-[#00000010] dark:bg-[#00000020] rounded-full py-1 px-4 flex items-center"
                  style={{
                    border: "none",
                    cursor: plausibleLoaded ? "pointer" : "not-allowed",
                    opacity: plausibleLoaded ? 1 : 0.5,
                  }}
                  title="No, needs improvement"
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    strokeWidth={1.5}
                    stroke="currentColor"
                    className="size-6 mr-2 group-hover:scale-110 group-hover:text-blue-500 transition-transform duration-200"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      d="M15.182 16.318A4.486 4.486 0 0 0 12.016 15a4.486 4.486 0 0 0-3.198 1.318M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0ZM9.75 9.75c0 .414-.168.75-.375.75S9 10.164 9 9.75 9.168 9 9.375 9s.375.336.375.75Zm-.375 0h.008v.015h-.008V9.75Zm5.625 0c0 .414-.168.75-.375.75s-.375-.336-.375-.75.168-.75.375-.75.375.336.375.75Zm-.375 0h.008v.015h-.008V9.75Z"
                    />
                  </svg>
                  No
                </button>

                <button
                  onClick={handleYes}
                  disabled={!plausibleLoaded}
                  className="group bg-[#00000010] dark:bg-[#00000020] rounded-full py-1 px-4 flex items-center"
                  style={{
                    border: "none",
                    cursor: plausibleLoaded ? "pointer" : "not-allowed",
                    opacity: plausibleLoaded ? 1 : 0.5,
                  }}
                  title="Yes, helpful!"
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    strokeWidth={1.5}
                    stroke="currentColor"
                    className="size-6 mr-2 group-hover:scale-110 group-hover:text-amber-500 transition-transform duration-200"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      d="M15.182 15.182a4.5 4.5 0 0 1-6.364 0M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0ZM9.75 9.75c0 .414-.168.75-.375.75S9 10.164 9 9.75 9.168 9 9.375 9s.375.336.375.75Zm-.375 0h.008v.015h-.008V9.75Zm5.625 0c0 .414-.168.75-.375.75s-.375-.336-.375-.75.168-.75.375-.75.375.336.375.75Zm-.375 0h.008v.015h-.008V9.75Z"
                    />
                  </svg>
                  Yes
                </button>
              </div>
            </div>
          </div>
        )}

        {feedback === null && (
          <div
            className="w-full rounded-lg border border-gray-300 dark:border-gray-700 bg-white dark:bg-[#ffffff10] px-2 py-2 shadow-sm transition-all duration-300"
            style={{
              display: "flex",
              flexDirection: "column",
            }}
          >
            <div className="mt-1 w-full flex flex-row items-center gap-3 flex-grow">
              <button
                className="button"
                style={{ color: "var(--color-fg)" }}
                onClick={() => {
                  setFeedback(null);
                  setComment("");
                }}
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  strokeWidth={1.5}
                  stroke="currentColor"
                  className="size-6"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    d="M10.5 19.5 3 12m0 0 7.5-7.5M3 12h18"
                  />
                </svg>
              </button>
              <textarea
                className="textarea textarea-bordered w-full resize-none rounded-md border border-gray-300 dark:border-gray-600 bg-transparent p-2 text-sm"
                rows="3"
                placeholder="Sorry to hear that — how can we improve this page?"
                value={comment}
                onChange={(e) => setComment(e.target.value)}
              />
              <div className="flex gap-2">
                <button
                  className="button button--sm button--primary h-fit"
                  onClick={handleSubmitComment}
                  disabled={comment.trim() === "" || !plausibleLoaded}
                  style={{
                    opacity: comment.trim() === "" ? 0.5 : 1,
                    cursor:
                      comment.trim() === "" || !plausibleLoaded
                        ? "not-allowed"
                        : "pointer",
                  }}
                >
                  Submit
                </button>
              </div>
            </div>
          </div>
        )}

        {(feedback === "yes" || feedback === "submitted") && (
          <div
            className="w-fit rounded-full border border-gray-300 dark:border-gray-700 bg-white dark:bg-[#ffffff10] px-2 py-2 shadow-sm transition-all duration-300"
            style={{
              display: "flex",
              flexDirection: "column",
            }}
          >
            <div className="flex items-center justify-center flex-grow">
              <p className="font-medium text-center m-0 mx-4">
                Thanks for your feedback! 🎉
              </p>
            </div>
          </div>
        )}
      </div>
    </>
  );
}
