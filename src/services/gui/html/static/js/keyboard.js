const Keyboard = window.SimpleKeyboard.default;
const keyboardTheme = "hg-theme-default dark";
const keyboard = new Keyboard({
  theme: keyboardTheme,
  onChange: (input) => onChange(input),
  onKeyPress: (button) => onKeyPress(button),
});

function onChange(input) {
  document.getElementById("inputName").value = input;
}

function onKeyPress(button) {
  if (button === "{shift}" || button === "{lock}") handleShift();
}

function handleShift() {
  let currentLayout = keyboard.options.layoutName;
  let shiftToggle = currentLayout === "default" ? "shift" : "default";
  keyboard.setOptions({ layoutName: shiftToggle });
}

/**
 * Keyboard show
 */
const inputName = document.getElementById("inputName");
const elemsHideOnKeyboard = document.getElementsByClassName("hide-on-keyboard");
inputName.addEventListener("focus", (event) => {
  showKeyboard();
});

/**
 * Keyboard show toggle
 */
document.addEventListener("click", (event) => {
  if (
    /**
     * Hide the keyboard when you're not clicking it or when clicking an input
     * If you have installed a "click outside" library, please use that instead.
     */
    keyboard.options.theme.includes("show-keyboard") &&
    event.target.id !== "inputName" &&
    !event.target.className.includes("hg-button") &&
    !event.target.className.includes("hg-row") &&
    !event.target.className.includes("simple-keyboard")
  ) {
    hideKeyboard();
  }
});

function showKeyboard() {
  keyboard.setOptions({
    theme: `${keyboardTheme} show-keyboard`,
  });
  Array.prototype.forEach.call(elemsHideOnKeyboard, (elem) =>
    elem.classList.add("d-none")
  );
}

function hideKeyboard() {
  keyboard.setOptions({
    theme: keyboardTheme,
  });
  Array.prototype.forEach.call(elemsHideOnKeyboard, (elem) =>
    elem.classList.remove("d-none")
  );
}
