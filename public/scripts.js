// Hide Spinner when not submitting
window.addEventListener('pageshow', hideSpinners)
function hideSpinners () {
  document.getElementsByClassName('loader')[0].style.display = 'none'
  document.getElementsByClassName('btn-grad')[0].style.display = 'block'
  document.getElementsByClassName('btn-grad-embed')[0].style.display = 'block'
}

window.onload = function () {
  // init checkbox values
  document.getElementById('artist').value = 'true'
  document.getElementById('album').value = 'true'
  document.getElementById('playcount').value = 'true'

  // randomise credits
  const p = document.getElementById('links')
  const spans = p.getElementsByTagName('span')
  const spanArr = Array.prototype.slice.call(spans)
  spanArr.sort(function () {
    return 0.5 - Math.random()
  })
  spanArr.forEach(function (span) {
    p.appendChild(span) // This will move the <span> element (containing an <a> tag) to the end of the list.
  })
  // We need to adjust the "and" position after the first link
  const andNode = document.createTextNode(' and ')
  p.insertBefore(andNode, p.children[1])

  // LocalStorage for username
  const username = document.querySelector('#username')
  username.value = localStorage.getItem('username') || ''
  document.querySelector('#form').addEventListener('submit', function () {
    localStorage.setItem('username', username.value)
  })
}

function updateValue (checkbox) {
  checkbox.value = checkbox.checked ? 'true' : 'false'
}

// Embed button js
function embedUrl () {
  const form = document.getElementById('form')
  const action = form.action
  const elems = form.elements
  let url = action
  let first = true
  for (let i = 0; i < elems.length; i++) {
    if (elems[i].type === 'submit' || elems[i].name === 'embed' || elems[i].id === 'fieldset' || elems[i].id === 'aspectRatio' || elems[i].id === 'advanced') continue
    if (elems[i].name == 'height' && elems[i].value.length == 0) continue
    if (elems[i].name == 'width' && elems[i].value.length == 0) continue
    if (first) {
      url += '?'
      first = false
    } else {
      url += '&'
    }
    url += elems[i].name + '=' + encodeURIComponent(elems[i].value)
    var embedData = '<img class="songstitch-collage" src="' + url + '">'
  }
  document.getElementById('embedUrl').textContent = embedData
  modal.style.display = 'block'
  return false // prevent the form from submitting
}
function copyToClipboard () {
  const urlText = document.getElementById('embedUrl').textContent
  navigator.clipboard
    .writeText(urlText)
    .then(function () {
      console.log('Copied to clipboard')
    })
    .catch(function () {
      console.error('Failed to copy text')
    })
}
var modal = document.getElementById('myModal')
const span = document.getElementsByClassName('close')[0]
span.onclick = function () {
  modal.style.display = 'none'
}
window.onclick = function (event) {
  if (event.target == modal) {
    modal.style.display = 'none'
  }
}

// close modal when ESC key is pressed
document.addEventListener('keydown', function (event) {
  if (event.key === 'Escape') {
    modal.style.display = 'none'
  }
})

/// Copy embed button logic
function createCopyButton (highlightDiv) {
  const button = document.createElement('button')
  button.className = 'copy-code-button'
  button.type = 'button'
  button.innerText = 'Copy'
  button.addEventListener('click', () =>
    copyCodeToClipboard(button, highlightDiv)
  )
  addCopyButtonToDom(button, highlightDiv)
}

async function copyCodeToClipboard (button, highlightDiv) {
  const codeToCopy = highlightDiv.querySelector(
    ':last-child > .chroma > code'
  ).innerText
  try {
    result = await navigator.permissions.query({ name: 'clipboard-write' })
    if (result.state == 'granted' || result.state == 'prompt') {
      await navigator.clipboard.writeText(codeToCopy)
    } else {
      copyCodeBlockExecCommand(codeToCopy, highlightDiv)
    }
  } catch (_) {
    copyCodeBlockExecCommand(codeToCopy, highlightDiv)
  } finally {
    codeWasCopied(button)
  }
}

function copyCodeBlockExecCommand (codeToCopy, highlightDiv) {
  const textArea = document.createElement('textArea')
  textArea.contentEditable = 'true'
  textArea.readOnly = 'false'
  textArea.className = 'copyable-text-area'
  textArea.value = codeToCopy
  highlightDiv.insertBefore(textArea, highlightDiv.firstChild)
  const range = document.createRange()
  range.selectNodeContents(textArea)
  const sel = window.getSelection()
  sel.removeAllRanges()
  sel.addRange(range)
  textArea.setSelectionRange(0, 999999)
  document.execCommand('copy')
  highlightDiv.removeChild(textArea)
}

function codeWasCopied (button) {
  button.blur()
  button.innerText = 'Copied!'
  setTimeout(function () {
    button.innerText = 'Copy'
  }, 2000)
}

function addCopyButtonToDom (button, highlightDiv) {
  highlightDiv.insertBefore(button, highlightDiv.firstChild)
  const wrapper = document.createElement('div')
  wrapper.className = 'highlight-wrapper'
  highlightDiv.parentNode.insertBefore(wrapper, highlightDiv)
  wrapper.appendChild(highlightDiv)
}

document.querySelectorAll('.highlight').forEach((highlightDiv) => createCopyButton(highlightDiv))

// Loader logic
document.getElementById('form').addEventListener('submit', function () {
  event.preventDefault()
  submitForm(event.target)
  document.getElementsByClassName('loader')[0].style.display = 'block'
  document.getElementsByClassName('btn-grad')[0].style.display = 'none'
  document.getElementsByClassName('btn-grad-embed')[0].style.display = 'none'
})

// Advanced Options
function toggleAdvancedOptions (checkBoxElement) {
  const advancedOptions = document.getElementById('advanced-options')
  if (checkBoxElement.checked) {
    advancedOptions.style.display = 'block'
    if (typeof tempWidth !== 'undefined') {
      document.getElementById('width').value = tempWidth
    } else {
      document.getElementById('width').value = 500
    }
    if (typeof tempHeight !== 'undefined') {
      document.getElementById('height').value = tempHeight
    } else {
      document.getElementById('height').value = 500
    }
  } else {
    advancedOptions.style.display = 'none'
    tempWidth = document.getElementById('width').value
    tempHeight = document.getElementById('height').value
    document.getElementById('width').value = ''
    document.getElementById('height').value = ''
  }
}

let aspectRatioChecked = false
document.getElementById('aspectRatio').addEventListener('change', function () {
  aspectRatioChecked = this.checked
  maintainAspectRatio('aspectRatio')
})

function maintainAspectRatio (input) {
  if (aspectRatioChecked) {
    numCols = document.getElementById('columns').value
    numRows = document.getElementById('rows').value
    if (Math.round(numRows) === 0 || Math.round(numCols) === 0) {
      return
    }
    if (input.id === 'width') {
      value = Math.round((input.value * numRows) / numCols)
      document.getElementById('height').value = value
    } else if (input.id === 'height') {
      value = Math.round((input.value * numCols) / numRows)
      document.getElementById('width').value = value
    } else {
      value = Math.round((document.getElementById('width').value * numRows) / numCols)
      document.getElementById('height').value = value
      value = Math.round((document.getElementById('height').value * numCols) / numRows)
      document.getElementById('width').value = value
    }
  }
}

function submitForm (form) {
  const params = new URLSearchParams()
  for (const field of form.elements) {
    if (field.name && field.value && field.name !== 'submit' && field.name !== 'advanced' && field.name !== 'aspectRatio' && field.name !== 'embed') {
      params.append(field.name, field.value)
    }
    window.location.href = '/collage?' + params.toString()
  }
}
