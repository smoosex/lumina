// Plugins

// Composables
import { createApp } from "vue";
import { registerPlugins } from "@/plugins";
import { applyMetadata } from "@/utils/metadata";
// Components
import App from "./App.vue";

// Styles
import "./style.css";

const app = createApp(App);

applyMetadata();

registerPlugins(app);

app.mount("#app");
