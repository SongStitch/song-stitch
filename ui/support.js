import { renderFooter, renderHeader } from "./shared.js";

const support = document.getElementById("support");

support.innerHTML = `
  <main>
    ${renderHeader()}
    <h1>Support</h1>
    <h3 style="text-align: center">Thank you for using SongStitch</h3>
    <h4 style="text-align: center">
      If you have any issues or questions, please contact us at
      <a class="href-links" href="mailto:songstitchsupport@theden.sh"
        >supportsongstitch@theden.sh</a
      >
      or create a
      <a class="href-links" href="https://github.com/SongStitch/song-stitch/issues"
        >GitHub issue</a
      >
    </h4>
    ${renderFooter()}
  </main>
`;
