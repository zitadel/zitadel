document.getElementById("back-button").addEventListener("click", goBack);
document.getElementById("back-button").style.visibility = wereInUserSelection()
  ? "visible"
  : "hidden";

function goBack() {
  history.back();
}

function wereInUserSelection() {
  return window.location.href.includes("userselection");
}
