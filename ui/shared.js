const LOGO_LIGHT = "/assets/images/songstitch_logo.png";
const LOGO_DARK = "/assets/images/songstitch_logo_dark.png";
const APPSTORE_ICON = "/assets/images/appstore_icon.png";
const KOFI_CUP = "/assets/images/cup-border.webp";
const DARK_MODE_ICON = "/assets/images/dark_mode.png";

function getContributorNames() {
  const names = ["BradLewis", "TheDen"];
  if (Math.random() > 0.5) {
    names.reverse();
  }
  return names;
}

function getContributorLink(name) {
  if (name === "TheDen") {
    return "https://theden.sh";
  }
  return `https://github.com/${name}`;
}

export function renderHeader() {
  return `
    <div class="header-img-container" id="top">
      <a class="img-link" href="/">
        <img
          src="${LOGO_LIGHT}"
          class="header-img"
          alt="SongStitch Logo"
          width="418"
          height="100"
        />
      </a>
    </div>
    <p class="subheading">Share your most played Last.FM music with an image collage!</p>
  `;
}

export function renderDarkModeToggle() {
  return `
    <div class="dark-mode-icon">
      <a href="#" class="img-link" id="dark-mode-toggle">
        <img
          src="${DARK_MODE_ICON}"
          class="darkmode-icon-img"
          alt="darkmode icon"
          width="512"
          height="512"
        />
      </a>
    </div>
  `;
}

export function initialiseDarkMode() {
  let darkMode = localStorage.getItem("darkMode") === "true";
  const toggleLink = document.getElementById("dark-mode-toggle");

  const applyDarkMode = () => {
    document.body.classList.toggle("dark-mode", darkMode);
    document.documentElement.classList.toggle("dark-mode", darkMode);

    const headerImg = document.querySelector(".header-img");
    if (headerImg) {
      headerImg.classList.toggle("dark-mode", darkMode);
      headerImg.src = darkMode ? LOGO_DARK : LOGO_LIGHT;
    }

    localStorage.setItem("darkMode", String(darkMode));
    window.dispatchEvent(
      new CustomEvent("songstitch:darkmodechange", {
        detail: { darkMode },
      }),
    );
  };

  if (darkMode) {
    applyDarkMode();
  }

  if (toggleLink) {
    toggleLink.addEventListener("click", (event) => {
      event.preventDefault();
      darkMode = !darkMode;
      applyDarkMode();
    });
  }
}

export function renderFooter() {
  const names = getContributorNames();

  return `
    <div class="footer">
      <div class="kofi-container">
        <span id="spanPreview" class="imgPreview">
          <div class="btn-container">
            <a
              title="Support me on ko-fi.com"
              class="kofi-button"
              href="https://ko-fi.com/P5P3NPGIU"
              target="_blank"
              rel="noopener"
            >
              <span class="kofitext">
                <img src="${KOFI_CUP}" alt="Ko-fi donations" class="kofiimg" />
                Support Us on Ko-fi
              </span>
            </a>
          </div>
        </span>
      </div>
      <p class="appstore-container">
        <a
          href="https://apps.apple.com/au/app/songstitch/id6450189672"
          target="_blank"
          rel="noopener"
        >
          <img
            class="appstore-icon"
            src="${APPSTORE_ICON}"
            alt="iOS App store link"
            width="150"
            height="50"
          />
        </a>
      </p>
      <p id="links" class="footer-text">
        Made with &#10084;&#65039; by
        <span>
          <a
            class="href-links"
            href="${getContributorLink(names[0])}"
            target="_blank"
            rel="noopener"
          >${names[0]}</a>
        </span>
        and
        <span>
          <a
            class="href-links"
            href="${getContributorLink(names[1])}"
            target="_blank"
            rel="noopener"
          >${names[1]}</a>
        </span>
      </p>
      <p class="gh-footer">
        <a
          class="gh-link"
          href="https://github.com/SongStitch/song-stitch"
          target="_blank"
          rel="noopener"
          aria-label="SongStitch GitHub repository"
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24">
            <path
              d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"
            />
          </svg>
        </a>
      </p>
      <p class="lastfm-footer">
        <a class="lastfm-link" href="https://last.fm/" target="_blank" rel="noopener">
          <span class="nonbold">Powered by</span>
          <svg
            fill="#c42b1d"
            width="20"
            height="12"
            viewBox="0 0 20 24"
            overflow="visible"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              d="M14.131 22.948l-1.172-3.193c0 0-1.912 2.131-4.771 2.131-2.537 0-4.333-2.203-4.333-5.729 0-4.511 2.276-6.125 4.515-6.125 3.224 0 4.245 2.089 5.125 4.772l1.161 3.667c1.161 3.561 3.365 6.421 9.713 6.421 4.548 0 7.631-1.391 7.631-5.068 0-2.968-1.697-4.511-4.844-5.244l-2.344-0.511c-1.624-0.371-2.104-1.032-2.104-2.131 0-1.249 0.985-1.984 2.604-1.984 1.767 0 2.704 0.661 2.865 2.24l3.661-0.444c-0.297-3.301-2.584-4.656-6.323-4.656-3.308 0-6.532 1.251-6.532 5.245 0 2.5 1.204 4.077 4.245 4.807l2.484 0.589c1.865 0.443 2.484 1.224 2.484 2.287 0 1.359-1.323 1.921-3.828 1.921-3.703 0-5.244-1.943-6.124-4.625l-1.204-3.667c-1.541-4.765-4.005-6.531-8.891-6.531-5.287-0.016-8.151 3.385-8.151 9.192 0 5.573 2.864 8.595 8.005 8.595 4.14 0 6.125-1.943 6.125-1.943z"
            />
          </svg>
        </a>
      </p>
    </div>
  `;
}
