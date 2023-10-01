import "htmx.org";
import "./scripts/themes";

import "./global.css";
import "./index.css";
import "./header.css";

import { fetchData } from "./lib/main";

const element = document.getElementById("data");

const result = fetchData();
console.log("Fetched data: " + result);

if (element) {
    element.innerText = result;
}
