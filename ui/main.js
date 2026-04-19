import {
  initialiseDarkMode,
  renderDarkModeToggle,
  renderFooter,
  renderHeader,
} from "./shared.js";

const STORAGE_KEY = "songstitchform";
const USERNAME_REGEX = /^[a-zA-Z][a-zA-Z0-9_-]{0,15}$/;

const DEFAULT_VALUES = {
  username: "",
  method: "album",
  period: "7day",
  track: true,
  artist: true,
  album: true,
  playcount: true,
  rows: 3,
  columns: 3,
  advancedOptions: false,
  showTextSize: false,
  showTextLocation: false,
  textSize: "12",
  textLocation: "topleft",
  showBoldtext: false,
  grayscaleImage: false,
  WebPLossyCompression: false,
};

const app = document.getElementById("app");

app.innerHTML = `
  <main>
    ${renderHeader()}
    <form id="collage-form" novalidate>
      ${renderDarkModeToggle()}
      <div class="form-group">
        <label class="form-heading" for="username">Last.fm Username</label>
        <input
          class="username"
          type="text"
          name="username"
          id="username"
          autocomplete="on"
          placeholder="Enter your username"
        />
        <div class="error" id="username-error" hidden></div>
      </div>

      <div class="form-row">
        <div class="form-group form-group-half">
          <label class="form-heading" for="method">Type</label>
          <select name="method" id="method">
            <option value="album">Top Albums</option>
            <option value="artist">Top Artists</option>
            <option value="track">Top Tracks</option>
          </select>
        </div>

        <div class="form-group form-group-half">
          <label class="form-heading" for="period">Time Period</label>
          <select name="period" id="period">
            <option value="7day">7 Days</option>
            <option value="1month">1 Month</option>
            <option value="3month">3 Months</option>
            <option value="6month">6 Months</option>
            <option value="12month">Year</option>
            <option value="overall">All Time</option>
          </select>
        </div>
      </div>

      <fieldset id="fieldset">
        <legend class="legend">Display Options</legend>

        <div class="checkbox-wrapper" id="track-wrapper">
          <input type="checkbox" class="switch" name="track" id="track" />
          <label class="checkbox-label" for="track">Track Name</label>
        </div>

        <div class="checkbox-wrapper" id="artist-wrapper">
          <input type="checkbox" class="switch" name="artist" id="artist" />
          <label class="checkbox-label" for="artist">Artist Name</label>
        </div>

        <div class="checkbox-wrapper" id="album-wrapper">
          <input type="checkbox" class="switch" name="album" id="album" />
          <label class="checkbox-label" for="album">Album Name</label>
        </div>

        <div class="checkbox-wrapper" id="playcount-wrapper">
          <input
            type="checkbox"
            class="switch"
            name="playcount"
            id="playcount"
          />
          <label class="checkbox-label" for="playcount">Play Count</label>
        </div>
      </fieldset>

      <fieldset id="grid-fieldset">
        <legend class="legend">Grid Size</legend>
        <div class="grid-size-section">
          <div class="grid-inputs">
            <div class="form-row">
              <div class="form-group form-group-half">
                <label class="label" for="rows">
                  Rows <span class="limit">(max <span id="rows-max">20</span>)</span>
                </label>
                <input
                  class="number-input"
                  inputmode="decimal"
                  type="number"
                  max="20"
                  min="0"
                  name="rows"
                  id="rows"
                />
                <div class="error error-text" id="rows-error" hidden></div>
              </div>

              <div class="form-group form-group-half">
                <label class="label" for="columns">
                  Columns <span class="limit">(max <span id="columns-max">20</span>)</span>
                </label>
                <input
                  class="number-input"
                  inputmode="decimal"
                  type="number"
                  max="20"
                  min="0"
                  name="columns"
                  id="columns"
                />
                <div class="error error-text" id="columns-error" hidden></div>
              </div>
            </div>
          </div>
          <div class="grid-preview" id="grid-preview"></div>
        </div>
      </fieldset>

      <fieldset id="advanced-fieldset">
        <legend class="legend">
          <div class="checkbox-wrapper advanced-legend-toggle" id="advanced-options-wrapper">
            <input
              type="checkbox"
              class="switch"
              name="advancedOptions"
              id="advancedOptions"
            />
            <label class="checkbox-label" for="advancedOptions"
              >Advanced Options</label
            >
          </div>
        </legend>

        <div class="advanced-options" id="advanced-options-content">
          <div class="checkbox-wrapper" id="grayscale-wrapper">
            <input
              type="checkbox"
              class="switch"
              name="grayscaleImage"
              id="grayscaleImage"
            />
            <label class="checkbox-label" for="grayscaleImage"
              >Grayscale</label
            >
          </div>

          <div class="checkbox-wrapper" id="bold-wrapper">
            <input
              type="checkbox"
              class="switch"
              name="showBoldtext"
              id="showBoldtext"
            />
            <label class="checkbox-label" for="showBoldtext">Bold Text</label>
          </div>

          <div class="checkbox-wrapper" id="text-size-wrapper">
            <input
              type="checkbox"
              class="switch"
              name="showTextSize"
              id="showTextSize"
            />
            <label class="checkbox-label" for="showTextSize"
              >Custom Font Size</label
            >
          </div>

          <div class="sub-options" id="text-size-options" hidden>
            <label class="advanced-option-label" for="textSize">Font Size</label>
            <select name="textSize" id="textSize">
              <option value="10">Extra Small</option>
              <option value="12">Small (default)</option>
              <option value="15">Medium</option>
              <option value="18">Large</option>
            </select>
          </div>

          <div class="checkbox-wrapper" id="text-location-wrapper">
            <input
              type="checkbox"
              class="switch"
              name="showTextLocation"
              id="showTextLocation"
            />
            <label class="checkbox-label" for="showTextLocation"
              >Custom Text Position</label
            >
          </div>

          <div class="sub-options" id="text-location-options" hidden>
            <label class="advanced-option-label" for="textLocation">Position</label>
            <select name="textLocation" id="textLocation">
              <option value="topleft">Top Left (default)</option>
              <option value="topcentre">Top Centre</option>
              <option value="topright">Top Right</option>
              <option value="bottomleft">Bottom Left</option>
              <option value="bottomcentre">Bottom Centre</option>
              <option value="bottomright">Bottom Right</option>
            </select>
          </div>

          <div id="webp-container">
            <div class="checkbox-wrapper" id="webp-wrapper">
              <input
                type="checkbox"
                class="switch"
                name="WebPLossyCompression"
                id="WebPLossyCompression"
              />
              <label class="checkbox-label" for="WebPLossyCompression"
                >WebP Compression</label
              >
            </div>
          </div>
        </div>
      </fieldset>

      <div class="loader-container">
        <div class="loader" id="form-loader" hidden></div>
      </div>

      <div class="form-actions">
        <button
          id="submit-button"
          class="btn-primary"
          type="submit"
        >Generate Collage</button>

        <button
          id="embed-button"
          class="btn-secondary"
          type="button"
        >Share / Embed</button>
      </div>

      <div class="reset-button">
        <a class="reset-text" href="#top" id="reset-link">Reset Form</a>
      </div>
    </form>

    ${renderFooter()}
  </main>

  <dialog id="embed-modal">
    <div id="modal-content">
      <span class="close" id="modal-close">&times;</span>
      <h3 class="modal-title">Share / Embed</h3>
      <div class="modal-section">
        <p class="modal-section-label">Direct Link</p>
        <a class="share-link-btn" id="share-link" href="" target="_blank" rel="noopener">Open Collage</a>
      </div>
      <div class="modal-section">
        <p class="modal-section-label">Embed Code</p>
        <p class="modal-hint">
          Paste this HTML to embed a live collage that updates automatically.
        </p>
      </div>
    </div>
    <div class="highlight" id="highlight">
      <button class="copy-code-button" type="button" id="copy-button">Copy</button>
      <pre class="chroma"><code id="embed-code"></code></pre>
    </div>
  </dialog>
`;

const form = document.getElementById("collage-form");
const usernameInput = document.getElementById("username");
const methodInput = document.getElementById("method");
const rowsInput = document.getElementById("rows");
const columnsInput = document.getElementById("columns");
const submitButton = document.getElementById("submit-button");
const loader = document.getElementById("form-loader");
const gridPreview = document.getElementById("grid-preview");
const advancedContent = document.getElementById("advanced-options-content");

const usernameError = document.getElementById("username-error");
const rowsError = document.getElementById("rows-error");
const columnsError = document.getElementById("columns-error");

const rowsMaxLabel = document.getElementById("rows-max");
const columnsMaxLabel = document.getElementById("columns-max");

const dialog = document.getElementById("embed-modal");
const shareLink = document.getElementById("share-link");
const embedCode = document.getElementById("embed-code");
const copyButton = document.getElementById("copy-button");

const state = {
  isSubmitting: false,
  submitted: false,
  touched: {
    username: false,
    rows: false,
    columns: false,
  },
  maxRows: 20,
  maxColumns: 20,
};

function toBoolean(value, defaultValue) {
  if (typeof value === "boolean") {
    return value;
  }
  if (typeof value === "string") {
    return value === "true";
  }
  if (typeof value === "number") {
    return value !== 0;
  }
  return defaultValue;
}

function toNumber(value, defaultValue) {
  const numericValue = Number(value);
  return Number.isFinite(numericValue) ? numericValue : defaultValue;
}

function loadPersistedValues() {
  const raw = localStorage.getItem(STORAGE_KEY);
  if (!raw) {
    return {};
  }

  try {
    const parsed = JSON.parse(raw);
    if (!parsed || typeof parsed !== "object") {
      return {};
    }

    if (parsed.data && typeof parsed.data === "object") {
      return parsed.data;
    }

    return parsed;
  } catch (_) {
    return {};
  }
}

function persistValues(values) {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(values));
}

function setFieldValue(name, value) {
  const field = form.elements.namedItem(name);
  if (!field) {
    return;
  }

  if (field.type === "checkbox") {
    field.checked = Boolean(value);
  } else {
    field.value = String(value);
  }
}

function applyValues(values) {
  const mergedValues = {
    ...DEFAULT_VALUES,
    ...values,
  };

  setFieldValue("username", mergedValues.username);
  setFieldValue("method", mergedValues.method);
  setFieldValue("period", mergedValues.period);
  setFieldValue("track", toBoolean(mergedValues.track, DEFAULT_VALUES.track));
  setFieldValue("artist", toBoolean(mergedValues.artist, DEFAULT_VALUES.artist));
  setFieldValue("album", toBoolean(mergedValues.album, DEFAULT_VALUES.album));
  setFieldValue(
    "playcount",
    toBoolean(mergedValues.playcount, DEFAULT_VALUES.playcount),
  );
  setFieldValue("rows", toNumber(mergedValues.rows, DEFAULT_VALUES.rows));
  setFieldValue("columns", toNumber(mergedValues.columns, DEFAULT_VALUES.columns));
  setFieldValue(
    "advancedOptions",
    toBoolean(mergedValues.advancedOptions, DEFAULT_VALUES.advancedOptions),
  );
  setFieldValue(
    "showTextSize",
    toBoolean(mergedValues.showTextSize, DEFAULT_VALUES.showTextSize),
  );
  setFieldValue(
    "showTextLocation",
    toBoolean(mergedValues.showTextLocation, DEFAULT_VALUES.showTextLocation),
  );
  setFieldValue("textSize", mergedValues.textSize || DEFAULT_VALUES.textSize);
  setFieldValue(
    "textLocation",
    mergedValues.textLocation || DEFAULT_VALUES.textLocation,
  );
  setFieldValue(
    "showBoldtext",
    toBoolean(mergedValues.showBoldtext, DEFAULT_VALUES.showBoldtext),
  );
  setFieldValue(
    "grayscaleImage",
    toBoolean(mergedValues.grayscaleImage, DEFAULT_VALUES.grayscaleImage),
  );
  setFieldValue(
    "WebPLossyCompression",
    toBoolean(
      mergedValues.WebPLossyCompression,
      DEFAULT_VALUES.WebPLossyCompression,
    ),
  );
}

function getParsedNumber(rawValue) {
  if (rawValue === "") {
    return null;
  }

  const numericValue = Number(rawValue);
  if (!Number.isFinite(numericValue)) {
    return null;
  }

  return numericValue;
}

function getValues() {
  return {
    username: form.elements.namedItem("username").value,
    method: form.elements.namedItem("method").value,
    period: form.elements.namedItem("period").value,
    track: form.elements.namedItem("track").checked,
    artist: form.elements.namedItem("artist").checked,
    album: form.elements.namedItem("album").checked,
    playcount: form.elements.namedItem("playcount").checked,
    rows: getParsedNumber(form.elements.namedItem("rows").value),
    columns: getParsedNumber(form.elements.namedItem("columns").value),
    advancedOptions: form.elements.namedItem("advancedOptions").checked,
    showTextSize: form.elements.namedItem("showTextSize").checked,
    showTextLocation: form.elements.namedItem("showTextLocation").checked,
    textSize: form.elements.namedItem("textSize").value,
    textLocation: form.elements.namedItem("textLocation").value,
    showBoldtext: form.elements.namedItem("showBoldtext").checked,
    grayscaleImage: form.elements.namedItem("grayscaleImage").checked,
    WebPLossyCompression: form.elements.namedItem("WebPLossyCompression").checked,
  };
}

function clampNumberInput(input, max, min = 0) {
  if (input.value === "") {
    return;
  }

  const numericValue = Number(input.value);
  if (!Number.isFinite(numericValue)) {
    return;
  }

  if (numericValue > max) {
    input.value = String(max);
  }

  if (numericValue < min) {
    input.value = String(min);
  }
}

function updateComputedState(values) {
  const styles = getComputedStyle(document.documentElement);
  const checkedColor = styles.getPropertyValue("--text").trim() || "#1a1a1a";
  const uncheckedColor =
    styles.getPropertyValue("--text-dim").trim() || "darkgrey";

  const showTrack = values.method === "track";
  const showAlbum = values.method !== "artist";

  state.maxRows = values.method === "album" ? 20 : 10;
  state.maxColumns = values.method === "album" ? 20 : 10;

  rowsInput.max = String(state.maxRows);
  columnsInput.max = String(state.maxColumns);
  rowsMaxLabel.textContent = String(state.maxRows);
  columnsMaxLabel.textContent = String(state.maxColumns);

  clampNumberInput(rowsInput, state.maxRows);
  clampNumberInput(columnsInput, state.maxColumns);

  document.getElementById("track-wrapper").style.display = showTrack
    ? "block"
    : "none";
  document.getElementById("album-wrapper").style.display = showAlbum
    ? "block"
    : "none";

  const advancedVisible = values.advancedOptions;
  advancedContent.hidden = !advancedVisible;

  document.getElementById("text-size-options").hidden =
    !advancedVisible || !values.showTextSize;
  document.getElementById("text-location-options").hidden =
    !advancedVisible || !values.showTextLocation;

  const webpContainer = document.getElementById("webp-container");
  webpContainer.hidden = values.grayscaleImage;

  // Render grid preview
  const previewRows = Math.min(Math.max(values.rows || 0, 0), state.maxRows);
  const previewCols = Math.min(
    Math.max(values.columns || 0, 0),
    state.maxColumns,
  );
  let gridHTML = "";
  if (previewRows > 0 && previewCols > 0) {
    gridHTML = `<div class="grid-preview-grid" style="grid-template-columns: repeat(${previewCols}, 1fr); grid-template-rows: repeat(${previewRows}, 1fr);">`;
    for (let i = 0; i < previewRows * previewCols; i++) {
      gridHTML += '<div class="grid-preview-cell"></div>';
    }
    gridHTML += "</div>";
    gridHTML += `<span class="grid-preview-label">${values.rows || 0} &times; ${values.columns || 0}</span>`;
  }
  gridPreview.innerHTML = gridHTML;

  document.querySelectorAll(".checkbox-wrapper").forEach((wrapper) => {
    const input = wrapper.querySelector('input[type="checkbox"]');
    const label = wrapper.querySelector(".checkbox-label");
    if (!input || !label) {
      return;
    }

    label.style.color = input.checked ? checkedColor : uncheckedColor;
  });
}

function validate(values) {
  const errors = {};

  if (!values.username) {
    errors.username = "Username is required";
  } else if (!USERNAME_REGEX.test(values.username)) {
    errors.username =
      "Username must be between 2 to 15 characters, begin with a letter and contain only letters, numbers, '_' or '-'";
  }

  if (values.rows === null) {
    errors.rows = "Number is required";
  } else if (!Number.isInteger(values.rows)) {
    errors.rows = "Expected integer, received float";
  } else if (values.rows < 1) {
    errors.rows = "Must be greater than 0";
  }

  if (values.columns === null) {
    errors.columns = "Number is required";
  } else if (!Number.isInteger(values.columns)) {
    errors.columns = "Expected integer, received float";
  } else if (values.columns < 1) {
    errors.columns = "Must be greater than 0";
  }

  return errors;
}

function shouldShowError(fieldName) {
  return state.submitted || state.touched[fieldName];
}

function renderErrors(errors) {
  const showUsernameError = Boolean(errors.username) && shouldShowError("username");
  usernameError.hidden = !showUsernameError;
  usernameError.textContent = showUsernameError ? errors.username : "";
  usernameInput.style.border = showUsernameError ? "2px solid red" : "";

  const showRowsError = Boolean(errors.rows) && shouldShowError("rows");
  rowsError.hidden = !showRowsError;
  rowsError.textContent = showRowsError ? errors.rows : "";
  rowsInput.classList.toggle("error", showRowsError);

  const showColumnsError =
    Boolean(errors.columns) && shouldShowError("columns");
  columnsError.hidden = !showColumnsError;
  columnsError.textContent = showColumnsError ? errors.columns : "";
  columnsInput.classList.toggle("error", showColumnsError);
}

function updateSubmittingState() {
  loader.hidden = !state.isSubmitting;
  submitButton.disabled = state.isSubmitting;
}

function generateUrl(values) {
  const params = new URLSearchParams();
  const showTrack = values.method === "track";
  const showAlbum = values.method !== "artist";

  params.append("username", values.username);
  params.append("method", values.method);
  params.append("period", values.period);

  if (showTrack) {
    params.append("track", String(values.track));
  }

  params.append("artist", String(values.artist));

  if (showAlbum) {
    params.append("album", String(values.album));
  }

  params.append("playcount", String(values.playcount));

  const rows = Math.min(values.rows ?? DEFAULT_VALUES.rows, state.maxRows);
  const columns = Math.min(
    values.columns ?? DEFAULT_VALUES.columns,
    state.maxColumns,
  );

  params.append("rows", String(rows));
  params.append("columns", String(columns));

  if (values.advancedOptions) {
    if (values.showTextSize) {
      params.append("fontsize", values.textSize);
    }

    if (values.showTextLocation) {
      params.append("textlocation", values.textLocation);
    }

    if (values.WebPLossyCompression) {
      params.append("webp", String(values.WebPLossyCompression));
    }

    if (values.showBoldtext) {
      params.append("boldfont", String(values.showBoldtext));
    }

    if (values.grayscaleImage) {
      params.append("grayscale", String(values.grayscaleImage));
    }
  }

  params.append("cacheid", Date.now().toString());
  return `/collage?${params.toString()}`;
}

function syncState() {
  const values = getValues();
  updateComputedState(values);
  const clampedValues = getValues();
  const errors = validate(clampedValues);
  renderErrors(errors);
  persistValues(clampedValues);
  return { values: clampedValues, errors };
}

function resetForm() {
  applyValues(DEFAULT_VALUES);
  state.submitted = false;
  state.touched = {
    username: false,
    rows: false,
    columns: false,
  };

  const { values } = syncState();
  persistValues(values);
}

function showEmbedModal() {
  const values = getValues();
  const url = `https://songstitch.art${generateUrl(values)}`;
  const embedHTML = `<img class="songstitch-collage" src="${url}">`;

  shareLink.href = url;
  embedCode.textContent = embedHTML;
  copyButton.textContent = "Copy";

  if (!dialog.open) {
    dialog.showModal();
  }
}

copyButton.addEventListener("click", () => {
  copyButton.textContent = "Copied!";

  const copyPromise = navigator.clipboard
    ? navigator.clipboard.writeText(embedCode.textContent)
    : Promise.resolve();

  copyPromise
    .catch(() => {})
    .finally(() => {
      window.setTimeout(() => {
        copyButton.textContent = "Copy";
      }, 2000);
    });
});

dialog.addEventListener("click", (event) => {
  if (event.target === dialog) {
    dialog.close();
  }
});

document.getElementById("modal-close").addEventListener("click", () => {
  dialog.close();
});

document.getElementById("embed-button").addEventListener("click", () => {
  showEmbedModal();
});

document.getElementById("reset-link").addEventListener("click", () => {
  resetForm();
});

form.addEventListener("input", (event) => {
  const fieldName = event.target && event.target.name;
  if (fieldName === "username" || fieldName === "rows" || fieldName === "columns") {
    state.touched[fieldName] = true;
  }

  syncState();
});

form.addEventListener("change", (event) => {
  const fieldName = event.target && event.target.name;
  if (fieldName === "username" || fieldName === "rows" || fieldName === "columns") {
    state.touched[fieldName] = true;
  }

  syncState();
});

form.addEventListener("submit", (event) => {
  event.preventDefault();

  state.submitted = true;
  state.isSubmitting = true;
  updateSubmittingState();

  const { values, errors } = syncState();

  if (Object.keys(errors).length > 0) {
    state.isSubmitting = false;
    updateSubmittingState();
    window.location.href = "#top";
    return;
  }

  const collageURL = generateUrl(values);
  window.open(collageURL, "_self");
});

window.addEventListener("pageshow", () => {
  state.isSubmitting = false;
  updateSubmittingState();
});

window.addEventListener("songstitch:darkmodechange", () => {
  syncState();
});

const storedValues = loadPersistedValues();
applyValues(storedValues);
syncState();
updateSubmittingState();
initialiseDarkMode();
