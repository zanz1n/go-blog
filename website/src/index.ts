import "htmx.org";
import "./global.css";
import "./index.css";

import { fetchData } from "./lib/main";

const element = document.getElementById("data");

const result = fetchData();
console.log("Fetched data: " + result);

if (element) {
    element.innerText = result;
}
