window.addEventListener('pageshow', () => toggleLoader(false));

// Form Submission
function submitForm(form) {
  const IGNORED_FIELDS = ['submit', 'advanced', 'aspectRatio', 'embed'];
  const params = new URLSearchParams();

  Array.from(form.elements)
    .filter(
      (field) =>
        field.name && field.value && !IGNORED_FIELDS.includes(field.name)
    )
    .forEach((field) => params.append(field.name, field.value));

  window.location.href = '/collage?' + params.toString();
}

const formElement = document.getElementById('form');
if (formElement) {
  formElement.addEventListener('submit', handleFormSubmit);
}

// Event handler for form submit
function handleFormSubmit(event) {
  event.preventDefault();
  toggleLoader(true);
  submitForm(event.target);
}

/**
 * Toggle the visibility of the loader and associated buttons
 * @param {boolean} isLoading - Whether the loader should be visible
 */
function toggleLoader(isLoading) {
  const elementClasses = [
    'loader-container',
    'loader',
    'btn-grad',
    'btn-grad-embed',
  ];
  const displayValue = isLoading
    ? ['grid', 'block', 'none', 'none']
    : ['none', 'none', 'block', 'block'];

  elementClasses.forEach((className, index) => {
    const element = document.getElementsByClassName(className)[0];
    if (element) {
      element.style.display = displayValue[index];
    } else {
      console.error(
        `Element with class ${className} could not be found in the DOM.`
      );
    }
  });
}

window.addEventListener('DOMContentLoaded', initializePage);

// Function to initialize the page after the DOM has been loaded
function initializePage() {
  initCheckboxValues();
  randomizeCredits();
  handleLocalStorage();
}

function initCheckboxValues() {
  ['artist', 'album', 'playcount'].forEach((id) => {
    const element = document.getElementById(id);
    if (element) {
      element.value = 'true';
    }
  });
  document.getElementById('compress').value = '';
}

function randomizeCredits() {
  const p = document.getElementById('links');
  const spanArr = Array.from(p.getElementsByTagName('span'));
  spanArr.sort(() => Math.random() - 0.5);
  spanArr.forEach((span) => p.appendChild(span));
  const andNode = document.createTextNode(' and ');
  if (p.children[1]) {
    p.insertBefore(andNode, p.children[1]);
  }
}

function handleLocalStorage() {
  const username = document.getElementById('username');
  if (username && formElement) {
    username.value = localStorage.getItem('username') || '';
    formElement.addEventListener('submit', () => {
      localStorage.setItem('username', username.value);
    });
  }
}

// checkbox value
function updateValue(checkbox) {
  checkbox.value = checkbox.checked ? 'true' : 'false';
}

// Embed button and modal
function embedUrl() {
  form = document.getElementById('form');
  action = form.action;
  elems = Array.from(form.elements);

  const filteredElems = elems.filter((el) => {
    excludedIds = ['fieldset', 'aspectRatio', 'advanced'];
    excludedNames = ['embed'];
    if (document.getElementById('width').value.length == 0) {
      excludedNames.push('width');
    }
    if (document.getElementById('height').value.length == 0) {
      excludedNames.push('height');
    }
    if (document.getElementById('compress').value == 'false') {
      excludedNames.push('compress');
    }
    return !(
      el.type === 'submit' ||
      excludedIds.includes(el.id) ||
      excludedNames.includes(el.name) ||
      el.value === ''
    );
  });

  const query = filteredElems
    .map((el) => `${el.name}=${encodeURIComponent(el.value)}`)
    .join('&');
  const url = `${action}?${query}`;

  const embedData = `<img class="songstitch-collage" src="${url}">`;

  document.getElementById('embedUrl').textContent = embedData;
  displayModal();

  return false; // prevent the form from submitting
}

function displayModal() {
  const modal = document.getElementById('modal');
  modal.style.display = 'block';
}

// copy to Clipboard
function copyToClipboard() {
  const urlText = document.getElementById('embedUrl').textContent;
  navigator.clipboard
    .writeText(urlText)
    .then(function () {})
    .catch(function () {
      console.error('Failed to copy text');
    });
}

const modal = document.getElementById('modal');
const span = document.getElementsByClassName('close')[0];
span.onclick = function () {
  modal.style.display = 'none';
};
// close modal when user clicks outside
window.onclick = function (event) {
  if (event.target == modal) {
    modal.style.display = 'none';
  }
};

// close modal when ESC key is pressed
document.addEventListener('keydown', function (event) {
  if (event.key === 'Escape') {
    modal.style.display = 'none';
  }
});

/// Copy embed button logic
function createCopyButton(highlightDiv) {
  const button = document.createElement('button');
  button.className = 'copy-code-button';
  button.type = 'button';
  button.innerText = 'Copy';
  button.addEventListener('click', () =>
    copyCodeToClipboard(button, highlightDiv)
  );
  addCopyButtonToDom(button, highlightDiv);
}

async function copyCodeToClipboard(button, highlightDiv) {
  const codeToCopy = highlightDiv.querySelector(
    ':last-child > .chroma > code'
  ).innerText;
  try {
    result = await navigator.permissions.query({ name: 'clipboard-write' });
    if (result.state == 'granted' || result.state == 'prompt') {
      await navigator.clipboard.writeText(codeToCopy);
    } else {
      copyCodeBlockExecCommand(codeToCopy, highlightDiv);
    }
  } catch (_) {
    copyCodeBlockExecCommand(codeToCopy, highlightDiv);
  } finally {
    codeWasCopied(button);
  }
}

function copyCodeBlockExecCommand(codeToCopy, highlightDiv) {
  const textArea = document.createElement('textArea');
  textArea.contentEditable = 'true';
  textArea.readOnly = 'false';
  textArea.className = 'copyable-text-area';
  textArea.value = codeToCopy;
  highlightDiv.insertBefore(textArea, highlightDiv.firstChild);
  const range = document.createRange();
  range.selectNodeContents(textArea);
  const sel = window.getSelection();
  sel.removeAllRanges();
  sel.addRange(range);
  textArea.setSelectionRange(0, 999999);
  document.execCommand('copy');
  highlightDiv.removeChild(textArea);
}

function codeWasCopied(button) {
  button.blur();
  button.innerText = 'Copied!';
  setTimeout(function () {
    button.innerText = 'Copy';
  }, 2000);
}

function addCopyButtonToDom(button, highlightDiv) {
  highlightDiv.insertBefore(button, highlightDiv.firstChild);
  const wrapper = document.createElement('div');
  wrapper.className = 'highlight-wrapper';
  highlightDiv.parentNode.insertBefore(wrapper, highlightDiv);
  wrapper.appendChild(highlightDiv);
}

document
  .querySelectorAll('.highlight')
  .forEach((highlightDiv) => createCopyButton(highlightDiv));

// Advanced Options
function toggleAdvancedOptions(checkBoxElement) {
  const advancedOptions = document.getElementById('advanced-options');
  if (checkBoxElement.checked) {
    advancedOptions.style.display = 'block';
    aspectRatioChecked = document.getElementById('aspectRatio').checked = true;
    validate('aspectRatio');
    if (typeof tempWidth !== 'undefined') {
      document.getElementById('width').value = tempWidth;
    } else {
      document.getElementById('width').value = 1500;
    }
    if (typeof tempHeight !== 'undefined') {
      document.getElementById('height').value = tempHeight;
    } else {
      document.getElementById('height').value = 1500;
    }
  } else {
    aspectRatioChecked = document.getElementById('aspectRatio').checked = false;
    validate('aspectRatio');
    advancedOptions.style.display = 'none';
    tempWidth = document.getElementById('width').value;
    tempHeight = document.getElementById('height').value;
    document.getElementById('width').value = '';
    document.getElementById('height').value = '';
  }
}

// input validation
maxResolution = document.getElementById('width').getAttribute('max');
maxGridSize = document.getElementById('rows').getAttribute('max');

function checkGridValues(inputValue, min = 0) {
  if (inputValue > maxGridSize) {
    return maxGridSize;
  } else if (inputValue < min) {
    return min;
  }
  return inputValue;
}

function checkAspectRatioValues(inputValue, min = 0) {
  if (inputValue > maxResolution) {
    return maxResolution;
  } else if (inputValue < min) {
    return min;
  }
  return inputValue;
}

function updateAndValidateValue(id, checkFunction) {
  const element = document.getElementById(id);
  const value = checkFunction(Number(element.value));
  element.value = value;
  return value;
}

function validate(input) {
  let numCols = updateAndValidateValue('columns', checkGridValues);
  let numRows = updateAndValidateValue('rows', checkGridValues);

  if (aspectRatioChecked) {
    height = updateAndValidateValue('height', checkAspectRatioValues);
    width = updateAndValidateValue('width', checkAspectRatioValues);
    numCols = document.getElementById('columns').value;
    numRows = document.getElementById('rows').value;
    height = document.getElementById('height').value;
    width = document.getElementById('width').value;
    if (
      Math.round(numRows) === 0 ||
      Math.round(numCols) === 0 ||
      Math.round(height) === 0 ||
      Math.round(width) === 0
    ) {
      return;
    }
    if (input.id === 'width') {
      value = Math.round((input.value * numRows) / numCols);
      document.getElementById('height').value = value;
    } else if (input.id === 'height') {
      value = Math.round((input.value * numCols) / numRows);
      document.getElementById('width').value = value;
    } else if (height > width) {
      value = Math.round((width * numRows) / numCols);
      document.getElementById('height').value = value;
    } else if (width >= height) {
      value = Math.round((height * numCols) / numRows);
      document.getElementById('width').value = value;
    }
  }
}

let aspectRatioChecked = (document.getElementById(
  'aspectRatio'
).checked = false);
document.getElementById('aspectRatio').addEventListener('change', function () {
  aspectRatioChecked = this.checked;
  validate('aspectRatio');
});

const maxForArtist = 10;
const maxForTrack = 15;
const maxForAlbum = 15;

function setInputValues(max) {
  document.querySelector('#rows').max = max;
  document.querySelector('#columns').max = max;
  maxText = document.getElementsByClassName('maxvalues');
  for (let i = 0; i < maxText.length; i++) {
    maxText[i].innerHTML = '(max. ' + max + ')';
  }
}

function checkCollageValue() {
  var selectBox = document.getElementById('method');
  var selectedValue = selectBox.options[selectBox.selectedIndex].value;
  if (selectedValue === 'artist') {
    document.querySelector(
      '#fieldset > div.checkbox-wrapper.album-checkbox'
    ).style.display = 'none';
    setInputValues(maxForArtist);
  } else if (selectedValue === 'track') {
    document.querySelector(
      '#fieldset > div.checkbox-wrapper.album-checkbox'
    ).style.display = 'block';
    setInputValues(maxForTrack);
  } else {
    document.querySelector(
      '#fieldset > div.checkbox-wrapper.album-checkbox'
    ).style.display = 'block';
    setInputValues(maxForAlbum);
  }
}
