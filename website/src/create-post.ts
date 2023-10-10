import Quill from "quill";
import "./scripts/themes";

import "./styles/global.css";
import "./styles/header.css";
import "./create-post.css";

new Quill("#post-editor", {
    theme: "snow",
});
