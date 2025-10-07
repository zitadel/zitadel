import React, { useState } from "react";
import OriginalApiItem from "@theme-original/ApiItem";

export default function ApiItemWrapper(props) {
  const [feedback, setFeedback] = useState(null);
  const [comment, setComment] = useState("");

  const sendPlausibleEvent = (feedbackValue, commentText = "") => {
    if (typeof window.plausible === "function") {
      window.plausible("Feedback Submitted", {
        props: {
          feedback: feedbackValue,
          comment: commentText,
          page: window.location.pathname,
        },
      });
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
        <div className="w-full max-w-2xl rounded-xl border border-gray-300 dark:border-gray-700 bg-white dark:bg-[#1a1a1a] p-6 shadow-sm">
          {/* Title and buttons aligned horizontally */}
          {feedback === null && (
            <div className="flex items-center justify-between flex-wrap gap-3">
              <p
                className="font-medium text-lg m-0"
                style={{ color: "var(--ifm-menu-color)" }}
              >
                Was this page useful?
              </p>
              <div className="flex gap-3">
                <button
                  onClick={handleYes}
                  className="button button--sm"
                  style={{
                    backgroundColor: "transparent",
                    color: "var(--ifm-color-primary)",
                    border: "1px solid var(--ifm-color-primary)",
                  }}
                >
                  Yes
                </button>
                <button
                  className="button button--sm button--secondary"
                  onClick={handleNo}
                >
                  No
                </button>
              </div>
            </div>
          )}

          {feedback === "no" && (
            <div className="mt-1 flex flex-col gap-3">
              <textarea
                className="textarea textarea-bordered w-full resize-none rounded-md border border-gray-300 dark:border-gray-600 bg-transparent p-2 text-sm"
                rows="3"
                placeholder="Sorry to hear that — how can we improve this page?"
                value={comment}
                onChange={(e) => setComment(e.target.value)}
              />
              <div className="flex gap-2">
                <button
                  className="button button--sm button--primary"
                  onClick={handleSubmitComment}
                >
                  Submit
                </button>
                <button
                  className="button button--sm button--secondary"
                  onClick={() => {
                    setFeedback(null);
                    setComment("");
                  }}
                >
                  Go Back
                </button>
              </div>
            </div>
          )}

          {(feedback === "yes" || feedback === "submitted") && (
            <p
              className="mt-4 font-medium"
              style={{ color: "var(--ifm-color-primary)" }}
            >
              Thanks for your feedback!
            </p>
          )}
        </div>
      </div>
    </>
  );
}
