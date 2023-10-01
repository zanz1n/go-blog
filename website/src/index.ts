import "htmx.org";

import { fetchData } from "./lib/main";
import "./global.css";
import "./index.css";

const element = document.getElementById("data");

const result = fetchData();
console.log("Fetched data: " + result);

if (element) {
    element.innerText = result;
}
